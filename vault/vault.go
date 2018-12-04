package vault

import (
	"encoding/json"
	vaultapi "github.com/hashicorp/vault/api"
	config "github.com/kubevault/operator/apis/config/v1alpha1"
	vaultauth "github.com/kubevault/operator/pkg/vault"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

type PodInfo struct {
	Name           string
	Namespace      string
	UID            string
	ServiceAccount string

	RefName      string
	RefNamespace string
}

func GetAppBindingVaultClient(pi *PodInfo) (*vaultapi.Client, error) {

	kubeClient, err := getKubeClient()
	if err != nil {
		return nil, err
	}

	app, err := getAppBinding(pi.RefName, pi.RefNamespace)
	if err != nil {
		return nil, err
	}

	var cf config.VaultServerConfiguration
	err = json.Unmarshal(app.Spec.Parameters.Raw, &cf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal parameters")
	}

	binding := app.DeepCopy()
	binding.Namespace = pi.Namespace

	if cf.UsePodServiceAccountForCSIDriver {
		binding.Spec.Secret = nil
		cf.ServiceAccountName = pi.ServiceAccount
	}

	rawData, err := json.Marshal(cf)
	if err != nil {
		return nil, err
	}

	binding.Spec.Parameters = &runtime.RawExtension{
		Raw: rawData,
	}

	return vaultauth.NewClientWithAppBinding(kubeClient, binding)

}

func getAppBinding(appName, appNamespace string) (*appcat.AppBinding, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	appClient, err := appcat_cs.NewForConfig(config)

	app, err := appClient.AppBindings(appNamespace).Get(appName, metav1.GetOptions{})
	return app, err
}

func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}
