package driver

import (
	"github.com/appscode/kutil/discovery"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func (c *Config) New() (*Driver, error) {
	if err := discovery.IsDefaultSupportedVersion(c.KubeClient); err != nil {
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
