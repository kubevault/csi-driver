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
package driver

import (
	"context"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kmodules.xyz/client-go/discovery"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

type config struct {
	Endpoint string
	NodeId   string
}

type Config struct {
	config

	ClientConfig *rest.Config
	KubeClient   kubernetes.Interface
	AppClient    appcat_cs.AppcatalogV1alpha1Interface
	CRDClient    crd_cs.ApiextensionsV1beta1Interface
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func isSupportedVersion(kc kubernetes.Interface) error {
	return discovery.IsSupportedVersion(
		kc,
		">= 1.14.0", // supported versions: 1.14.x and higher
		discovery.DefaultBlackListedVersions,
		discovery.DefaultBlackListedMultiMasterVersions)
}

func (c *Config) EnsureCRDs() error {
	crds := []*crd_api.CustomResourceDefinition{
		appcat.AppBinding{}.CustomResourceDefinition(),
	}
	return crdutils.RegisterCRDs(context.TODO(), c.KubeClient.Discovery(), c.CRDClient, crds)
}

func (c *Config) New() (*Driver, error) {
	if err := isSupportedVersion(c.KubeClient); err != nil {
		return nil, err
	}

	return &Driver{
		config:     c.config,
		mounter:    &mounter{},
		kubeClient: c.KubeClient,
		appClient:  c.AppClient,
		log: logrus.New().WithFields(logrus.Fields{
			"node-id": c.NodeId,
		}),
		ch: make(map[string]*vaultapi.Renewer),
	}, nil
}
