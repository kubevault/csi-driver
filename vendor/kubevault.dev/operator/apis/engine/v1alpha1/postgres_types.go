package v1alpha1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindPostgresRole = "PostgresRole"
	ResourcePostgresRole     = "postgresrole"
	ResourcePostgresRoles    = "postgresroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgresroles,singular=postgresrole,categories={vault,appscode,all}
// +kubebuilder:subresource:status
type PostgresRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresRoleSpec   `json:"spec,omitempty"`
	Status            PostgresRoleStatus `json:"status,omitempty"`
}

// PostgresRoleSpec contains connection information, postgres role info etc
type PostgresRoleSpec struct {
	// VaultRef is the name of a AppBinding referencing to a Vault Server
	VaultRef core.LocalObjectReference `json:"vaultRef"`

	// DatabaseRef specifies the database appbinding reference in any namespace
	DatabaseRef *appcat.AppReference `json:"databaseRef,omitempty"`

	// Specifies the database name under which the role will be created
	DatabaseName string `json:"databaseName,omitempty"`

	// Specifies the path where secret engine is enabled
	Path string `json:"path,omitempty"`

	// links:
	// 	- https://www.vaultproject.io/api/secret/databases/index.html
	//	- https://www.vaultproject.io/api/secret/databases/postgresql.html

	// Specifies the TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to system/engine default TTL time
	DefaultTTL string `json:"defaultTTL,omitempty"`

	// Specifies the maximum TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to system/engine default TTL time.
	MaxTTL string `json:"maxTTL,omitempty"`

	// https://www.vaultproject.io/api/secret/databases/postgresql.html#creation_statements
	// Specifies the database statements executed to create and configure a user.
	CreationStatements []string `json:"creationStatements"`

	// https://www.vaultproject.io/api/secret/databases/postgresql.html#revocation_statements
	// Specifies the database statements to be executed to revoke a user.
	RevocationStatements []string `json:"revocationStatements,omitempty"`

	// https://www.vaultproject.io/api/secret/databases/postgresql.html#rollback_statements
	// Specifies the database statements to be executed rollback a create operation in the event of an error.
	RollbackStatements []string `json:"rollbackStatements,omitempty"`

	// https://www.vaultproject.io/api/secret/databases/postgresql.html#renew_statements
	// Specifies the database statements to be executed to renew a user.
	RenewStatements []string `json:"renewStatements,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PostgresRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of PostgresRole objects
	Items []PostgresRole `json:"items,omitempty"`
}

type PostgresRolePhase string

type PostgresRoleStatus struct {
	// ObservedGeneration is the most recent generation observed for this PostgresROle. It corresponds to the
	// PostgresROle's generation, which is updated on mutation by the API Server.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Specifies the phase of the PostgresRole
	Phase PostgresRolePhase `json:"phase,omitempty"`

	// Represents the latest available observations of a PostgresRoleBinding current state.
	Conditions []PostgresRoleCondition `json:"conditions,omitempty"`
}

// PostgresRoleCondition describes the state of a PostgresRole at a certain point.
type PostgresRoleCondition struct {
	// Type of PostgresRole condition.
	Type string `json:"type,omitempty"`

	// Status of the condition, one of True, False, Unknown.
	Status core.ConditionStatus `json:"status,omitempty"`

	// The reason for the condition's.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}
