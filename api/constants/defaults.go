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

package constants

import (
	"time"

	"github.com/gravitational/teleport"
)

const (
	// Namespace is default namespace
	Namespace = "default"

	// ServerKeepAliveTTL is a period between server keep alives,
	// when servers announce only presence withough sending full data
	ServerKeepAliveTTL = 60 * time.Second

	// DefaultDialTimeout is a default TCP dial timeout we set for our
	// connection attempts
	DefaultDialTimeout = 30 * time.Second

	// KeepAliveCountMax is the number of keep-alive messages that can be sent
	// without receiving a response from the client before the client is
	// disconnected. The max count mirrors ClientAliveCountMax of sshd.
	KeepAliveCountMax = 3

	// MaxCertDuration limits maximum duration of validity of issued cert
	MaxCertDuration = 30 * time.Hour
)

// EnhancedEvents returns the default list of enhanced events.
func EnhancedEvents() []string {
	return []string{
		teleport.EnhancedRecordingCommand,
		teleport.EnhancedRecordingNetwork,
	}
}
