/*
Copyright 2021 Kuma authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kube_core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kumahq/kuma/pkg/plugins/runtime/k8s/metadata"
	"github.com/kumahq/kuma/pkg/transparentproxy/kubernetes"
)

var _ = Describe("kubernetes", func() {
	type testCaseKumactl struct {
		transparentProxyV1 bool
		pod                *kube_core.Pod
		commandLine        []string
	}

	DescribeTable("should generate kumactl command line", func(given testCaseKumactl) {
		podRedirect, err := kubernetes.NewPodRedirectForPod(given.transparentProxyV1, given.pod)
		Expect(err).ToNot(HaveOccurred())

		commandLine := podRedirect.AsKumactlCommandLine()
		Expect(commandLine).To(Equal(given.commandLine))
	},
		Entry("should generate", testCaseKumactl{
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaBuiltinDNS:                                 metadata.AnnotationEnabled,
						metadata.KumaBuiltinDNSPort:                             "25053",
						metadata.KumaTrafficExcludeOutboundPorts:                "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:  "25100",
						metadata.KumaTrafficExcludeInboundPorts:                 "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:   "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6: "25206",
						metadata.KumaSidecarUID:                                 "12345",
						metadata.KumaTrafficExcludeOutboundUDPPortsForUIDs:      "11001:1;11002:2",
						metadata.KumaTrafficExcludeOutboundTCPPortsForUIDs:      "11003:3",
						metadata.KumaTrafficExcludeOutboundPortsForUIDs:         "0;12",
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "true",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
				"--exclude-outbound-ports-for-uids", "0",
				"--exclude-outbound-ports-for-uids", "12",
				"--exclude-outbound-ports-for-uids", "tcp:11003:3",
				"--exclude-outbound-ports-for-uids", "udp:11001:1",
				"--exclude-outbound-ports-for-uids", "udp:11002:2",
				"--redirect-all-dns-traffic",
				"--redirect-dns-port", "25053",
			},
		}),
		Entry("should generate with deprecated dns annotation", testCaseKumactl{
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaTrafficExcludeOutboundPorts:                "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:  "25100",
						metadata.KumaTrafficExcludeInboundPorts:                 "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:   "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6: "25206",
						metadata.KumaSidecarUID:                                 "12345",
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "true",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
			},
		}),

		Entry("should generate no builtin DNS", testCaseKumactl{
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaTrafficExcludeOutboundPorts:                "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:  "25100",
						metadata.KumaTrafficExcludeInboundPorts:                 "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:   "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6: "25206",
						metadata.KumaSidecarUID:                                 "12345",
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "true",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
			},
		}),
		Entry("should generate engine v1", testCaseKumactl{
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaTrafficExcludeOutboundPorts:                "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:  "25100",
						metadata.KumaTrafficExcludeInboundPorts:                 "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:   "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6: "25206",
						metadata.KumaSidecarUID:                                 "12345",
						metadata.KumaTransparentProxyingEngineV1:                metadata.AnnotationEnabled,
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "true",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
				"--use-transparent-proxy-engine-v1",
			},
		}),
		Entry("should generate engine v1 if enabled even without annotation", testCaseKumactl{
			transparentProxyV1: true,
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaTrafficExcludeOutboundPorts:                "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:  "25100",
						metadata.KumaTrafficExcludeInboundPorts:                 "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:   "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6: "25206",
						metadata.KumaSidecarUID:                                 "12345",
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "true",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
				"--use-transparent-proxy-engine-v1",
			},
		}),
		Entry("should generate for Gateway", testCaseKumactl{
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaBuiltinDNS:                                 metadata.AnnotationEnabled,
						metadata.KumaBuiltinDNSPort:                             "25053",
						metadata.KumaTrafficExcludeOutboundPorts:                "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:  "25100",
						metadata.KumaGatewayAnnotation:                          metadata.AnnotationEnabled,
						metadata.KumaTrafficExcludeInboundPorts:                 "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:   "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6: "25206",
						metadata.KumaSidecarUID:                                 "12345",
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "false",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
				"--redirect-all-dns-traffic",
				"--redirect-dns-port", "25053",
			},
		}),
		Entry("should generate for ebpf transparent proxy", testCaseKumactl{
			pod: &kube_core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						metadata.KumaBuiltinDNS:                                  metadata.AnnotationEnabled,
						metadata.KumaBuiltinDNSPort:                              "25053",
						metadata.KumaTrafficExcludeOutboundPorts:                 "11000",
						metadata.KumaTransparentProxyingOutboundPortAnnotation:   "25100",
						metadata.KumaGatewayAnnotation:                           metadata.AnnotationEnabled,
						metadata.KumaTrafficExcludeInboundPorts:                  "12000",
						metadata.KumaTransparentProxyingInboundPortAnnotation:    "25204",
						metadata.KumaTransparentProxyingInboundPortAnnotationV6:  "25206",
						metadata.KumaSidecarUID:                                  "12345",
						metadata.KumaTransparentProxyingEbpf:                     metadata.AnnotationEnabled,
						metadata.KumaTransparentProxyingEbpfInstanceIPEnvVarName: "FOO_BAR",
						metadata.KumaTransparentProxyingEbpfBPFFSPath:            "/baz/bar/foo",
						metadata.KumaTransparentProxyingEbpfCgroupPath:           "/foo/bar/baz",
						metadata.KumaTransparentProxyingEbpfProgramsSourcePath:   "/foo",
					},
				},
			},
			commandLine: []string{
				"--redirect-outbound-port", "25100",
				"--redirect-inbound=" + "false",
				"--redirect-inbound-port", "25204",
				"--redirect-inbound-port-v6", "25206",
				"--kuma-dp-uid", "12345",
				"--exclude-inbound-ports", "12000",
				"--exclude-outbound-ports", "11000",
				"--verbose",
				"--skip-resolv-conf",
				"--redirect-all-dns-traffic",
				"--redirect-dns-port", "25053",
				"--ebpf-enabled",
				"--ebpf-instance-ip", "$(FOO_BAR)",
				"--ebpf-bpffs-path", "/baz/bar/foo",
				"--ebpf-cgroup-path", "/foo/bar/baz",
				"--ebpf-programs-source-path", "/foo",
			},
		}),
	)
})
