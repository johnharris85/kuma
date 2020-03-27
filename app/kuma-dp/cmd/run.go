package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Kong/kuma/pkg/catalog"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	kumadp_config "github.com/Kong/kuma/app/kuma-dp/pkg/config"
	"github.com/Kong/kuma/app/kuma-dp/pkg/dataplane/accesslogs"
	"github.com/Kong/kuma/app/kuma-dp/pkg/dataplane/envoy"
	"github.com/Kong/kuma/app/kuma-dp/pkg/dataplane/sds"
	catalog_client "github.com/Kong/kuma/pkg/catalog/client"
	kuma_cmd "github.com/Kong/kuma/pkg/cmd"
	"github.com/Kong/kuma/pkg/config"
	kuma_dp "github.com/Kong/kuma/pkg/config/app/kuma-dp"
	config_types "github.com/Kong/kuma/pkg/config/types"
	"github.com/Kong/kuma/pkg/core"
	"github.com/Kong/kuma/pkg/core/runtime/component"
	util_net "github.com/Kong/kuma/pkg/util/net"
)

type CatalogClientFactory func(string) (catalog_client.CatalogClient, error)

var (
	runLog = dataplaneLog.WithName("run")
	// overridable by tests
	bootstrapGenerator   = envoy.NewRemoteBootstrapGenerator(&http.Client{Timeout: 10 * time.Second})
	catalogClientFactory = catalog_client.NewCatalogClient
)

