package api

import "context"

type contextKey string

const (
	// ContextDelegator is a delegator for access requests set in the context
	// of the request
	ContextDelegator contextKey = "delegator"
)

// GetDelegator attempts to load the context value AccessRequestDelegator,
// returning the empty string if no value was found.
func GetDelegator(ctx context.Context) string {
	delegator, ok := ctx.Value(ContextDelegator).(string)
	if !ok {
		return ""
	}
	return delegator
}

// WithDelegator creates a child context with the AccessRequestDelegator
// value set.  Optionally used by AuthServer.SetAccessRequestState to log
// a delegating identity.
func WithDelegator(ctx context.Context, delegator string) context.Context {
	return context.WithValue(ctx, ContextDelegator, delegator)
}
