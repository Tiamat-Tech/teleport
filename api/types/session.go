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
	"time"

	"github.com/gravitational/teleport/api/constants"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

func (ws *WebSessionV2) GetKind() string {
	return ws.Kind
}

func (ws *WebSessionV2) GetSubKind() string {
	return ws.SubKind
}

func (ws *WebSessionV2) SetSubKind(subKind string) {
	ws.SubKind = subKind
}

func (ws *WebSessionV2) GetVersion() string {
	return ws.Version
}

func (ws *WebSessionV2) GetName() string {
	return ws.Metadata.Name
}

func (ws *WebSessionV2) SetName(name string) {
	ws.Metadata.Name = name
}

func (ws *WebSessionV2) Expiry() time.Time {
	return ws.Metadata.Expiry()
}

func (ws *WebSessionV2) SetExpiry(expiry time.Time) {
	ws.Metadata.SetExpiry(expiry)
}

func (ws *WebSessionV2) SetTTL(clock clockwork.Clock, ttl time.Duration) {
	ws.Metadata.SetTTL(clock, ttl)
}

func (ws *WebSessionV2) GetMetadata() Metadata {
	return ws.Metadata
}

func (ws *WebSessionV2) GetResourceID() int64 {
	return ws.Metadata.GetID()
}

func (ws *WebSessionV2) SetResourceID(id int64) {
	ws.Metadata.SetID(id)
}

