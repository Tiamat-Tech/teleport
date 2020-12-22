package types

import (
	fmt "fmt"
	time "time"

	"github.com/gravitational/teleport/lib/utils"
	"github.com/jonboulle/clockwork"
)

// RoleConditionType specifies if it's an allow rule (true) or deny rule (false).
type RoleConditionType bool

const (
	// Allow is the set of conditions that allow access.
	Allow RoleConditionType = true
	// Deny is the set of conditions that prevent access.
	Deny RoleConditionType = false
)

// Role contains a set of permissions or settings
type Role interface {
	// Resource provides common resource methods.
	Resource
	// CheckAndSetDefaults checks and set default values for any missing fields.
	CheckAndSetDefaults() error
	// Equals returns true if the roles are equal. Roles are equal if options and
	// conditions match.
	Equals(other Role) bool
	// ApplyTraits applies the passed in traits to any variables within the role
	// and returns itself.
	ApplyTraits(map[string][]string) Role

	// GetOptions gets role options.
	GetOptions() RoleOptions
	// SetOptions sets role options
	SetOptions(opt RoleOptions)

	// GetLogins gets *nix system logins for allow or deny condition.
	GetLogins(RoleConditionType) []string
	// SetLogins sets *nix system logins for allow or deny condition.
	SetLogins(RoleConditionType, []string)

	// GetNamespaces gets a list of namespaces this role is allowed or denied access to.
	GetNamespaces(RoleConditionType) []string
	// GetNamespaces sets a list of namespaces this role is allowed or denied access to.
	SetNamespaces(RoleConditionType, []string)

	// GetNodeLabels gets the map of node labels this role is allowed or denied access to.
	GetNodeLabels(RoleConditionType) Labels
	// SetNodeLabels sets the map of node labels this role is allowed or denied access to.
	SetNodeLabels(RoleConditionType, Labels)

	// GetAppLabels gets the map of app labels this role is allowed or denied access to.
	GetAppLabels(RoleConditionType) Labels
	// SetAppLabels sets the map of app labels this role is allowed or denied access to.
	SetAppLabels(RoleConditionType, Labels)

	// GetClusterLabels gets the map of cluster labels this role is allowed or denied access to.
	GetClusterLabels(RoleConditionType) Labels
	// SetClusterLabels sets the map of cluster labels this role is allowed or denied access to.
	SetClusterLabels(RoleConditionType, Labels)

	// GetKubernetesLabels gets the map of kubernetes labels this role is
	// allowed or denied access to.
	GetKubernetesLabels(RoleConditionType) Labels
	// SetKubernetesLabels sets the map of kubernetes labels this role is
	// allowed or denied access to.
	SetKubernetesLabels(RoleConditionType, Labels)

	// GetRules gets all allow or deny rules.
	GetRules(rct RoleConditionType) []Rule
	// SetRules sets an allow or deny rule.
	SetRules(rct RoleConditionType, rules []Rule)

	// GetKubeGroups returns kubernetes groups
	GetKubeGroups(RoleConditionType) []string
	// SetKubeGroups sets kubernetes groups for allow or deny condition.
	SetKubeGroups(RoleConditionType, []string)

	// GetKubeUsers returns kubernetes users to impersonate
	GetKubeUsers(RoleConditionType) []string
	// SetKubeUsers sets kubernetes users to impersonate for allow or deny condition.
	SetKubeUsers(RoleConditionType, []string)

	// GetAccessRequestConditions gets allow/deny conditions for access requests.
	GetAccessRequestConditions(RoleConditionType) AccessRequestConditions
	// SetAccessRequestConditions sets allow/deny conditions for access requests.
	SetAccessRequestConditions(RoleConditionType, AccessRequestConditions)
}

