/*
Copyright 2015-2019 Gravitational, Inc.

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

package services

import (
	"encoding/json"

	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
)

// ReverseTunnel is SSH reverse tunnel established between a local Proxy
// and a remote Proxy. It helps to bypass firewall restrictions, so local
// clusters don't need to have the cluster involved
type ReverseTunnel interface {
	// Resource provides common methods for resource objects
	Resource
	// GetClusterName returns name of the cluster
	GetClusterName() string
	// SetClusterName sets cluster name
	SetClusterName(name string)
	// GetType gets the type of ReverseTunnel.
	GetType() TunnelType
	// SetType sets the type of ReverseTunnel.
	SetType(TunnelType)
	// GetDialAddrs returns list of dial addresses for this cluster
	GetDialAddrs() []string
	// Check checks tunnel for errors
	Check() error
	// CheckAndSetDefaults checks and set default values for any missing fields.
	CheckAndSetDefaults() error
}

// NewReverseTunnel returns new version of reverse tunnel
func NewReverseTunnel(clusterName string, dialAddrs []string) ReverseTunnel {
	return &ReverseTunnelV2{
		Kind:    KindReverseTunnel,
		Version: V2,
		Metadata: Metadata{
			Name:      clusterName,
			Namespace: defaults.Namespace,
		},
		Spec: ReverseTunnelSpecV2{
			ClusterName: clusterName,
			DialAddrs:   dialAddrs,
		},
	}
}

// UnmarshalReverseTunnel unmarshals reverse tunnel from JSON or YAML,
// sets defaults and checks the schema
func UnmarshalReverseTunnel(data []byte, opts ...MarshalOption) (ReverseTunnel, error) {
	if len(data) == 0 {
		return nil, trace.BadParameter("missing tunnel data")
	}
	var h ResourceHeader
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	cfg, err := collectOptions(opts)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	switch h.Version {
	case "":
		var r ReverseTunnelV1
		err := json.Unmarshal(data, &r)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		v2 := r.V2()
		if cfg.ID != 0 {
			v2.SetResourceID(cfg.ID)
		}
		return r.V2(), nil
	case V2:
		var r ReverseTunnelV2
		if cfg.SkipValidation {
			if err := utils.FastUnmarshal(data, &r); err != nil {
				return nil, trace.BadParameter(err.Error())
			}
		} else {
			if err := utils.UnmarshalWithSchema(GetReverseTunnelSchema(), &r, data); err != nil {
				return nil, trace.BadParameter(err.Error())
			}
		}
		if err := r.CheckAndSetDefaults(); err != nil {
			return nil, trace.Wrap(err)
		}
		if cfg.ID != 0 {
			r.SetResourceID(cfg.ID)
		}
		if !cfg.Expires.IsZero() {
			r.SetExpiry(cfg.Expires)
		}
		return &r, nil
	}
	return nil, trace.BadParameter("reverse tunnel version %v is not supported", h.Version)
}

var tunnelMarshaler ReverseTunnelMarshaler = &TeleportTunnelMarshaler{}

func SetReerseTunnelMarshaler(m ReverseTunnelMarshaler) {
	marshalerMutex.Lock()
	defer marshalerMutex.Unlock()
	tunnelMarshaler = m
}

func GetReverseTunnelMarshaler() ReverseTunnelMarshaler {
	marshalerMutex.Lock()
	defer marshalerMutex.Unlock()
	return tunnelMarshaler
}

// ReverseTunnelMarshaler implements marshal/unmarshal of reverse tunnel implementations
type ReverseTunnelMarshaler interface {
	// UnmarshalReverseTunnel unmarshals reverse tunnel from binary representation
	UnmarshalReverseTunnel(bytes []byte, opts ...MarshalOption) (ReverseTunnel, error)
	// MarshalReverseTunnel marshals reverse tunnel to binary representation
	MarshalReverseTunnel(ReverseTunnel, ...MarshalOption) ([]byte, error)
}

type TeleportTunnelMarshaler struct{}

// UnmarshalReverseTunnel unmarshals reverse tunnel from JSON or YAML
func (*TeleportTunnelMarshaler) UnmarshalReverseTunnel(bytes []byte, opts ...MarshalOption) (ReverseTunnel, error) {
	return UnmarshalReverseTunnel(bytes, opts...)
}

// MarshalRole marshalls role into JSON
func (*TeleportTunnelMarshaler) MarshalReverseTunnel(rt ReverseTunnel, opts ...MarshalOption) ([]byte, error) {
	cfg, err := collectOptions(opts)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	type tunv1 interface {
		V1() *ReverseTunnelV1
	}
	type tunv2 interface {
		V2() *ReverseTunnelV2
	}
	version := cfg.GetVersion()
	switch version {
	case V1:
		v, ok := rt.(tunv1)
		if !ok {
			return nil, trace.BadParameter("don't know how to marshal %v", V1)
		}
		return json.Marshal(v.V1())
	case V2:
		v, ok := rt.(tunv2)
		if !ok {
			return nil, trace.BadParameter("don't know how to marshal %v", V2)
		}
		v2 := v.V2()
		if !cfg.PreserveResourceID {
			// avoid modifying the original object
			// to prevent unexpected data races
			copy := *v2
			copy.SetResourceID(0)
			v2 = &copy
		}
		return utils.FastMarshal(v2)
	default:
		return nil, trace.BadParameter("version %v is not supported", version)
	}
}
