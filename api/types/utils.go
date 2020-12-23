package types

import (
	"net/url"
	"strings"
	"time"

	"github.com/gravitational/trace"
)

// These functions are copied over from utils.go to phase out importing of utils.
// TODO: replace originals in /lib/utils with aliases to copies.

// CheckParseAddr takes a string and returns true if it can be parsed into a utils.NetAddr
func CheckParseAddr(a string) error {
	if a == "" {
		return trace.BadParameter("missing parameter address")
	}
	if !strings.Contains(a, "://") {
		return nil
	}
	u, err := url.Parse(a)
	if err != nil {
		return trace.BadParameter("failed to parse %q: %v", a, err)
	}
	switch u.Scheme {
	case "tcp", "unix", "http", "https":
		return nil
	default:
		return trace.BadParameter("'%v': unsupported scheme: '%v'", a, u.Scheme)
	}
}

// CopyStrings makes a deep copy of the passed in string slice and returns
// the copy.
func CopyStrings(in []string) []string {
	if in == nil {
		return nil
	}

	out := make([]string, len(in))
	copy(out, in)

	return out
}

// StringSlicesEqual returns true if string slices equal
func StringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// UTC converts time to UTC timezone
func UTC(t *time.Time) {
	if t == nil {
		return
	}

	if t.IsZero() {
		// to fix issue with timezones for tests
		*t = time.Time{}
		return
	}
	*t = t.UTC()
}