// Equals returns true if the roles are equal. Roles are equal if options,
// namespaces, logins, labels, and conditions match.
func (r *RoleV3) Equals(other Role) bool {
	if !r.GetOptions().Equals(other.GetOptions()) {
		return false
	}

	for _, condition := range []RoleConditionType{Allow, Deny} {
		if !utils.StringSlicesEqual(r.GetLogins(condition), other.GetLogins(condition)) {
			return false
		}
		if !utils.StringSlicesEqual(r.GetNamespaces(condition), other.GetNamespaces(condition)) {
			return false
		}
		if !r.GetNodeLabels(condition).Equals(other.GetNodeLabels(condition)) {
			return false
		}
		if !r.GetAppLabels(condition).Equals(other.GetAppLabels(condition)) {
			return false
		}
		if !RuleSlicesEqual(r.GetRules(condition), other.GetRules(condition)) {
			return false
		}
		if !r.GetClusterLabels(condition).Equals(other.GetClusterLabels(condition)) {
			return false
		}
		if !r.GetKubernetesLabels(condition).Equals(other.GetKubernetesLabels(condition)) {
			return false
		}
	}

	return true
}

// GetVersion returns resource version
func (r *RoleV3) GetVersion() string {
	return r.Version
}

// GetKind returns resource kind
func (r *RoleV3) GetKind() string {
	return r.Kind
}

// GetSubKind returns resource sub kind
func (r *RoleV3) GetSubKind() string {
	return r.SubKind
}

// SetSubKind sets resource subkind
func (r *RoleV3) SetSubKind(s string) {
	r.SubKind = s
}

// GetResourceID returns resource ID
func (r *RoleV3) GetResourceID() int64 {
	return r.Metadata.ID
}

// SetResourceID sets resource ID
func (r *RoleV3) SetResourceID(id int64) {
	r.Metadata.ID = id
}

// TODO: ApplyTraits would go here

// SetExpiry sets expiry time for the object.
func (r *RoleV3) SetExpiry(expires time.Time) {
	r.Metadata.SetExpiry(expires)
}

// Expiry returns the expiry time for the object.
func (r *RoleV3) Expiry() time.Time {
	return r.Metadata.Expiry()
}

// SetTTL sets TTL header using realtime clock.
func (r *RoleV3) SetTTL(clock clockwork.Clock, ttl time.Duration) {
	r.Metadata.SetTTL(clock, ttl)
}

// SetName sets the role name and is a shortcut for SetMetadata().Name.
func (r *RoleV3) SetName(s string) {
	r.Metadata.Name = s
}

// GetName gets the role name and is a shortcut for GetMetadata().Name.
func (r *RoleV3) GetName() string {
	return r.Metadata.Name
}

// GetMetadata returns role metadata.
func (r *RoleV3) GetMetadata() Metadata {
	return r.Metadata
}

// GetOptions gets role options.
func (r *RoleV3) GetOptions() RoleOptions {
	return r.Spec.Options
}

// SetOptions sets role options.
func (r *RoleV3) SetOptions(options RoleOptions) {
	r.Spec.Options = options
}

// GetLogins gets system logins for allow or deny condition.
func (r *RoleV3) GetLogins(rct RoleConditionType) []string {
	if rct == Allow {
		return r.Spec.Allow.Logins
	}
	return r.Spec.Deny.Logins
}

// SetLogins sets system logins for allow or deny condition.
func (r *RoleV3) SetLogins(rct RoleConditionType, logins []string) {
	lcopy := CopyStrings(logins)

	if rct == Allow {
		r.Spec.Allow.Logins = lcopy
	} else {
		r.Spec.Deny.Logins = lcopy
	}
}

// GetKubeGroups returns kubernetes groups
func (r *RoleV3) GetKubeGroups(rct RoleConditionType) []string {
	if rct == Allow {
		return r.Spec.Allow.KubeGroups
	}
	return r.Spec.Deny.KubeGroups
}

// SetKubeGroups sets kubernetes groups for allow or deny condition.
func (r *RoleV3) SetKubeGroups(rct RoleConditionType, groups []string) {
	lcopy := CopyStrings(groups)

	if rct == Allow {
		r.Spec.Allow.KubeGroups = lcopy
	} else {
		r.Spec.Deny.KubeGroups = lcopy
	}
}

// GetKubeUsers returns kubernetes users
func (r *RoleV3) GetKubeUsers(rct RoleConditionType) []string {
	if rct == Allow {
		return r.Spec.Allow.KubeUsers
	}
	return r.Spec.Deny.KubeUsers
}

