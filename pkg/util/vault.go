/*
Copyright The KubeVault Authors.

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
package util

import (
	"context"
	"encoding/json"

	config "kubevault.dev/operator/apis/config/v1alpha1"
	vaultauth "kubevault.dev/operator/pkg/vault"
	authtype "kubevault.dev/operator/pkg/vault/auth/types"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

const (
	VaultRole = "secrets.csi.kubevault.com/vault-role"
)

type PodInfo struct {
	KubeClient kubernetes.Interface
	AppClient  appcat_cs.AppcatalogV1alpha1Interface

	Name           string
	Namespace      string
	UID            string
	ServiceAccount string

	RefName      string
	RefNamespace string
}

func NewVaultClient(pi *PodInfo) (*vaultapi.Client, error) {
	app, err := pi.AppClient.AppBindings(pi.RefNamespace).Get(context.TODO(), pi.RefName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var cf config.VaultServerConfiguration
	err = json.Unmarshal(app.Spec.Parameters.Raw, &cf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal parameters")
	}

	// If Kubernetes authentication information is given,
	// perform Kubernetes auth to the Vault servers.
	if cf.Kubernetes != nil {
		authInfo := &authtype.AuthInfo{
			VaultApp: app,
		}

		if cf.Kubernetes.UsePodServiceAccountForCSIDriver {
			// Get pod's service account
			sa, err := pi.KubeClient.CoreV1().ServiceAccounts(pi.Namespace).Get(context.TODO(), pi.ServiceAccount, metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get pod's service account: %s/%s", pi.Namespace, pi.ServiceAccount)
			}
			authInfo.ServiceAccountRef = &core.ObjectReference{
				Name:      sa.Name,
				Namespace: sa.Namespace,
			}

			// Get the role name from service account annotations.
			// Kubernetes authentication will be performed in the Vault server against this role.
			if pbRole, ok := sa.Annotations[VaultRole]; ok {
				authInfo.VaultRole = pbRole
			} else {
				return nil, errors.Errorf("pod's service account: %s/%s is missing annotation: %s", sa.Namespace, sa.Name, VaultRole)
			}
		} else {
			authInfo.ServiceAccountRef = &core.ObjectReference{
				Name:      cf.Kubernetes.ServiceAccountName,
				Namespace: app.Namespace,
			}
			authInfo.VaultRole = cf.VaultRole
		}

		return vaultauth.NewClientWithAppBinding(pi.KubeClient, authInfo)
	}

	return vaultauth.NewClient(pi.KubeClient, pi.AppClient, &v1alpha1.AppReference{
		Name:      app.Name,
		Namespace: app.Namespace,
	})
}
