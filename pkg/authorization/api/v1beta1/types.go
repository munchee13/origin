package v1beta1

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api/v1beta3"
	kruntime "github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	kutil "github.com/GoogleCloudPlatform/kubernetes/pkg/util"
)

// Authorization is calculated against
// 1. all deny RoleBinding PolicyRules in the master namespace - short circuit on match
// 2. all allow RoleBinding PolicyRules in the master namespace - short circuit on match
// 3. all deny RoleBinding PolicyRules in the namespace - short circuit on match
// 4. all allow RoleBinding PolicyRules in the namespace - short circuit on match
// 5. deny by default

// PolicyRule holds information that describes a policy rule, but does not contain information
// about who the rule applies to or which namespace the rule applies to.
type PolicyRule struct {
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule.  VerbAll represents all kinds.
	Verbs []string `json:"verbs"`
	// AttributeRestrictions will vary depending on what the Authorizer/AuthorizationAttributeBuilder pair supports.
	// If the Authorizer does not recognize how to handle the AttributeRestrictions, the Authorizer should report an error.
	AttributeRestrictions kruntime.RawExtension `json:"attributeRestrictions"`
	// ResourceKinds is a list of resources this rule applies to.  ResourceAll represents all resources.
	// DEPRECATED
	ResourceKinds []string `json:"resourceKinds,omitempty"`
	// Resources is a list of resources this rule applies to.  ResourceAll represents all resources.
	Resources []string `json:"resources"`
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	ResourceNames []string `json:"resourceNames,omitempty"`
}

// Role is a logical grouping of PolicyRules that can be referenced as a unit by RoleBindings.
type Role struct {
	kapi.TypeMeta   `json:",inline"`
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// Rules holds all the PolicyRules for this Role
	Rules []PolicyRule `json:"rules"`
}

// RoleBinding references a Role, but not contain it.  It can reference any Role in the same namespace or in the global namespace.
// It adds who information via Users and Groups and namespace information by which namespace it exists in.  RoleBindings in a given
// namespace only have effect in that namespace (excepting the master namespace which has power in all namespaces).
type RoleBinding struct {
	kapi.TypeMeta   `json:",inline"`
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// UserNames holds all the usernames directly bound to the role
	UserNames []string `json:"userNames"`
	// GroupNames holds all the groups directly bound to the role
	GroupNames []string `json:"groupNames"`

	// Since Policy is a singleton, this is sufficient knowledge to locate a role
	// RoleRefs can only reference the current namespace and the global namespace
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef kapi.ObjectReference `json:"roleRef"`
}

// Policy is a object that holds all the Roles for a particular namespace.  There is at most
// one Policy document per namespace.
type Policy struct {
	kapi.TypeMeta   `json:",inline"`
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// LastModified is the last time that any part of the Policy was created, updated, or deleted
	LastModified kutil.Time `json:"lastModified"`

	// Roles holds all the Roles held by this Policy, mapped by Role.Name
	Roles []NamedRole `json:"roles"`
}

// PolicyBinding is a object that holds all the RoleBindings for a particular namespace.  There is
// one PolicyBinding document per referenced Policy namespace
type PolicyBinding struct {
	kapi.TypeMeta   `json:",inline"`
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// LastModified is the last time that any part of the PolicyBinding was created, updated, or deleted
	LastModified kutil.Time `json:"lastModified"`

	// PolicyRef is a reference to the Policy that contains all the Roles that this PolicyBinding's RoleBindings may reference
	PolicyRef kapi.ObjectReference `json:"policyRef"`
	// RoleBindings holds all the RoleBindings held by this PolicyBinding, mapped by RoleBinding.Name
	RoleBindings []NamedRoleBinding `json:"roleBindings"`
}

// ResourceAccessReviewResponse describes who can perform the action
type ResourceAccessReviewResponse struct {
	kapi.TypeMeta `json:",inline"`

	// Namespace is the namespace used for the access review
	Namespace string `json:"namespace,omitempty"`
	// Users is the list of users who can perform the action
	Users []string `json:"users"`
	// Groups is the list of groups who can perform the action
	Groups []string `json:"groups"`
}

// ResourceAccessReview is a means to request a list of which users and groups are authorized to perform the
// action specified by spec
type ResourceAccessReview struct {
	kapi.TypeMeta `json:",inline"`

	// Verb is one of: get, list, watch, create, update, delete
	Verb string `json:"verb"`
	// Resource is one of the existing resource types
	Resource string `json:"resource"`
	// Content is the actual content of the request for create and update
	Content kruntime.RawExtension `json:"content,omitempty"`
	// ResourceName is the name of the resource being requested for a "get" or deleted for a "delete"
	ResourceName string `json:"resourceName,omitempty"`
}

type NamedRole struct {
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type NamedRoleBinding struct {
	Name        string      `json:"name"`
	RoleBinding RoleBinding `json:"roleBinding"`
}

// SubjectAccessReviewResponse describes whether or not a user or group can perform an action
type SubjectAccessReviewResponse struct {
	kapi.TypeMeta `json:",inline"`

	// Namespace is the namespace used for the access review
	Namespace string `json:"namespace,omitempty"`
	// Allowed is required.  True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string `json:"reason,omitempty"`
}

// SubjectAccessReview is an object for requesting information about whether a user or group can perform an action
type SubjectAccessReview struct {
	kapi.TypeMeta `json:",inline"`

	// Verb is one of: get, list, watch, create, update, delete
	Verb string `json:"verb"`
	// Resource is one of the existing resource types
	Resource string `json:"resource"`
	// User is optional.  If both User and Groups are empty, the current authenticated user is used.
	User string `json:"user"`
	// Groups is optional.  Groups is the list of groups to which the User belongs.
	Groups []string `json:"groups"`
	// Content is the actual content of the request for create and update
	Content kruntime.RawExtension `json:"content,omitempty"`
	// ResourceName is the name of the resource being requested for a "get" or deleted for a "delete"
	ResourceName string `json:"resourceName"`
}

// PolicyList is a collection of Policies
type PolicyList struct {
	kapi.TypeMeta `json:",inline"`
	kapi.ListMeta `json:"metadata,omitempty"`
	Items         []Policy `json:"items"`
}

// PolicyBindingList is a collection of PolicyBindings
type PolicyBindingList struct {
	kapi.TypeMeta `json:",inline"`
	kapi.ListMeta `json:"metadata,omitempty"`
	Items         []PolicyBinding `json:"items"`
}
