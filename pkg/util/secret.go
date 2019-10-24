package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	vaultEng "kubevault.dev/operator/pkg/vault/secret/engines"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// FetchSecret writes the secret on specified targetPath
// Returns LeaseID and error
func FetchSecret(engineName string, vc *vaultapi.Client, opts map[string]string, targetPath string) (*vaultapi.Secret, error) {
	engine, err := vaultEng.NewSecretManager(engineName)
	if err != nil {
		return nil, err
	}

	if err = engine.SetOptions(vc, opts); err != nil {
		return nil, err
	}

	secret, err := engine.GetSecret()
	if err != nil {
		return nil, err
	}

	if err = ensurePath(targetPath); err != nil {
		return secret, err
	}

	for key, val := range secret.Data {
		if val != nil {
			if err := writeData(targetPath, key, val); err != nil {
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

func writeData(dir, key string, value interface{}) error {
	keyPath := path.Join(dir, key)
	return ioutil.WriteFile(keyPath, []byte(fmt.Sprintf("%v", value)), 0644)
}

func ensurePath(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.Errorf("Failed to mkdir %v: %v", path, err)
	}
	return nil
}