// SetKubeUsers sets kubernetes user for allow or deny condition.
func (r *RoleV3) SetKubeUsers(rct RoleConditionType, users []string) {
	lcopy := CopyStrings(users)

	if rct == Allow {
		r.Spec.Allow.KubeUsers = lcopy
	} else {
		r.Spec.Deny.KubeUsers = lcopy
	}
}

// GetAccessRequestConditions gets conditions for access requests.
func (r *RoleV3) GetAccessRequestConditions(rct RoleConditionType) AccessRequestConditions {
	cond := r.Spec.Deny.Request
	if rct == Allow {
		cond = r.Spec.Allow.Request
	}
	if cond == nil {
		return AccessRequestConditions{}
	}
	return *cond
}

// SetAccessRequestConditions sets allow/deny conditions for access requests.
func (r *RoleV3) SetAccessRequestConditions(rct RoleConditionType, cond AccessRequestConditions) {
	if rct == Allow {
		r.Spec.Allow.Request = &cond
	} else {
		r.Spec.Deny.Request = &cond
	}
}

// GetNamespaces gets a list of namespaces this role is allowed or denied access to.
func (r *RoleV3) GetNamespaces(rct RoleConditionType) []string {
	if rct == Allow {
		return r.Spec.Allow.Namespaces
	}
	return r.Spec.Deny.Namespaces
}

// SetNamespaces sets a list of namespaces this role is allowed or denied access to.
func (r *RoleV3) SetNamespaces(rct RoleConditionType, namespaces []string) {
	ncopy := CopyStrings(namespaces)

	if rct == Allow {
		r.Spec.Allow.Namespaces = ncopy
	} else {
		r.Spec.Deny.Namespaces = ncopy
	}
}

// GetNodeLabels gets the map of node labels this role is allowed or denied access to.
func (r *RoleV3) GetNodeLabels(rct RoleConditionType) Labels {
	if rct == Allow {
		return r.Spec.Allow.NodeLabels
	}
	return r.Spec.Deny.NodeLabels
}

// SetNodeLabels sets the map of node labels this role is allowed or denied access to.
func (r *RoleV3) SetNodeLabels(rct RoleConditionType, labels Labels) {
	if rct == Allow {
		r.Spec.Allow.NodeLabels = labels.Clone()
	} else {
		r.Spec.Deny.NodeLabels = labels.Clone()
	}
}

// GetAppLabels gets the map of app labels this role is allowed or denied access to.
func (r *RoleV3) GetAppLabels(rct RoleConditionType) Labels {
	if rct == Allow {
		return r.Spec.Allow.AppLabels
	}
	return r.Spec.Deny.AppLabels
}

// SetAppLabels sets the map of node labels this role is allowed or denied access to.
func (r *RoleV3) SetAppLabels(rct RoleConditionType, labels Labels) {
	if rct == Allow {
		r.Spec.Allow.AppLabels = labels.Clone()
	} else {
		r.Spec.Deny.AppLabels = labels.Clone()
	}
}

// GetClusterLabels gets the map of cluster labels this role is allowed or denied access to.
func (r *RoleV3) GetClusterLabels(rct RoleConditionType) Labels {
	if rct == Allow {
		return r.Spec.Allow.ClusterLabels
	}
	return r.Spec.Deny.ClusterLabels
}

// SetClusterLabels sets the map of cluster labels this role is allowed or denied access to.
func (r *RoleV3) SetClusterLabels(rct RoleConditionType, labels Labels) {
	if rct == Allow {
		r.Spec.Allow.ClusterLabels = labels.Clone()
	} else {
		r.Spec.Deny.ClusterLabels = labels.Clone()
	}
}

// GetKubernetesLabels gets the map of app labels this role is allowed or denied access to.
func (r *RoleV3) GetKubernetesLabels(rct RoleConditionType) Labels {
	if rct == Allow {
		return r.Spec.Allow.KubernetesLabels
	}
	return r.Spec.Deny.KubernetesLabels
}

