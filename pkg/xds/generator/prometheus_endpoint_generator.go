package generator

import (
	"net"
	"strings"

	"github.com/pkg/errors"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	manager_dataplane "github.com/kumahq/kuma/pkg/core/managers/apis/dataplane"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	core_xds "github.com/kumahq/kuma/pkg/core/xds"
	xds_context "github.com/kumahq/kuma/pkg/xds/context"
	envoy_common "github.com/kumahq/kuma/pkg/xds/envoy"
	envoy_clusters "github.com/kumahq/kuma/pkg/xds/envoy/clusters"
	envoy_listeners "github.com/kumahq/kuma/pkg/xds/envoy/listeners"
	envoy_names "github.com/kumahq/kuma/pkg/xds/envoy/names"
)

// OriginPrometheus is a marker to indicate by which ProxyGenerator resources were generated.
const OriginPrometheus = "prometheus"

// PrometheusEndpointGenerator generates an inbound Envoy listener
// that forwards HTTP requests into the `/stats/prometheus`
// endpoint of the Envoy Admin API.
//
// When generating such a listener, it's important not to overshadow
// a port that is already in use by the application or other Envoy listeners.
// In the latter case we prefer not generate Prometheus endpoint at all
// rather than introduce undeterministic behavior.
type PrometheusEndpointGenerator struct {
}

var HasMetricsHijacker = func(version *mesh_proto.Version) bool {
	return strings.HasPrefix(version.GetKumaDp().GetVersion(), "1.2.")
}

func (g PrometheusEndpointGenerator) Generate(ctx xds_context.Context, proxy *core_xds.Proxy) (*core_xds.ResourceSet, error) {
	prometheusEndpoint, err := proxy.Dataplane.GetPrometheusEndpoint(ctx.Mesh.Resource)
	if err != nil {
		return nil, errors.Wrap(err, "could not get prometheus endpoint")
	}
	if prometheusEndpoint == nil {
		// Prometheus metrics must be enabled Mesh-wide for Prometheus endpoint to be generated.
		return nil, nil
	}
	if proxy.Metadata.GetAdminPort() == 0 {
		// It's not possible to export Prometheus metrics if Envoy Admin API has not been enabled on that dataplane.

		// TODO(yskopets): find a way to communicate this to users
		return nil, nil
	}

	prometheusEndpointAddress := proxy.Dataplane.Spec.GetNetworking().Address
	prometheusEndpointIP := net.ParseIP(prometheusEndpointAddress)

	if proxy.Dataplane.UsesInterface(prometheusEndpointIP, prometheusEndpoint.Port) {
		// If the Prometheus endpoint would otherwise overshadow one of interfaces of that Dataplane,
		// we prefer not to do that.

		// TODO(yskopets): find a way to communicate this to users
		return nil, nil
	}

	resources := core_xds.NewResourceSet()

	// backwards compatibility
	// Kuma DP <1.2.x does not have metrics hijacker builtin therefore we configure those Envoys with just Envoy prometheus endpoint.
	var clusterName string
	var statsPath string

	if HasMetricsHijacker(proxy.Metadata.Version) {
		metricsHijackerClusterName := envoy_names.GetMetricsHijackerClusterName()
		cluster, err := envoy_clusters.NewClusterBuilder(proxy.APIVersion).
			Configure(envoy_clusters.StaticClusterUnixSocket(metricsHijackerClusterName,
				envoy_common.MetricsHijackerSocketName(proxy.Dataplane.Meta.GetName(), proxy.Dataplane.Meta.GetMesh()))).
			Build()
		if err != nil {
			return nil, err
		}

		resources.Add(&core_xds.Resource{
			Name:     cluster.GetName(),
			Origin:   OriginPrometheus,
			Resource: cluster,
		})

		clusterName = metricsHijackerClusterName
		statsPath = "/"
	} else {
		clusterName = envoy_names.GetEnvoyAdminClusterName()
		statsPath = "/stats/prometheus"
	}

	// Cluster is generated by AdminProxyGenerator
	prometheusListenerName := envoy_names.GetPrometheusListenerName()

	inbound, err := manager_dataplane.PrometheusInbound(proxy.Dataplane, ctx.Mesh.Resource)
	if err != nil {
		return nil, errors.Wrap(err, "could not get prometheus inbound interface")
	}

	iface := proxy.Dataplane.Spec.GetNetworking().ToInboundInterface(inbound)
	var listener envoy_common.NamedResource
	if secureMetrics(prometheusEndpoint, ctx.Mesh.Resource) {
		listener, err = envoy_listeners.NewListenerBuilder(proxy.APIVersion).
			Configure(envoy_listeners.InboundListener(prometheusListenerName, prometheusEndpointAddress, prometheusEndpoint.Port, core_xds.SocketAddressProtocolTCP)).
			// generate filter chain that does not require mTLS when DP scrapes itself (for example DP next to Prometheus Server)
			Configure(envoy_listeners.FilterChain(envoy_listeners.NewFilterChainBuilder(proxy.APIVersion).
				Configure(envoy_listeners.SourceMatcher(proxy.Dataplane.Spec.GetNetworking().Address)).
				Configure(envoy_listeners.StaticEndpoints(prometheusListenerName,
					[]*envoy_common.StaticEndpointPath{
						{
							ClusterName: clusterName,
							Path:        prometheusEndpoint.Path,
							RewritePath: statsPath,
						},
					})),
			)).
			Configure(envoy_listeners.FilterChain(envoy_listeners.NewFilterChainBuilder(proxy.APIVersion).
				Configure(envoy_listeners.StaticEndpoints(prometheusListenerName, []*envoy_common.StaticEndpointPath{
					{
						ClusterName: clusterName,
						Path:        prometheusEndpoint.Path,
						RewritePath: statsPath,
					},
				})).
				Configure(envoy_listeners.ServerSideMTLS(ctx)).
				Configure(envoy_listeners.NetworkRBAC(prometheusListenerName, ctx.Mesh.Resource.MTLSEnabled(), proxy.Policies.TrafficPermissions[iface]...)),
			)).
			Build()
	} else {
		listener, err = envoy_listeners.NewListenerBuilder(proxy.APIVersion).
			Configure(envoy_listeners.InboundListener(prometheusListenerName, prometheusEndpointAddress, prometheusEndpoint.Port, core_xds.SocketAddressProtocolTCP)).
			Configure(envoy_listeners.FilterChain(envoy_listeners.NewFilterChainBuilder(proxy.APIVersion).
				Configure(envoy_listeners.StaticEndpoints(prometheusListenerName, []*envoy_common.StaticEndpointPath{
					{
						ClusterName: clusterName,
						Path:        prometheusEndpoint.Path,
						RewritePath: statsPath,
					},
				})),
			)).
			Build()
	}
	if err != nil {
		return nil, err
	}

	resources.Add(&core_xds.Resource{
		Name:     listener.GetName(),
		Origin:   OriginPrometheus,
		Resource: listener,
	})
	return resources, nil
}

func secureMetrics(cfg *mesh_proto.PrometheusMetricsBackendConfig, mesh *core_mesh.MeshResource) bool {
	return !cfg.SkipMTLS.GetValue() && mesh.MTLSEnabled()
}
