package util

import (
	"encoding/json"

	config "kubevault.dev/operator/apis/config/v1alpha1"
	vaultauth "kubevault.dev/operator/pkg/vault"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
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
	app, err := pi.AppClient.AppBindings(pi.RefNamespace).Get(pi.RefName, metav1.GetOptions{})
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

	return vaultauth.NewClientWithAppBinding(pi.KubeClient, binding)

}