// SetKubernetesLabels sets the map of node labels this role is allowed or denied access to.
func (r *RoleV3) SetKubernetesLabels(rct RoleConditionType, labels Labels) {
	if rct == Allow {
		r.Spec.Allow.KubernetesLabels = labels.Clone()
	} else {
		r.Spec.Deny.KubernetesLabels = labels.Clone()
	}
}

// GetRules gets all allow or deny rules.
func (r *RoleV3) GetRules(rct RoleConditionType) []Rule {
	if rct == Allow {
		return r.Spec.Allow.Rules
	}
	return r.Spec.Deny.Rules
}

// SetRules sets an allow or deny rule.
func (r *RoleV3) SetRules(rct RoleConditionType, in []Rule) {
	rcopy := CopyRulesSlice(in)

	if rct == Allow {
		r.Spec.Allow.Rules = rcopy
	} else {
		r.Spec.Deny.Rules = rcopy
	}
}

// TODO: CheckAndSetDefaults

// String returns the human readable representation of a role.
func (r *RoleV3) String() string {
	return fmt.Sprintf("Role(Name=%v,Options=%v,Allow=%+v,Deny=%+v)",
		r.GetName(), r.Spec.Options, r.Spec.Allow, r.Spec.Deny)
}

// Equals checks if all the key/values in the RoleOptions map match.
func (o RoleOptions) Equals(other RoleOptions) bool {
	return (o.ForwardAgent.Value() == other.ForwardAgent.Value() &&
		o.MaxSessionTTL.Value() == other.MaxSessionTTL.Value() &&
		BoolDefaultTrue(o.PortForwarding) == BoolDefaultTrue(other.PortForwarding) &&
		o.CertificateFormat == other.CertificateFormat &&
		o.ClientIdleTimeout.Value() == other.ClientIdleTimeout.Value() &&
		o.DisconnectExpiredCert.Value() == other.DisconnectExpiredCert.Value() &&
		utils.StringSlicesEqual(o.BPF, other.BPF))
}

// Equals returns true if the role conditions (logins, namespaces, labels,
// and rules) are equal and false if they are not.
func (r *RoleConditions) Equals(o RoleConditions) bool {
	if !utils.StringSlicesEqual(r.Logins, o.Logins) {
		return false
	}
	if !utils.StringSlicesEqual(r.Namespaces, o.Namespaces) {
		return false
	}
	if !r.NodeLabels.Equals(o.NodeLabels) {
		return false
	}
	if !r.AppLabels.Equals(o.AppLabels) {
		return false
	}
	if !r.KubernetesLabels.Equals(o.KubernetesLabels) {
		return false
	}
	if len(r.Rules) != len(o.Rules) {
		return false
	}
	for i := range r.Rules {
		if !r.Rules[i].Equals(o.Rules[i]) {
			return false
		}
	}
	return true
}

// BoolOption is a wrapper around bool
// that can take multiple values:
// * true, false and non-set (when pointer is nil)
// and can marshal itself to protobuf equivalent BoolValue
type BoolOption struct {
	// Value is a value of the option
	Value bool
}

// BoolDefaultTrue returns true if v is not set (pointer is nil)
// otherwise returns real boolean value
func BoolDefaultTrue(v *BoolOption) bool {
	if v == nil {
		return true
	}
	return v.Value
}

// CopyRulesSlice copies input slice of Rules and returns the copy
func CopyRulesSlice(in []Rule) []Rule {
	out := make([]Rule, len(in))
	copy(out, in)
	return out
}

// RuleSlicesEqual returns true if two rule slices are equal
func RuleSlicesEqual(a, b []Rule) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

// Labels is a wrapper around map
// that can marshal and unmarshal itself
// from scalar and list values
type Labels map[string]utils.Strings

// Clone returns non-shallow copy of the labels set
func (l Labels) Clone() Labels {
	if l == nil {
		return nil
	}
	out := make(Labels, len(l))
	for key, vals := range l {
		cvals := make([]string, len(vals))
		copy(cvals, vals)
		out[key] = cvals
	}
	return out
}
