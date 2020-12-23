package services

import (
	"github.com/gravitational/teleport/api/constants"
	"github.com/gravitational/teleport/api/types"
)

// The following types are implemented in /api/types, and imported/wrapped here.
// The new structs are used to wrap the imported types with additional methods.
// The other types are basic imports and can be removed if their references are updated.

type AccessRequest = types.AccessRequest
type AccessRequestV3 = types.AccessRequestV3
type AccessRequestSpecV3 = types.AccessRequestSpecV3
type AccessRequestFilter = types.AccessRequestFilter
type AccessRequestConditions struct{ types.AccessRequestConditions }
type RequestState = types.RequestState

type App = types.App

type Duration = types.Duration

type Event = types.Event

type PluginData = types.PluginData
type PluginDataV3 = types.PluginDataV3
type PluginDataSpecV3 = types.PluginDataSpecV3
type PluginDataFilter = types.PluginDataFilter
type PluginDataEntry = types.PluginDataEntry
type PluginDataUpdateParams = types.PluginDataUpdateParams

type ProvisionTokenV1 struct{ types.ProvisionTokenV1 }
type ProvisionTokenV2 struct{ types.ProvisionTokenV2 }
type ProvisionTokenSpecV2 = types.ProvisionTokenSpecV2

type RemoteClusterV3 struct{ types.RemoteClusterV3 }

type Resource = types.Resource
type ResourceHeader struct{ types.ResourceHeader }
type Metadata = types.Metadata

type Role = types.Role
type RoleV3 = types.RoleV3
type RoleSpecV3 = types.RoleSpecV3
type RoleConditions = types.RoleConditions
type RoleConditionType = types.RoleConditionType
type RoleOptions = types.RoleOptions
type Rule = types.Rule
type Labels = types.Labels

type WebSessionV2 struct{ types.WebSessionV2 }
type WebSessionSpecV2 = types.WebSessionSpecV2

// Some functions and variables also need to be imported from the types package
var (
	NewRule       = types.NewRule
	NewBoolOption = types.NewBoolOption

	MaxDuration           = types.MaxDuration
	NewDuration           = types.NewDuration
	IsValidLabelKey       = types.IsValidLabelKey
	MetadataSchema        = types.MetadataSchema
	CopyRulesSlice        = types.CopyRulesSlice
	RequestState_NONE     = types.RequestState_NONE
	RequestState_PENDING  = types.RequestState_PENDING
	RequestState_APPROVED = types.RequestState_APPROVED
	RequestState_DENIED   = types.RequestState_DENIED
)

// The following Constants are imported from api to simplify
// refactoring. These could be removed and their references updated.
const (
	DefaultAPIGroup               = constants.DefaultAPIGroup
	ActionRead                    = constants.ActionRead
	ActionWrite                   = constants.ActionWrite
	Wildcard                      = constants.Wildcard
	KindNamespace                 = constants.KindNamespace
	KindUser                      = constants.KindUser
	KindKeyPair                   = constants.KindKeyPair
	KindHostCert                  = constants.KindHostCert
	KindJWT                       = constants.KindJWT
	KindLicense                   = constants.KindLicense
	KindRole                      = constants.KindRole
	KindAccessRequest             = constants.KindAccessRequest
	KindPluginData                = constants.KindPluginData
	KindOIDC                      = constants.KindOIDC
	KindSAML                      = constants.KindSAML
	KindGithub                    = constants.KindGithub
	KindOIDCRequest               = constants.KindOIDCRequest
	KindSAMLRequest               = constants.KindSAMLRequest
	KindGithubRequest             = constants.KindGithubRequest
	KindSession                   = constants.KindSession
	KindSSHSession                = constants.KindSSHSession
	KindWebSession                = constants.KindWebSession
	KindAppSession                = constants.KindAppSession
	KindEvent                     = constants.KindEvent
	KindAuthServer                = constants.KindAuthServer
	KindProxy                     = constants.KindProxy
	KindNode                      = constants.KindNode
	KindAppServer                 = constants.KindAppServer
	KindToken                     = constants.KindToken
	KindCertAuthority             = constants.KindCertAuthority
	KindReverseTunnel             = constants.KindReverseTunnel
	KindOIDCConnector             = constants.KindOIDCConnector
	KindSAMLConnector             = constants.KindSAMLConnector
	KindGithubConnector           = constants.KindGithubConnector
	KindConnectors                = constants.KindConnectors
	KindClusterAuthPreference     = constants.KindClusterAuthPreference
	MetaNameClusterAuthPreference = constants.MetaNameClusterAuthPreference
	KindClusterConfig             = constants.KindClusterConfig
	KindSemaphore                 = constants.KindSemaphore
	MetaNameClusterConfig         = constants.MetaNameClusterConfig
	KindClusterName               = constants.KindClusterName
	MetaNameClusterName           = constants.MetaNameClusterName
	KindStaticTokens              = constants.KindStaticTokens
	MetaNameStaticTokens          = constants.MetaNameStaticTokens
	KindTrustedCluster            = constants.KindTrustedCluster
	KindAuthConnector             = constants.KindAuthConnector
	KindTunnelConnection          = constants.KindTunnelConnection
	KindRemoteCluster             = constants.KindRemoteCluster
	KindResetPasswordToken        = constants.KindResetPasswordToken
	KindResetPasswordTokenSecrets = constants.KindResetPasswordTokenSecrets
	KindIdentity                  = constants.KindIdentity
	KindState                     = constants.KindState
	KindKubeService               = constants.KindKubeService
	V3                            = constants.V3
	V2                            = constants.V2
	V1                            = constants.V1
	VerbList                      = constants.VerbList
	VerbCreate                    = constants.VerbCreate
	VerbRead                      = constants.VerbRead
	VerbReadNoSecrets             = constants.VerbReadNoSecrets
	VerbUpdate                    = constants.VerbUpdate
	VerbDelete                    = constants.VerbDelete
	VerbRotate                    = constants.VerbRotate
)
