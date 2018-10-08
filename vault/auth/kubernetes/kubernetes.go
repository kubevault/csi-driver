package kubernetes

import (
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/kubevault/csi-driver/vault/auth"
	"github.com/pkg/errors"
	"fmt"
	"io/ioutil"
	"context"
)
type AuthInfo struct {
	vaultClient *vaultapi.Client
	pod PodInfo
	authRole string
}

var _ Authentication = &AuthInfo{}

const(
	UID = "kubernetes"
	path = "v1/auth/kubernetes"
	tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)
func init()  {
	RegisterAuthMethod(UID, func(info PodInfo, client *vaultapi.Client) (Authentication, error) {
		return &AuthInfo{
			vaultClient:client,
			pod:info,
		}, nil
	})
}



func (ai *AuthInfo) GetLoginToken() (string, error) {
	req := fmt.Sprintf("%s/login", path)

	jwt, err := ai.GetJWT()
	if err != nil {
		return "", nil
	}

	r := ai.vaultClient.NewRequest("POST", req)
	if err := r.SetJSONBody(map[string]interface{}{
		"role": ai.authRole,
		"jwt": jwt,
	}); err != nil {
		return "", err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := ai.vaultClient.RawRequestWithContext(ctx, r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	secret, err :=  vaultapi.ParseSecret(resp.Body)

	if err != nil {
		return "", nil
	}
	if secret != nil {
		return "", errors.Errorf("secret not found")
	}
	return secret.Auth.ClientToken, nil
}

func (ai *AuthInfo) SetRole(role string)  {
	ai.authRole = role
}

func (ai *AuthInfo) GetSecret(p string) (*vaultapi.Secret, error)  {
	req := ai.vaultClient.NewRequest("GET", p)
	resp, err := ai.vaultClient.RawRequest(req)
	if err != nil {
		return nil, err
	}
	return vaultapi.ParseSecret(resp.Body)
}

func (ai *AuthInfo) GetJWT()(string, error)  {
	data, err := ioutil.ReadFile(tokenPath)
	return string(data), err
}