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
	"regexp"
	"time"

	"github.com/gravitational/teleport/api/constants"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// Resource represents common properties for all resources.
type Resource interface {
	// GetKind returns resource kind
	GetKind() string
	// GetSubKind returns resource subkind
	GetSubKind() string
	// SetSubKind sets resource subkind
	SetSubKind(string)
	// GetVersion returns resource version
	GetVersion() string
	// GetName returns the name of the resource
	GetName() string
	// SetName sets the name of the resource
	SetName(string)
	// Expiry returns object expiry setting
	Expiry() time.Time
	// SetExpiry sets object expiry
	SetExpiry(time.Time)
	// SetTTL sets Expires header using current clock
	SetTTL(clock clockwork.Clock, ttl time.Duration)
	// GetMetadata returns object metadata
	GetMetadata() Metadata
	// GetResourceID returns resource ID
	GetResourceID() int64
	// SetResourceID sets resource ID
	SetResourceID(int64)
}

// GetVersion returns resource version
func (h *ResourceHeader) GetVersion() string {
	return h.Version
}

// GetResourceID returns resource ID
func (h *ResourceHeader) GetResourceID() int64 {
	return h.Metadata.ID
}

// SetResourceID sets resource ID
func (h *ResourceHeader) SetResourceID(id int64) {
	h.Metadata.ID = id
}

// GetName returns the name of the resource
func (h *ResourceHeader) GetName() string {
	return h.Metadata.Name
}

// SetName sets the name of the resource
func (h *ResourceHeader) SetName(v string) {
	h.Metadata.SetName(v)
}

// Expiry returns object expiry setting
func (h *ResourceHeader) Expiry() time.Time {
	return h.Metadata.Expiry()
}

// SetExpiry sets object expiry
func (h *ResourceHeader) SetExpiry(t time.Time) {
	h.Metadata.SetExpiry(t)
}

// SetTTL sets Expires header using current clock
func (h *ResourceHeader) SetTTL(clock clockwork.Clock, ttl time.Duration) {
	h.Metadata.SetTTL(clock, ttl)
}

// GetMetadata returns object metadata
func (h *ResourceHeader) GetMetadata() Metadata {
	return h.Metadata
}

// GetKind returns resource kind
func (h *ResourceHeader) GetKind() string {
	return h.Kind
}

// GetSubKind returns resource subkind
func (h *ResourceHeader) GetSubKind() string {
	return h.SubKind
}

// SetSubKind sets resource subkind
func (h *ResourceHeader) SetSubKind(s string) {
	h.SubKind = s
}

// GetID returns resource ID
func (m *Metadata) GetID() int64 {
	return m.ID
}

// SetID sets resource ID
func (m *Metadata) SetID(id int64) {
	m.ID = id
}

// GetMetadata returns object metadata
func (m *Metadata) GetMetadata() Metadata {
	return *m
}

// GetName returns the name of the resource
func (m *Metadata) GetName() string {
	return m.Name
}

// SetName sets the name of the resource
func (m *Metadata) SetName(name string) {
	m.Name = name
}

// SetExpiry sets expiry time for the object
func (m *Metadata) SetExpiry(expires time.Time) {
	m.Expires = &expires
}

// Expiry returns object expiry setting.
func (m *Metadata) Expiry() time.Time {
	if m.Expires == nil {
		return time.Time{}
	}
	return *m.Expires
}

// SetTTL sets Expires header using realtime clock
func (m *Metadata) SetTTL(clock clockwork.Clock, ttl time.Duration) {
	expireTime := clock.Now().UTC().Add(ttl)
	m.Expires = &expireTime
}

// CheckAndSetDefaults checks validity of all parameters and sets defaults
func (m *Metadata) CheckAndSetDefaults() error {
	if m.Name == "" {
		return trace.BadParameter("missing parameter Name")
	}
	if m.Namespace == "" {
		m.Namespace = constants.Namespace
	}

	// adjust expires time to utc if it's set
	if m.Expires != nil {
		UTC(m.Expires)
	}

	for key := range m.Labels {
		if !IsValidLabelKey(key) {
			return trace.BadParameter("invalid label key: %q", key)
		}
	}

	return nil
}

const labelPattern = `^[a-zA-Z/.0-9_*-]+$`

var validLabelKey = regexp.MustCompile(labelPattern)

// IsValidLabelKey checks if the supplied string matches the
// label key regexp.
func IsValidLabelKey(s string) bool {
	return validLabelKey.MatchString(s)
}

// MetadataSchema is a schema for resource metadata
var MetadataSchema = fmt.Sprintf(baseMetadataSchema, labelPattern)

// DefaultDefinitions the default list of JSON schema definitions which is none.
const DefaultDefinitions = ``

const baseMetadataSchema = `{
  "type": "object",
  "additionalProperties": false,
  "default": {},
  "required": ["name"],
  "properties": {
    "name": {"type": "string"},
    "namespace": {"type": "string", "default": "default"},
    "description": {"type": "string"},
    "expires": {"type": "string"},
    "id": {"type": "integer"},
    "labels": {
      "type": "object",
      "additionalProperties": false,
      "patternProperties": {
         "%s":  { "type": "string" }
      }
    }
  }
}`

// V2SchemaTemplate is a template JSON Schema for V2 style objects
const V2SchemaTemplate = `{
	"type": "object",
	"additionalProperties": false,
	"required": ["kind", "spec", "metadata", "version"],
	"properties": {
	  "kind": {"type": "string"},
	  "sub_kind": {"type": "string"},
	  "version": {"type": "string", "default": "v2"},
	  "metadata": %v,
	  "spec": %v
	}%v
  }`
