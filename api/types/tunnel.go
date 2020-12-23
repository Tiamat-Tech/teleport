/*
Copyright 2020 Gravitational, Inc.

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

package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/gravitational/teleport/api/constants"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// GetVersion returns resource version
func (r *ReverseTunnelV2) GetVersion() string {
	return r.Version
}

// GetKind returns resource kind
func (r *ReverseTunnelV2) GetKind() string {
	return r.Kind
}

// GetSubKind returns resource sub kind
func (r *ReverseTunnelV2) GetSubKind() string {
	return r.SubKind
}

// SetSubKind sets resource subkind
func (r *ReverseTunnelV2) SetSubKind(s string) {
	r.SubKind = s
}

// GetResourceID returns resource ID
func (r *ReverseTunnelV2) GetResourceID() int64 {
	return r.Metadata.ID
}

// SetResourceID sets resource ID
func (r *ReverseTunnelV2) SetResourceID(id int64) {
	r.Metadata.ID = id
}

// GetMetadata returns object metadata
func (r *ReverseTunnelV2) GetMetadata() Metadata {
	return r.Metadata
}

// SetExpiry sets expiry time for the object
func (r *ReverseTunnelV2) SetExpiry(expires time.Time) {
	r.Metadata.SetExpiry(expires)
}

// Expiry returns object expiry setting
func (r *ReverseTunnelV2) Expiry() time.Time {
	return r.Metadata.Expiry()
}

// SetTTL sets Expires header using realtime clock
func (r *ReverseTunnelV2) SetTTL(clock clockwork.Clock, ttl time.Duration) {
	r.Metadata.SetTTL(clock, ttl)
}

// GetName returns the name of the User
func (r *ReverseTunnelV2) GetName() string {
	return r.Metadata.Name
}

// SetName sets the name of the User
func (r *ReverseTunnelV2) SetName(e string) {
	r.Metadata.Name = e
}

// CheckAndSetDefaults checks and sets defaults
func (r *ReverseTunnelV2) CheckAndSetDefaults() error {
	err := r.Metadata.CheckAndSetDefaults()
	if err != nil {
		return trace.Wrap(err)
	}

	err = r.Check()
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// SetClusterName sets name of a cluster
func (r *ReverseTunnelV2) SetClusterName(name string) {
	r.Spec.ClusterName = name
}

// GetClusterName returns name of the cluster
func (r *ReverseTunnelV2) GetClusterName() string {
	return r.Spec.ClusterName
}

// GetType gets the type of ReverseTunnel.
func (r *ReverseTunnelV2) GetType() TunnelType {
	if string(r.Spec.Type) == "" {
		return ProxyTunnel
	}
	return r.Spec.Type
}

// SetType sets the type of ReverseTunnel.
func (r *ReverseTunnelV2) SetType(tt TunnelType) {
	r.Spec.Type = tt
}

// GetDialAddrs returns list of dial addresses for this cluster
func (r *ReverseTunnelV2) GetDialAddrs() []string {
	return r.Spec.DialAddrs
}

// V2 returns V2 version of the resource
func (r *ReverseTunnelV2) V2() *ReverseTunnelV2 {
	return r
}

// V1 returns V1 version of the resource
func (r *ReverseTunnelV2) V1() *ReverseTunnelV1 {
	return &ReverseTunnelV1{
		DomainName: r.Spec.ClusterName,
		DialAddrs:  r.Spec.DialAddrs,
	}
}

// ReverseTunnelV1 is V1 version of reverse tunnel
type ReverseTunnelV1 struct {
	// DomainName is a domain name of remote cluster we are connecting to
	DomainName string `json:"domain_name"`
	// DialAddrs is a list of remote address to establish a connection to
	// it's always SSH over TCP
	DialAddrs []string `json:"dial_addrs"`
}

// V1 returns V1 version of the resource
func (r *ReverseTunnelV1) V1() *ReverseTunnelV1 {
	return r
}

// V2 returns V2 version of reverse tunnel
func (r *ReverseTunnelV1) V2() *ReverseTunnelV2 {
	return &ReverseTunnelV2{
		Kind:    constants.KindReverseTunnel,
		Version: constants.V2,
		Metadata: Metadata{
			Name:      r.DomainName,
			Namespace: constants.Namespace,
		},
		Spec: ReverseTunnelSpecV2{
			ClusterName: r.DomainName,
			Type:        ProxyTunnel,
			DialAddrs:   r.DialAddrs,
		},
	}
}

// Check returns nil if all parameters are good, error otherwise
func (r *ReverseTunnelV2) Check() error {
	if r.Version == "" {
		return trace.BadParameter("missing reverse tunnel version")
	}
	if strings.TrimSpace(r.Spec.ClusterName) == "" {
		return trace.BadParameter("Reverse tunnel validation error: empty cluster name")
	}

	if len(r.Spec.DialAddrs) == 0 {
		return trace.BadParameter("Invalid dial address for reverse tunnel '%v'", r.Spec.ClusterName)
	}

	for _, addr := range r.Spec.DialAddrs {
		if err := CheckParseAddr(addr); err != nil {
			return trace.Wrap(err)
		}
	}

	return nil
}

// GetReverseTunnelSchema returns role schema with optionally injected
// schema for extensions
func GetReverseTunnelSchema() string {
	return fmt.Sprintf(V2SchemaTemplate, MetadataSchema, ReverseTunnelSpecV2Schema, DefaultDefinitions)
}

// ReverseTunnelSpecV2Schema is JSON schema for reverse tunnel spec
const ReverseTunnelSpecV2Schema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["cluster_name", "dial_addrs"],
  "properties": {
    "cluster_name": {"type": "string"},
    "type": {"type": "string"},
    "dial_addrs": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  }
}`

const (
	// NodeTunnel is a tunnel where the node connects to the proxy (dial back).
	NodeTunnel TunnelType = "node"

	// ProxyTunnel is a tunnel where a proxy connects to the proxy (trusted cluster).
	ProxyTunnel TunnelType = "proxy"

	// AppTunnel is a tunnel where the application proxy dials back to the proxy.
	AppTunnel TunnelType = "app"

	// KubeTunnel is a tunnel where the kubernetes service dials back to the proxy.
	KubeTunnel TunnelType = "kube"
)

// TunnelType is the type of tunnel. Either node or proxy.
type TunnelType string
