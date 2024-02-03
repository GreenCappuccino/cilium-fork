// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package v2alpha1

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cilium/cilium/pkg/option"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories={cilium},singular="ciliumenvoyhttpfilter",path="ciliumenvoyhttpfilters",scope="Namespaced",shortName={cehf}
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",description="The age of the identity",name="Age",type=date
// +kubebuilder:storageversion

type CiliumEnvoyHTTPFilter struct {
	// +k8s:openapi-gen=false
	// +deepequal-gen=false
	metav1.TypeMeta `json:",inline"`
	// +k8s:openapi-gen=false
	// +deepequal-gen=false
	metav1.ObjectMeta `json:"metadata"`

	// +k8s:openapi-gen=false
	Spec CiliumEnvoyHTTPFilterSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +deepequal-gen=false

// CiliumEnvoyHTTPFilterList is a list of CiliumEnvoyHTTPFilter objects.
type CiliumEnvoyHTTPFilterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items is a list of CiliumEnvoyConfig.
	Items []CiliumEnvoyHTTPFilter `json:"items"`
}

type CiliumEnvoyHTTPFilterSpec struct {
	// HTTPFilters is a list of HTTPFilter to be inserted in the HTTP connection manager filter chain
	//
	// +kubebuilder:validation:Optional
	HTTPFilters []*HTTPFilter `json:"httpFilters,omitempty"`
}

// HTTPFilter is an Envoy extensions.filters.network.http_connection_manager.v3.HttpFilter
//
// +kubebuilder:validation:XValidation:message="HTTPFilter must have exactly 1 of typedConfig or configDiscovery",rule="(has(self.typedConfig) || has(self.configDiscovery)) && !(has(self.typedConfig) && has(self.configDiscovery))"
type HTTPFilter struct {
	// Name is the name of the filter configuration.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// TypedConfig is filter specific configuration which depends on the filter being instantiated.
	//
	// +kubebuilder:validation:Optional
	TypedConfig TypedConfig `json:"typedConfig,omitempty"`
	// ConfigDiscovery is a configuration source specifier for an extension configuration discovery service.
	//
	// Warning: Note that this is not validated extensively for now.
	//
	// +kubebuilder:validation:Optional
	ConfigDiscovery ExtensionConfigSource `json:"configDiscovery,omitempty"`
	// IsOptional, if set to true, allows clients that do not support this filter to ignore the filter but otherwise accept the config. Otherwise, clients that do not support this filter must reject the config.
	//
	// +kubebuilder:validation:Optional
	IsOptional bool `json:"isOptional"`
	// Disabled, if set to true, makes the filter disabled by default, and must be explicitly enabled by setting per filter configuration in the route configuration.
	//
	// +kubebuilder:validation:Optional
	Disabled bool `json:"disabled"`
}

// TypedConfig is a stand-in for Envoy's HTTP Filter typed_config
//
// +kubebuilder:pruning:PreserveUnknownFields
type TypedConfig struct {
	*anypb.Any `json:"-"`
}

// DeepCopyInto deep copies 'in' into 'out'.
func (in *TypedConfig) DeepCopyInto(out *TypedConfig) {
	out.Any, _ = proto.Clone(in.Any).(*anypb.Any)
}

// DeepEqual returns 'true' if 'a' and 'b' are equal.
func (a *TypedConfig) DeepEqual(b *TypedConfig) bool {
	return proto.Equal(a.Any, b.Any)
}

// MarshalJSON ensures that the unstructured object produces proper
// JSON when passed to Go's standard JSON library.
func (u *TypedConfig) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(u.Any)
}

// UnmarshalJSON ensures that the unstructured object properly decodes
// JSON when passed to Go's standard JSON library.
func (u *TypedConfig) UnmarshalJSON(b []byte) (err error) {
	// TypedConfig resources are not validated in K8s, recover from possible panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("HTTP Filter JSON decoding paniced: %v", r)
		}
	}()
	u.Any = &anypb.Any{}
	err = protojson.Unmarshal(b, u.Any)
	if err != nil {
		var buf bytes.Buffer
		json.Indent(&buf, b, "", "\t")
		log.Warningf("Ignoring invalid CiliumEnvoyHTTPFilter JSON (%s): %s",
			err, buf.String())
	} else if option.Config.Debug {
		log.Debugf("HTTP Filter unmarshaled TypedConfig Resource: %v", prototext.Format(u.Any))
	}
	return nil
}

// ExtensionConfigSource is a stand-in for Envoy's config.core.v3.ExtensionConfigSource
//
// +kubebuilder:pruning:PreserveUnknownFields
type ExtensionConfigSource struct {
	*anypb.Any `json:"-"`
}

// DeepCopyInto deep copies 'in' into 'out'.
func (in *ExtensionConfigSource) DeepCopyInto(out *ExtensionConfigSource) {
	out.Any, _ = proto.Clone(in.Any).(*anypb.Any)
}

// DeepEqual returns 'true' if 'a' and 'b' are equal.
func (a *ExtensionConfigSource) DeepEqual(b *ExtensionConfigSource) bool {
	return proto.Equal(a.Any, b.Any)
}

// MarshalJSON ensures that the unstructured object produces proper
// JSON when passed to Go's standard JSON library.
func (u *ExtensionConfigSource) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(u.Any)
}

// UnmarshalJSON ensures that the unstructured object properly decodes
// JSON when passed to Go's standard JSON library.
func (u *ExtensionConfigSource) UnmarshalJSON(b []byte) (err error) {
	// ExtensionConfigSource resources are not validated in K8s, recover from possible panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("HTTP Filter JSON decoding paniced: %v", r)
		}
	}()
	u.Any = &anypb.Any{}
	err = protojson.Unmarshal(b, u.Any)
	if err != nil {
		var buf bytes.Buffer
		json.Indent(&buf, b, "", "\t")
		log.Warningf("Ignoring invalid CiliumEnvoyHTTPFilter JSON (%s): %s",
			err, buf.String())
	} else if option.Config.Debug {
		log.Debugf("HTTP Filter unmarshaled ExtensionConfigSource Resource: %v", prototext.Format(u.Any))
	}
	return nil
}
