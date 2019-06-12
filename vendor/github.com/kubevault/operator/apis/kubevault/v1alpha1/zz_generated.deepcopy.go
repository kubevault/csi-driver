// +build !ignore_autogenerated

/*
Copyright 2019 The Kube Vault Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	apiv1 "kmodules.xyz/monitoring-agent-api/api/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthConfig) DeepCopyInto(out *AuthConfig) {
	*out = *in
	if in.AuditNonHMACRequestKeys != nil {
		in, out := &in.AuditNonHMACRequestKeys, &out.AuditNonHMACRequestKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AuditNonHMACResponseKeys != nil {
		in, out := &in.AuditNonHMACResponseKeys, &out.AuditNonHMACResponseKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PassthroughRequestHeaders != nil {
		in, out := &in.PassthroughRequestHeaders, &out.PassthroughRequestHeaders
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthConfig.
func (in *AuthConfig) DeepCopy() *AuthConfig {
	if in == nil {
		return nil
	}
	out := new(AuthConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthMethod) DeepCopyInto(out *AuthMethod) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(AuthConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthMethod.
func (in *AuthMethod) DeepCopy() *AuthMethod {
	if in == nil {
		return nil
	}
	out := new(AuthMethod)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthMethodStatus) DeepCopyInto(out *AuthMethodStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthMethodStatus.
func (in *AuthMethodStatus) DeepCopy() *AuthMethodStatus {
	if in == nil {
		return nil
	}
	out := new(AuthMethodStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AwsKmsSsmSpec) DeepCopyInto(out *AwsKmsSsmSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AwsKmsSsmSpec.
func (in *AwsKmsSsmSpec) DeepCopy() *AwsKmsSsmSpec {
	if in == nil {
		return nil
	}
	out := new(AwsKmsSsmSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureKeyVault) DeepCopyInto(out *AzureKeyVault) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureKeyVault.
func (in *AzureKeyVault) DeepCopy() *AzureKeyVault {
	if in == nil {
		return nil
	}
	out := new(AzureKeyVault)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureSpec) DeepCopyInto(out *AzureSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureSpec.
func (in *AzureSpec) DeepCopy() *AzureSpec {
	if in == nil {
		return nil
	}
	out := new(AzureSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendStorageSpec) DeepCopyInto(out *BackendStorageSpec) {
	*out = *in
	if in.Inmem != nil {
		in, out := &in.Inmem, &out.Inmem
		*out = new(InmemSpec)
		**out = **in
	}
	if in.Etcd != nil {
		in, out := &in.Etcd, &out.Etcd
		*out = new(EtcdSpec)
		**out = **in
	}
	if in.Gcs != nil {
		in, out := &in.Gcs, &out.Gcs
		*out = new(GcsSpec)
		**out = **in
	}
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		*out = new(S3Spec)
		**out = **in
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(AzureSpec)
		**out = **in
	}
	if in.PostgreSQL != nil {
		in, out := &in.PostgreSQL, &out.PostgreSQL
		*out = new(PostgreSQLSpec)
		**out = **in
	}
	if in.MySQL != nil {
		in, out := &in.MySQL, &out.MySQL
		*out = new(MySQLSpec)
		**out = **in
	}
	if in.File != nil {
		in, out := &in.File, &out.File
		*out = new(FileSpec)
		**out = **in
	}
	if in.DynamoDB != nil {
		in, out := &in.DynamoDB, &out.DynamoDB
		*out = new(DynamoDBSpec)
		**out = **in
	}
	if in.Swift != nil {
		in, out := &in.Swift, &out.Swift
		*out = new(SwiftSpec)
		**out = **in
	}
	if in.Consul != nil {
		in, out := &in.Consul, &out.Consul
		*out = new(ConsulSpec)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendStorageSpec.
func (in *BackendStorageSpec) DeepCopy() *BackendStorageSpec {
	if in == nil {
		return nil
	}
	out := new(BackendStorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsulSpec) DeepCopyInto(out *ConsulSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsulSpec.
func (in *ConsulSpec) DeepCopy() *ConsulSpec {
	if in == nil {
		return nil
	}
	out := new(ConsulSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamoDBSpec) DeepCopyInto(out *DynamoDBSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamoDBSpec.
func (in *DynamoDBSpec) DeepCopy() *DynamoDBSpec {
	if in == nil {
		return nil
	}
	out := new(DynamoDBSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EtcdSpec) DeepCopyInto(out *EtcdSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EtcdSpec.
func (in *EtcdSpec) DeepCopy() *EtcdSpec {
	if in == nil {
		return nil
	}
	out := new(EtcdSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSpec) DeepCopyInto(out *FileSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSpec.
func (in *FileSpec) DeepCopy() *FileSpec {
	if in == nil {
		return nil
	}
	out := new(FileSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GcsSpec) DeepCopyInto(out *GcsSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GcsSpec.
func (in *GcsSpec) DeepCopy() *GcsSpec {
	if in == nil {
		return nil
	}
	out := new(GcsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleKmsGcsSpec) DeepCopyInto(out *GoogleKmsGcsSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleKmsGcsSpec.
func (in *GoogleKmsGcsSpec) DeepCopy() *GoogleKmsGcsSpec {
	if in == nil {
		return nil
	}
	out := new(GoogleKmsGcsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InmemSpec) DeepCopyInto(out *InmemSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InmemSpec.
func (in *InmemSpec) DeepCopy() *InmemSpec {
	if in == nil {
		return nil
	}
	out := new(InmemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesSecretSpec) DeepCopyInto(out *KubernetesSecretSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesSecretSpec.
func (in *KubernetesSecretSpec) DeepCopy() *KubernetesSecretSpec {
	if in == nil {
		return nil
	}
	out := new(KubernetesSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModeSpec) DeepCopyInto(out *ModeSpec) {
	*out = *in
	if in.KubernetesSecret != nil {
		in, out := &in.KubernetesSecret, &out.KubernetesSecret
		*out = new(KubernetesSecretSpec)
		**out = **in
	}
	if in.GoogleKmsGcs != nil {
		in, out := &in.GoogleKmsGcs, &out.GoogleKmsGcs
		*out = new(GoogleKmsGcsSpec)
		**out = **in
	}
	if in.AwsKmsSsm != nil {
		in, out := &in.AwsKmsSsm, &out.AwsKmsSsm
		*out = new(AwsKmsSsmSpec)
		**out = **in
	}
	if in.AzureKeyVault != nil {
		in, out := &in.AzureKeyVault, &out.AzureKeyVault
		*out = new(AzureKeyVault)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModeSpec.
func (in *ModeSpec) DeepCopy() *ModeSpec {
	if in == nil {
		return nil
	}
	out := new(ModeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MySQLSpec) DeepCopyInto(out *MySQLSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MySQLSpec.
func (in *MySQLSpec) DeepCopy() *MySQLSpec {
	if in == nil {
		return nil
	}
	out := new(MySQLSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgreSQLSpec) DeepCopyInto(out *PostgreSQLSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgreSQLSpec.
func (in *PostgreSQLSpec) DeepCopy() *PostgreSQLSpec {
	if in == nil {
		return nil
	}
	out := new(PostgreSQLSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3Spec) DeepCopyInto(out *S3Spec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3Spec.
func (in *S3Spec) DeepCopy() *S3Spec {
	if in == nil {
		return nil
	}
	out := new(S3Spec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SwiftSpec) DeepCopyInto(out *SwiftSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SwiftSpec.
func (in *SwiftSpec) DeepCopy() *SwiftSpec {
	if in == nil {
		return nil
	}
	out := new(SwiftSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSPolicy) DeepCopyInto(out *TLSPolicy) {
	*out = *in
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSPolicy.
func (in *TLSPolicy) DeepCopy() *TLSPolicy {
	if in == nil {
		return nil
	}
	out := new(TLSPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UnsealerSpec) DeepCopyInto(out *UnsealerSpec) {
	*out = *in
	in.Mode.DeepCopyInto(&out.Mode)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UnsealerSpec.
func (in *UnsealerSpec) DeepCopy() *UnsealerSpec {
	if in == nil {
		return nil
	}
	out := new(UnsealerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServer) DeepCopyInto(out *VaultServer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServer.
func (in *VaultServer) DeepCopy() *VaultServer {
	if in == nil {
		return nil
	}
	out := new(VaultServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VaultServer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerCondition) DeepCopyInto(out *VaultServerCondition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerCondition.
func (in *VaultServerCondition) DeepCopy() *VaultServerCondition {
	if in == nil {
		return nil
	}
	out := new(VaultServerCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerList) DeepCopyInto(out *VaultServerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VaultServer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerList.
func (in *VaultServerList) DeepCopy() *VaultServerList {
	if in == nil {
		return nil
	}
	out := new(VaultServerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VaultServerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerSpec) DeepCopyInto(out *VaultServerSpec) {
	*out = *in
	if in.ConfigSource != nil {
		in, out := &in.ConfigSource, &out.ConfigSource
		*out = new(v1.VolumeSource)
		(*in).DeepCopyInto(*out)
	}
	if in.DataSources != nil {
		in, out := &in.DataSources, &out.DataSources
		*out = make([]v1.VolumeSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(TLSPolicy)
		(*in).DeepCopyInto(*out)
	}
	in.Backend.DeepCopyInto(&out.Backend)
	if in.Unsealer != nil {
		in, out := &in.Unsealer, &out.Unsealer
		*out = new(UnsealerSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.AuthMethods != nil {
		in, out := &in.AuthMethods, &out.AuthMethods
		*out = make([]AuthMethod, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Monitor != nil {
		in, out := &in.Monitor, &out.Monitor
		*out = new(apiv1.AgentSpec)
		(*in).DeepCopyInto(*out)
	}
	in.PodTemplate.DeepCopyInto(&out.PodTemplate)
	in.ServiceTemplate.DeepCopyInto(&out.ServiceTemplate)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerSpec.
func (in *VaultServerSpec) DeepCopy() *VaultServerSpec {
	if in == nil {
		return nil
	}
	out := new(VaultServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultServerStatus) DeepCopyInto(out *VaultServerStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = (*in).DeepCopy()
	}
	in.VaultStatus.DeepCopyInto(&out.VaultStatus)
	if in.UpdatedNodes != nil {
		in, out := &in.UpdatedNodes, &out.UpdatedNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]VaultServerCondition, len(*in))
		copy(*out, *in)
	}
	if in.AuthMethodStatus != nil {
		in, out := &in.AuthMethodStatus, &out.AuthMethodStatus
		*out = make([]AuthMethodStatus, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultServerStatus.
func (in *VaultServerStatus) DeepCopy() *VaultServerStatus {
	if in == nil {
		return nil
	}
	out := new(VaultServerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultStatus) DeepCopyInto(out *VaultStatus) {
	*out = *in
	if in.Standby != nil {
		in, out := &in.Standby, &out.Standby
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Sealed != nil {
		in, out := &in.Sealed, &out.Sealed
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Unsealed != nil {
		in, out := &in.Unsealed, &out.Unsealed
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultStatus.
func (in *VaultStatus) DeepCopy() *VaultStatus {
	if in == nil {
		return nil
	}
	out := new(VaultStatus)
	in.DeepCopyInto(out)
	return out
}
