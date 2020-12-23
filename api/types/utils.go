package types

// These functions are copied over from utils.go to phase out importing of utils.
// TODO: replace originals in /lib/utils with aliases to copies.

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
