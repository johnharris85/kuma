package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/kumahq/kuma/tools/resource-gen/genutils"
)

// CustomResourceTemplate for creating a Kubernetes CRD to wrap a Kuma resource.
var CustomResourceTemplate = template.Must(template.New("custom-resource").Parse(`
// Generated by tools/resource-gen
// Run "make generate" to update this file.

{{ $tk := "` + "`" + `" }}

// nolint:whitespace
package {{.PolicyVersion}}

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"


	"github.com/kumahq/kuma/pkg/plugins/policies/{{.KumactlSingular}}/api/{{.PolicyVersion}}"
	"github.com/kumahq/kuma/pkg/plugins/resources/k8s/native/pkg/model"
	{{- if not .SkipRegistration }}
	"github.com/kumahq/kuma/pkg/plugins/resources/k8s/native/pkg/registry"
	{{- end }}
)

{{- if not .SkipKubernetesWrappers }}

// +kubebuilder:object:root=true
{{- if .ScopeNamespace }}
// +kubebuilder:resource:categories=kuma,scope=Namespaced
{{- else }}
// +kubebuilder:resource:categories=kuma,scope=Cluster
{{- end}}
{{- if .StorageVersion }}
// +kubebuilder:storageversion
{{- end}}
type {{.ResourceType}} struct {
	metav1.TypeMeta   {{ $tk }}json:",inline"{{ $tk }}
	metav1.ObjectMeta {{ $tk }}json:"metadata,omitempty"{{ $tk }}

    // Mesh is the name of the Kuma mesh this resource belongs to.
	// It may be omitted for cluster-scoped resources.
	//
    // +kubebuilder:validation:Optional
	Mesh string {{ $tk }}json:"mesh,omitempty"{{ $tk }}

{{- if eq .ResourceType "DataplaneInsight" }}
	// Status is the status the Kuma resource.
    // +kubebuilder:validation:Optional
	Status   *apiextensionsv1.JSON {{ $tk }}json:"status,omitempty"{{ $tk }}
{{- else}}
	// Spec is the specification of the Kuma {{ .ProtoType }} resource.
    // +kubebuilder:validation:Optional
	Spec   *{{.PolicyVersion}}.{{.ProtoType}} {{ $tk }}json:"spec,omitempty"{{ $tk }}
{{- end}}
}

// +kubebuilder:object:root=true
{{- if .ScopeNamespace }}
// +kubebuilder:resource:scope=Cluster
{{- else }}
// +kubebuilder:resource:scope=Namespaced
{{- end}}
type {{.ResourceType}}List struct {
	metav1.TypeMeta {{ $tk }}json:",inline"{{ $tk }}
	metav1.ListMeta {{ $tk }}json:"metadata,omitempty"{{ $tk }}
	Items           []{{.ResourceType}} {{ $tk }}json:"items"{{ $tk }}
}

{{- if not .SkipRegistration}}
func init() {
	SchemeBuilder.Register(&{{.ResourceType}}{}, &{{.ResourceType}}List{})
}
{{- end}}

func (cb *{{.ResourceType}}) GetObjectMeta() *metav1.ObjectMeta {
	return &cb.ObjectMeta
}

func (cb *{{.ResourceType}}) SetObjectMeta(m *metav1.ObjectMeta) {
	cb.ObjectMeta = *m
}

func (cb *{{.ResourceType}}) GetMesh() string {
	return cb.Mesh
}

func (cb *{{.ResourceType}}) SetMesh(mesh string) {
	cb.Mesh = mesh
}

func (cb *{{.ResourceType}}) GetSpec() proto.Message {
{{- if eq .ResourceType "DataplaneInsight" }}
	return cb.Status
{{- else}}
	return  cb.Spec
{{- end}}
}

func (cb *{{.ResourceType}}) SetSpec(spec proto.Message) {
	if spec == nil {
{{- if eq .ResourceType "DataplaneInsight" }}
		cb.Status = nil
{{ else }}
		cb.Spec = nil
{{- end }}
		return
	}

	if _, ok := spec.(*{{.PolicyVersion}}.{{.ProtoType}}); !ok {
		panic(fmt.Sprintf("unexpected protobuf message type %T", spec))
	}

{{ if eq .ResourceType "DataplaneInsight" }}
	cb.Status = spec.(*{{.PolicyVersion}}.{{.ProtoType}})
{{ else}}
	cb.Spec = spec.(*{{.PolicyVersion}}.{{.ProtoType}})
{{- end}}
}

func (cb *{{.ResourceType}}) Scope() model.Scope {
{{- if .ScopeNamespace }}
	return model.ScopeNamespace
{{- else }}
	return model.ScopeCluster
{{- end }}
}

func (l *{{.ResourceType}}List) GetItems() []model.KubernetesObject {
	result := make([]model.KubernetesObject, len(l.Items))
	for i := range l.Items {
		result[i] = &l.Items[i]
	}
	return result
}

{{if not .SkipRegistration}}
func init() {
	registry.RegisterObjectType(&{{.PolicyVersion}}.{{.ProtoType}}{}, &{{.ResourceType}}{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       "{{.ResourceType}}",
		},
	})
	registry.RegisterListType(&{{.PolicyVersion}}.{{.ProtoType}}{}, &{{.ResourceType}}List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       "{{.ResourceType}}List",
		},
	})
}
{{- end }} {{/* .SkipRegistration */}}
{{- end }} {{/* .SkipKubernetesWrappers */}}
`))

var GroupVersionInfoTemplate = template.Must(template.New("groupversion-info").Parse(`
// Package {{.PolicyVersion}} contains API Schema definitions for the mesh {{.PolicyVersion}} API group
// +groupName=kuma.io
package {{.PolicyVersion}}

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "kuma.io", Version: "{{.PolicyVersion}}"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
`))

func generateCRD(
	p *protogen.Plugin,
	file *protogen.File,
) error {
	var infos []genutils.ResourceInfo
	for _, msg := range file.Messages {
		infos = append(infos, genutils.ToResourceInfo(msg.Desc))
	}

	if len(infos) > 1 {
		return errors.Errorf("only one Kuma resource per proto file is allowed")
	}

	info := infos[0]

	outBuf := bytes.Buffer{}
	if err := CustomResourceTemplate.Execute(&outBuf, struct {
		genutils.ResourceInfo
		PolicyVersion string
	}{
		ResourceInfo:  info,
		PolicyVersion: string(file.GoPackageName),
	}); err != nil {
		return err
	}

	out, err := format.Source(outBuf.Bytes())
	if err != nil {
		return err
	}

	typesGenerator := p.NewGeneratedFile(fmt.Sprintf("k8s/%s/zz_generated.types.go", string(file.GoPackageName)), file.GoImportPath)
	if _, err := typesGenerator.Write(out); err != nil {
		return err
	}

	gviOutBuf := bytes.Buffer{}
	if err := GroupVersionInfoTemplate.Execute(&gviOutBuf, struct {
		PolicyVersion string
	}{
		PolicyVersion: string(file.GoPackageName),
	}); err != nil {
		return err
	}

	gvi, err := format.Source(gviOutBuf.Bytes())
	if err != nil {
		return err
	}

	gviGenerator := p.NewGeneratedFile(fmt.Sprintf("k8s/%s/groupversion_info.go", string(file.GoPackageName)), file.GoImportPath)
	if _, err := gviGenerator.Write(gvi); err != nil {
		return err
	}
	return nil
}
