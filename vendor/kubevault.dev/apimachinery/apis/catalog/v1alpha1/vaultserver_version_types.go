/*
Copyright AppsCode Inc. and Contributors

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindVaultServerVersion = "VaultServerVersion"
	ResourceVaultServerVersion     = "vaultserverversion"
	ResourceVaultServerVersions    = "vaultserverversions"
)

// VaultServerVersion defines a vaultserver version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vaultserverversions,singular=vaultserverversion,scope=Cluster,shortName=vsv,categories={vault,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="VAULT_IMAGE",type="string",JSONPath=".spec.vault.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VaultServerVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              VaultServerVersionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// VaultServerVersionSpec is the spec for postgres version
type VaultServerVersionSpec struct {
	// Version
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// Vault Image
	Vault VaultServerVersionVault `json:"vault" protobuf:"bytes,2,opt,name=vault"`
	// Unsealer Image
	Unsealer VaultServerVersionUnsealer `json:"unsealer" protobuf:"bytes,3,opt,name=unsealer"`
	// Exporter Image
	Exporter VaultServerVersionExporter `json:"exporter" protobuf:"bytes,4,opt,name=exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,5,opt,name=deprecated"`
}

// VaultServerVersionVault is the vault image
type VaultServerVersionVault struct {
	// Image is the Docker image name
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
	// ImagePullPolicy one of Always, Never, IfNotPresent. It defaults to Always if :latest is used, or IfNotPresent overwise.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,2,opt,name=imagePullPolicy,casttype=PullPolicy"`
}

// VaultServerVersionUnsealer is the image for the vault unsealer
type VaultServerVersionUnsealer struct {
	// Image is the Docker image name
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
	// ImagePullPolicy one of Always, Never, IfNotPresent. It defaults to Always if :latest is used, or IfNotPresent overwise.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,2,opt,name=imagePullPolicy,casttype=PullPolicy"`
}

// VaultServerVersionExporter is the image for the vault exporter
type VaultServerVersionExporter struct {
	// Image is the Docker image name
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
	// ImagePullPolicy one of Always, Never, IfNotPresent. It defaults to Always if :latest is used, or IfNotPresent overwise.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,2,opt,name=imagePullPolicy,casttype=PullPolicy"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// VaultServerVersionList is a list of VaultserverVersions
type VaultServerVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of VaultServerVersion CRD objects
	Items []VaultServerVersion `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
