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

package api

import (
	"github.com/gravitational/teleport/api/auth"
	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/trace"
)

func EventToGRPC(in services.Event) (*auth.Event, error) {
	eventType, err := eventTypeToGRPC(in.Type)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	out := auth.Event{
		Type: eventType,
	}
	if in.Type == backend.OpInit {
		return &out, nil
	}
	switch r := in.Resource.(type) {
	case *services.ResourceHeader:
		out.Resource = &auth.Event_ResourceHeader{
			ResourceHeader: r,
		}
	case *services.CertAuthorityV2:
		out.Resource = &auth.Event_CertAuthority{
			CertAuthority: r,
		}
	case *services.StaticTokensV2:
		out.Resource = &auth.Event_StaticTokens{
			StaticTokens: r,
		}
	case *services.ProvisionTokenV2:
		out.Resource = &auth.Event_ProvisionToken{
			ProvisionToken: r,
		}
	case *services.ClusterNameV2:
		out.Resource = &auth.Event_ClusterName{
			ClusterName: r,
		}
	case *services.ClusterConfigV3:
		out.Resource = &auth.Event_ClusterConfig{
			ClusterConfig: r,
		}
	case *services.UserV2:
		out.Resource = &auth.Event_User{
			User: r,
		}
	case *services.RoleV3:
		out.Resource = &auth.Event_Role{
			Role: r,
		}
	case *services.Namespace:
		out.Resource = &auth.Event_Namespace{
			Namespace: r,
		}
	case *services.ServerV2:
		out.Resource = &auth.Event_Server{
			Server: r,
		}
	case *services.ReverseTunnelV2:
		out.Resource = &auth.Event_ReverseTunnel{
			ReverseTunnel: r,
		}
	case *services.TunnelConnectionV2:
		out.Resource = &auth.Event_TunnelConnection{
			TunnelConnection: r,
		}
	case *services.AccessRequestV3:
		out.Resource = &auth.Event_AccessRequest{
			AccessRequest: r,
		}
	case *services.WebSessionV2:
		if r.GetSubKind() != services.KindAppSession {
			return nil, trace.BadParameter("only %v supported", services.KindAppSession)
		}
		out.Resource = &auth.Event_AppSession{
			AppSession: r,
		}
	case *services.RemoteClusterV3:
		out.Resource = &auth.Event_RemoteCluster{
			RemoteCluster: r,
		}
	default:
		return nil, trace.BadParameter("resource type %T is not supported", in.Resource)
	}
	return &out, nil
}

func eventTypeToGRPC(in backend.OpType) (auth.Operation, error) {
	switch in {
	case backend.OpInit:
		return auth.Operation_INIT, nil
	case backend.OpPut:
		return auth.Operation_PUT, nil
	case backend.OpDelete:
		return auth.Operation_DELETE, nil
	default:
		return -1, trace.BadParameter("event type %v is not supported", in)
	}
}

func EventFromGRPC(in auth.Event) (*services.Event, error) {
	eventType, err := eventTypeFromGRPC(in.Type)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	out := services.Event{
		Type: eventType,
	}
	if eventType == backend.OpInit {
		return &out, nil
	}
	if r := in.GetResourceHeader(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetCertAuthority(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetStaticTokens(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetProvisionToken(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetClusterName(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetClusterConfig(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetUser(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetRole(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetNamespace(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetServer(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetReverseTunnel(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetTunnelConnection(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetAccessRequest(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetAppSession(); r != nil {
		out.Resource = r
		return &out, nil
	} else if r := in.GetRemoteCluster(); r != nil {
		out.Resource = r
		return &out, nil
	} else {
		return nil, trace.BadParameter("received unsupported resource %T", in.Resource)
	}
}

func eventTypeFromGRPC(in auth.Operation) (backend.OpType, error) {
	switch in {
	case auth.Operation_INIT:
		return backend.OpInit, nil
	case auth.Operation_PUT:
		return backend.OpPut, nil
	case auth.Operation_DELETE:
		return backend.OpDelete, nil
	default:
		return backend.OpInvalid, trace.BadParameter("unsupported operation type: %v", in)
	}
}