func newRunCmd() *cobra.Command {
	cfg := kuma_dp.DefaultConfig()
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launch Dataplane (Envoy)",
		Long:  `Launch Dataplane (Envoy).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// only support configuration via environment variables and args
			if err := config.Load("", &cfg); err != nil {
				runLog.Error(err, "unable to load configuration")
				return err
			}
			if conf, err := config.ToYAML(&cfg); err == nil {
				runLog.Info("effective configuration", "config", string(conf))
			} else {
				runLog.Error(err, "unable to format effective configuration", "config", cfg)
				return err
			}

			catalogClient, err := catalogClientFactory(cfg.ControlPlane.ApiServer.URL)
			if err != nil {
				return errors.Wrap(err, "could not create catalog client")
			}
			catalog, err := catalogClient.Catalog()
			if err != nil {
				return errors.Wrap(err, "could retrieve catalog")
			}
			if isDpTokenRequired(cfg, catalog) {
				if cfg.DataplaneRuntime.TokenPath == "" {
					return errors.New("Kuma CP is configured with Dataplane Token Server therefore the Dataplane Token is required. " +
						"Generate token using 'kumactl generate dataplane-token > /path/file' and provide it via --dataplane-token-file=/path/file argument to Kuma DP")
				}
				if err := kumadp_config.ValidateTokenPath(cfg.DataplaneRuntime.TokenPath); err != nil {
					return err
				}
			}

			if !cfg.Dataplane.AdminPort.Empty() {
				// unless a user has explicitly opted out of Envoy Admin API, pick a free port from the range
				adminPort, err := util_net.PickTCPPort("127.0.0.1", cfg.Dataplane.AdminPort.Lowest(), cfg.Dataplane.AdminPort.Highest())
				if err != nil {
					return errors.Wrapf(err, "unable to find a free port in the range %q for Envoy Admin API to listen on", cfg.Dataplane.AdminPort)
				}
				cfg.Dataplane.AdminPort = config_types.MustExactPort(adminPort)
				runLog.Info("picked a free port for Envoy Admin API to listen on", "port", cfg.Dataplane.AdminPort)
			}

			if cfg.DataplaneRuntime.ConfigDir == "" {
				tmpDir, err := ioutil.TempDir("", "kuma-dp-")
				if err != nil {
					runLog.Error(err, "unable to create a temporary directory to store generated Envoy config at")
					return err
				}
				defer func() {
					if err := os.RemoveAll(tmpDir); err != nil {
						runLog.Error(err, "unable to remove a temporary directory with a generated Envoy config")
					}
				}()
				cfg.DataplaneRuntime.ConfigDir = tmpDir
				runLog.Info("generated Envoy configuration will be stored in a temporary directory", "dir", tmpDir)
			}

			if cfg.SDS.Address == "" {
				cfg.SDS.Address = fmt.Sprintf("unix:///tmp/kuma-sds-%s-%s.sock", cfg.Dataplane.Name, cfg.Dataplane.Mesh)
			}

			bootstrapConfig, err := bootstrapGenerator(catalog.Apis.Bootstrap.Url, cfg)
			if err != nil {
				return errors.Wrapf(err, "failed to generate Envoy bootstrap config")
			}
			dataplane := envoy.New(envoy.Opts{
				Config:          cfg,
				BootstrapConfig: bootstrapConfig,
				Stdout:          cmd.OutOrStdout(),
				Stderr:          cmd.OutOrStderr(),
			})
			server := accesslogs.NewAccessLogServer(cfg.Dataplane)

			componentMgr := component.NewManager()
			if err := componentMgr.Add(server, dataplane); err != nil {
				return err
			}
			if cfg.SDS.Type == kuma_dp.SdsDpVault {
				sdsServer, err := sds.NewVaultSdsServer(cfg, bootstrapConfig.Node.Cluster)
				if err != nil {
					return errors.Wrap(err, "could not create Vault SDS Server")
				}
				if err := componentMgr.Add(sdsServer); err != nil {
					return err
				}
			}

			runLog.Info("starting Kuma DP")
			if err := componentMgr.Start(core.SetupSignalHandler()); err != nil {
				runLog.Error(err, "error while running Kuma DP")
				return err
			}
			runLog.Info("stopping Kuma DP")
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfg.Dataplane.Name, "name", cfg.Dataplane.Name, "Name of the Dataplane")
	cmd.PersistentFlags().Var(&cfg.Dataplane.AdminPort, "admin-port", `Port (or range of ports to choose from) for Envoy Admin API to listen on. Empty value indicates that Envoy Admin API should not be exposed over TCP. Format: "9901 | 9901-9999 | 9901- | -9901"`)
	cmd.PersistentFlags().StringVar(&cfg.Dataplane.Mesh, "mesh", cfg.Dataplane.Mesh, "Mesh that Dataplane belongs to")
	cmd.PersistentFlags().StringVar(&cfg.ControlPlane.ApiServer.URL, "cp-address", cfg.ControlPlane.ApiServer.URL, "URL of the Control Plane API Server")
	cmd.PersistentFlags().StringVar(&cfg.DataplaneRuntime.BinaryPath, "binary-path", cfg.DataplaneRuntime.BinaryPath, "Binary path of Envoy executable")
	cmd.PersistentFlags().StringVar(&cfg.DataplaneRuntime.ConfigDir, "config-dir", cfg.DataplaneRuntime.ConfigDir, "Directory in which Envoy config will be generated")
	cmd.PersistentFlags().StringVar(&cfg.DataplaneRuntime.TokenPath, "dataplane-token-file", cfg.DataplaneRuntime.TokenPath, "Path to a file with dataplane token (use 'kumactl generate dataplane-token' to get one)")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Type, "sds-type", cfg.SDS.Type, kuma_cmd.UsageOptions("Type of Secret Discovery Server that DP will work with", kuma_dp.SdsDpVault, kuma_dp.SdsCp))
	cmd.PersistentFlags().StringVar(&cfg.SDS.Address, "sds-address", cfg.SDS.Address, `Address of the SDS server. Only Unix socket is supported for now ex. "unix:///tmp/server.sock"`)
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.Address, "vault-address", cfg.SDS.Vault.Address, "Address of Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.AgentAddress, "vault-agent-address", cfg.SDS.Vault.AgentAddress, "Agent Address of Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.Token, "vault-token", cfg.SDS.Vault.Token, "Token used for authentication with Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.Namespace, "vault-namespace", cfg.SDS.Vault.Namespace, "Vault namespace")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.TLS.CaCertPath, "vault-ca-cert-path", cfg.SDS.Vault.TLS.CaCertPath, "Path to TLS certificate that will be used to connect to Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.TLS.CaCertDir, "vault-ca-cert-dir", cfg.SDS.Vault.TLS.CaCertDir, "Path to directory of TLS certificates that will be used to connect to Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.TLS.ClientCertPath, "vault-client-cert-path", cfg.SDS.Vault.TLS.ClientCertPath, "Path to client TLS certificate that will be used to connect to Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.TLS.ClientKeyPath, "vault-client-key-path", cfg.SDS.Vault.TLS.ClientKeyPath, "Path to client TLS key that will be used to connect to Vault")
	cmd.PersistentFlags().StringVar(&cfg.SDS.Vault.TLS.ServerName, "vault-tls-server-name", cfg.SDS.Vault.TLS.ServerName, "If set, it is used to set the SNI host when connecting via TLS")
	cmd.PersistentFlags().BoolVar(&cfg.SDS.Vault.TLS.SkipVerify, "vault-tls-skip-verify", cfg.SDS.Vault.TLS.SkipVerify, "Disables TLS verification")
	return cmd
}

func isDpTokenRequired(cfg kuma_dp.Config, catalog catalog.Catalog) bool {
	return cfg.SDS.Type == kuma_dp.SdsCp && catalog.Apis.DataplaneToken.Enabled()
}
