package kubernetes

import (
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/kubevault/csi-driver/vault/auth"
	vaultauth "github.com/kubevault/operator/pkg/vault/auth"
)

type AuthInfo struct {
	vaultClient  *vaultapi.Client
	pod          PodInfo
	authRole     string
	refName      string
	refNamespace string
}

var _ Authentication = &AuthInfo{}

const (
	UID  = "kubernetes"
	path = "v1/auth/kubernetes"
)

func init() {
	RegisterAuthMethod(UID, func(info PodInfo, client *vaultapi.Client) (Authentication, error) {
		return &AuthInfo{
			vaultClient: client,
			pod:         info,
		}, nil
	})
}

func (ai *AuthInfo) GetLoginToken() (string, error) {
	kubeClient, err := getKubeClient()
	if err != nil {
		return "", err
	}
	app, err := getAppBinding(ai.refName, ai.refNamespace)
	if err != nil {
		return "", err
	}
	binding := app.DeepCopy()
	binding.Spec.Secret.Name, err = getServiceAccountSecret(kubeClient, ai.pod.ServiceAccount, ai.pod.Namespace)
	if err != nil {
		return "", err
	}

	vAuth, err := vaultauth.NewAuth(kubeClient, binding)
	return vAuth.Login()
}

func (ai *AuthInfo) SetRef(name, namespace string) {
	ai.refName = name
	ai.refNamespace = namespace
}

func (ai *AuthInfo) GetSecret(p string) (*vaultapi.Secret, error) {
	req := ai.vaultClient.NewRequest("GET", p)
	resp, err := ai.vaultClient.RawRequest(req)
	if err != nil {
		return nil, err
	}
	return vaultapi.ParseSecret(resp.Body)
}
