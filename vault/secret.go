package vault

import (
	"encoding/json"
	vaultapi "github.com/hashicorp/vault/api"
	vaultEng "github.com/kubevault/operator/pkg/vault/secret/engines"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ReasdSecretData writes the secret on specified targetPath
// Returns LeaseID and error
func ReadSecretData(engineName string, vc *vaultapi.Client, opts map[string]string, targetPath string) (*vaultapi.Secret, error) {
	engine, err := vaultEng.NewSecretManager(engineName)
	if err != nil {
		return nil, err
	}

	if err = engine.SetOptions(vc, opts); err != nil {
		return nil, err
	}

	data, err := engine.GetSecret()
	if err != nil {
		return nil, err
	}

	secret := &vaultapi.Secret{}
	if err = json.Unmarshal(data, secret); err != nil {
		return nil, err
	}

	if err = ensurePath(targetPath); err != nil {
		return secret, err
	}

	for key, val := range secret.Data {
		if val != nil {
			if err := writeData(key, val.(string), targetPath); err != nil {
				return secret, err
			}
		}
	}
	return secret, nil
}

func SetRenewal(vc *vaultapi.Client, secret *vaultapi.Secret) (*vaultapi.Renewer, error) {
	renewer, err := vc.NewRenewer(&vaultapi.RenewerInput{
		Secret: secret,
	})
	if err != nil {
		return nil, err
	}
	renewer.Renew()
	return renewer, nil
}

func StopRenew(renewer *vaultapi.Renewer) {
	renewer.Stop()
}

func writeData(key, value, dir string) error {
	keyPath := path.Join(dir, key)
	return ioutil.WriteFile(keyPath, []byte(strings.TrimSpace(value)), 0644)
}

func ensurePath(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.Errorf("Failed to mkdir %v: %v", path, err)
	}
	return nil
}
