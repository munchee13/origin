package registry

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	"github.com/openshift/origin/pkg/auth/api"
	"github.com/openshift/origin/pkg/oauth/registry/clientauthorization"
	"github.com/openshift/origin/pkg/oauth/scope"
)

type ClientAuthorizationGrantChecker struct {
	registry clientauthorization.Registry
}

func NewClientAuthorizationGrantChecker(registry clientauthorization.Registry) *ClientAuthorizationGrantChecker {
	return &ClientAuthorizationGrantChecker{registry}
}

func (c *ClientAuthorizationGrantChecker) HasAuthorizedClient(user user.Info, grant *api.Grant) (approved bool, err error) {
	id := c.registry.ClientAuthorizationName(user.GetName(), grant.Client.GetId())
	authorization, err := c.registry.GetClientAuthorization(id)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(authorization.UserUID) != 0 && authorization.UserUID != user.GetUID() {
		return false, fmt.Errorf("user %s UID %s does not match stored client authorization value for UID %s", user.GetName(), user.GetUID(), authorization.UserUID)
	}
	// TODO: improve this to allow the scope implementation to determine overlap
	if !scope.Covers(authorization.Scopes, scope.Split(grant.Scope)) {
		return false, nil
	}
	return true, nil
}
