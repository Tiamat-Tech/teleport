package types

import "github.com/gravitational/teleport/lib/backend"

// Event represents an event that happened in the backend
type Event struct {
	// Type is the event type
	Type backend.OpType
	// Resource is a modified or deleted resource
	// in case of deleted resources, only resource header
	// will be provided
	Resource Resource
}
