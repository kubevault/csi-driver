package kubernetes

import (
	"bytes"
	"encoding/json"
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/kubevault/csi-driver/vault/auth"
	"github.com/pkg/errors"
	"net/http"
)

type AuthInfo struct {
	vaultClient *vaultapi.Client
	pod         PodInfo
	authRole    string
	vaultUrl    string
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
	url := fmt.Sprintf("%s/%s/login", ai.vaultUrl, path)

	jwt, err := GetJWT(ai.pod.ServiceAccount, ai.pod.Namespace)
	if err != nil {
		return "", err
	}

	body := map[string]interface{}{
		"role": ai.authRole,
		"jwt":  jwt,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	secret, err := getLoginSecret(url, data)
	if err != nil {
		return "", err
	}

	if secret.Auth == nil {
		return "", errors.Errorf("secret  auth not found")
	}

	return secret.Auth.ClientToken, nil
}

func (ai *AuthInfo) SetRole(role string) {
	ai.authRole = role
}

func (ai *AuthInfo) SetVaultUrl(url string) {
	ai.vaultUrl = url
}

func (ai *AuthInfo) GetSecret(p string) (*vaultapi.Secret, error) {
	req := ai.vaultClient.NewRequest("GET", p)
	resp, err := ai.vaultClient.RawRequest(req)
	if err != nil {
		return nil, err
	}
	return vaultapi.ParseSecret(resp.Body)
}

func getLoginSecret(url string, data []byte) (*vaultapi.Secret, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return vaultapi.ParseSecret(resp.Body)
}
