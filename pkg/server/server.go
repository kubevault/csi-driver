package server

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"kubevault.dev/csi-driver/pkg/driver"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

type VaultDriverConfig struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   *driver.Config
}

// VaultDriver contains state for a Kubernetes cluster master/api server.
type VaultDriver struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
	Driver           *driver.Driver
}

func (op *VaultDriver) Run(stopCh <-chan struct{}) error {
	go utilruntime.Must(op.Driver.Run())
	return op.GenericAPIServer.PrepareRun().Run(stopCh)
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *driver.Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *VaultDriverConfig) Complete() CompletedConfig {
	completedCfg := completedConfig{
		c.GenericConfig.Complete(),
		c.ExtraConfig,
	}

	completedCfg.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "1",
	}

	return CompletedConfig{&completedCfg}
}

// New returns a new instance of VaultDriver from the given config.
func (c completedConfig) New() (*VaultDriver, error) {
	genericServer, err := c.GenericConfig.New("vault-csi-driver", genericapiserver.NewEmptyDelegate()) // completion is done in Complete, no need for a second time
	if err != nil {
		return nil, err
	}

	driver, err := c.ExtraConfig.New()
	if err != nil {
		return nil, err
	}

	s := &VaultDriver{
		GenericAPIServer: genericServer,
		Driver:           driver,
	}
	return s, nil
}