// CheckAndSetDefaults checks and set default values for any missing fields.
func (ws *WebSessionV2) CheckAndSetDefaults() error {
	err := ws.Metadata.CheckAndSetDefaults()
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// String returns string representation of the session.
func (ws *WebSessionV2) String() string {
	return fmt.Sprintf("WebSession(kind=%v,name=%v,id=%v)", ws.GetKind(), ws.GetUser(), ws.GetName())
}

// SetUser sets user associated with this session
func (ws *WebSessionV2) SetUser(u string) {
	ws.Spec.User = u
}

// GetUser returns the user this session is associated with
func (ws *WebSessionV2) GetUser() string {
	return ws.Spec.User
}

// GetShortName returns visible short name used in logging
func (ws *WebSessionV2) GetShortName() string {
	if len(ws.Metadata.Name) < 4 {
		return "<undefined>"
	}
	return ws.Metadata.Name[:4]
}

// GetTLSCert returns PEM encoded TLS certificate associated with session
func (ws *WebSessionV2) GetTLSCert() []byte {
	return ws.Spec.TLSCert
}

// GetPub is returns public certificate signed by auth server
func (ws *WebSessionV2) GetPub() []byte {
	return ws.Spec.Pub
}

// GetPriv returns private OpenSSH key used to auth with SSH nodes
func (ws *WebSessionV2) GetPriv() []byte {
	return ws.Spec.Priv
}

// SetPriv sets private key
func (ws *WebSessionV2) SetPriv(priv []byte) {
	ws.Spec.Priv = priv
}

// BearerToken is a special bearer token used for additional
// bearer authentication
func (ws *WebSessionV2) GetBearerToken() string {
	return ws.Spec.BearerToken
}

// SetBearerTokenExpiryTime sets bearer token expiry time
func (ws *WebSessionV2) SetBearerTokenExpiryTime(tm time.Time) {
	ws.Spec.BearerTokenExpires = tm
}

// SetExpiryTime sets session expiry time
func (ws *WebSessionV2) SetExpiryTime(tm time.Time) {
	ws.Spec.Expires = tm
}

// GetBearerTokenExpiryTime - absolute time when token expires
func (ws *WebSessionV2) GetBearerTokenExpiryTime() time.Time {
	return ws.Spec.BearerTokenExpires
}

// GetExpiryTime - absolute time when web session expires
func (ws *WebSessionV2) GetExpiryTime() time.Time {
	return ws.Spec.Expires
}

// V2 returns V2 version of the resource
func (ws *WebSessionV2) V2() *WebSessionV2 {
	return ws
}

// V1 returns V1 version of the object
func (ws *WebSessionV2) V1() *WebSessionV1 {
	return &WebSessionV1{
		ID:          ws.Metadata.Name,
		Priv:        ws.Spec.Priv,
		Pub:         ws.Spec.Pub,
		BearerToken: ws.Spec.BearerToken,
		Expires:     ws.Spec.Expires,
	}
}

// WebSession stores key and value used to authenticate with SSH
// nodes on behalf of user
type WebSessionV1 struct {
	// ID is session ID
	ID string `json:"id"`
	// User is a user this web session is associated with
	User string `json:"user"`
	// Pub is a public certificate signed by auth server
	Pub []byte `json:"pub"`
	// Priv is a private OpenSSH key used to auth with SSH nodes
	Priv []byte `json:"priv,omitempty"`
	// BearerToken is a special bearer token used for additional
	// bearer authentication
	BearerToken string `json:"bearer_token"`
	// Expires - absolute time when token expires
	Expires time.Time `json:"expires"`
}

// V1 returns V1 version of the resource
func (ws *WebSessionV1) V1() *WebSessionV1 {
	return ws
}

// V2 returns V2 version of the resource
func (ws *WebSessionV1) V2() *WebSessionV2 {
	return &WebSessionV2{
		Kind:    constants.KindWebSession,
		Version: constants.V2,
		Metadata: Metadata{
			Name:      ws.ID,
			Namespace: constants.Namespace,
		},
		Spec: WebSessionSpecV2{
			Pub:                ws.Pub,
			Priv:               ws.Priv,
			BearerToken:        ws.BearerToken,
			Expires:            ws.Expires,
			BearerTokenExpires: ws.Expires,
		},
	}
}

// SetName sets session name
func (ws *WebSessionV1) SetName(name string) {
	ws.ID = name
}

// SetUser sets user associated with this session
func (ws *WebSessionV1) SetUser(u string) {
	ws.User = u
}

// GetUser returns the user this session is associated with
func (ws *WebSessionV1) GetUser() string {
	return ws.User
}

// GetShortName returns visible short name used in logging
func (ws *WebSessionV1) GetShortName() string {
	if len(ws.ID) < 4 {
		return "<undefined>"
	}
	return ws.ID[:4]
}

// GetName returns session name
func (ws *WebSessionV1) GetName() string {
	return ws.ID
}

// GetPub is returns public certificate signed by auth server
func (ws *WebSessionV1) GetPub() []byte {
	return ws.Pub
}

// GetPriv returns private OpenSSH key used to auth with SSH nodes
func (ws *WebSessionV1) GetPriv() []byte {
	return ws.Priv
}

// BearerToken is a special bearer token used for additional
// bearer authentication
func (ws *WebSessionV1) GetBearerToken() string {
	return ws.BearerToken
}

// Expires - absolute time when token expires
func (ws *WebSessionV1) GetExpiryTime() time.Time {
	return ws.Expires
}

// SetExpiryTime sets session expiry time
func (ws *WebSessionV1) SetExpiryTime(tm time.Time) {
	ws.Expires = tm
}

// GetBearerRoken - absolute time when token expires
func (ws *WebSessionV1) GetBearerTokenExpiryTime() time.Time {
	return ws.Expires
}

// SetBearerTokenExpiryTime sets session expiry time
func (ws *WebSessionV1) SetBearerTokenExpiryTime(tm time.Time) {
	ws.Expires = tm
}

// GetWebSessionSchema returns JSON Schema for web session
func GetWebSessionSchema() string {
	return GetWebSessionSchemaWithExtensions("")
}

// GetWebSessionSchemaWithExtensions returns JSON Schema for web session with user-supplied extensions
func GetWebSessionSchemaWithExtensions(extension string) string {
	return fmt.Sprintf(V2SchemaTemplate, MetadataSchema, fmt.Sprintf(WebSessionSpecV2Schema, extension), DefaultDefinitions)
}

// WebSessionSpecV2Schema is JSON schema for cert authority V2
const WebSessionSpecV2Schema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["pub", "bearer_token", "bearer_token_expires", "expires", "user"],
  "properties": {
    "user": {"type": "string"},
    "pub": {"type": "string"},
    "priv": {"type": "string"},
    "tls_cert": {"type": "string"},
    "bearer_token": {"type": "string"},
    "bearer_token_expires": {"type": "string"},
    "expires": {"type": "string"}%v
  }
}`
