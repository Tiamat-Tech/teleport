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

package services

import (
	"fmt"

	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
)

type PluginDataMarshaler interface {
	MarshalPluginData(req PluginData, opts ...MarshalOption) ([]byte, error)
	UnmarshalPluginData(bytes []byte, opts ...MarshalOption) (PluginData, error)
}

type pluginDataMarshaler struct{}

func (m *pluginDataMarshaler) MarshalPluginData(data PluginData, opts ...MarshalOption) ([]byte, error) {
	cfg, err := collectOptions(opts)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	switch r := data.(type) {
	case *PluginDataV3:
		if !cfg.PreserveResourceID {
			// avoid modifying the original object
			// to prevent unexpected data races
			cp := *r
			cp.SetResourceID(0)
			r = &cp
		}
		return utils.FastMarshal(r)
	default:
		return nil, trace.BadParameter("unrecognized plugin data type: %T", data)
	}
}

func (m *pluginDataMarshaler) UnmarshalPluginData(raw []byte, opts ...MarshalOption) (PluginData, error) {
	cfg, err := collectOptions(opts)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	var data PluginDataV3
	if cfg.SkipValidation {
		if err := utils.FastUnmarshal(raw, &data); err != nil {
			return nil, trace.Wrap(err)
		}
	} else {
		if err := utils.UnmarshalWithSchema(GetPluginDataSchema(), &data, raw); err != nil {
			return nil, trace.Wrap(err)
		}
	}
	if err := data.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}
	if cfg.ID != 0 {
		data.SetResourceID(cfg.ID)
	}
	if !cfg.Expires.IsZero() {
		data.SetExpiry(cfg.Expires)
	}
	return &data, nil
}

var pluginDataMarshalerInstance PluginDataMarshaler = &pluginDataMarshaler{}

func GetPluginDataMarshaler() PluginDataMarshaler {
	marshalerMutex.Lock()
	defer marshalerMutex.Unlock()
	return pluginDataMarshalerInstance
}

const PluginDataSpecSchema = `{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"entries": { "type":"object" }
	}
}`

func GetPluginDataSchema() string {
	return fmt.Sprintf(V2SchemaTemplate, MetadataSchema, PluginDataSpecSchema, DefaultDefinitions)
}
