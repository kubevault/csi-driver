package aws

import (
	vaultapi "github.com/hashicorp/vault/api"
	"fmt"
	"os"
	"github.com/pkg/errors"
	"path"
	"io/ioutil"
	"strings"
)
func (se *EngineInfo) InitializeEngine(vc *vaultapi.Client, secretName, secretDir string) {
	se.vc = vc
	se.secretName = secretName
	se.secretDir = secretDir
}


func (se *EngineInfo) ReadSecret() error {
	path := fmt.Sprintf("/v1/aws/creds/%s", se.secretName)
	req := se.vc.NewRequest("GET", path)
	resp, err := se.vc.RawRequest(req)
	if err != nil {
		return err
	}
	secret, err := vaultapi.ParseSecret(resp.Body)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(se.secretDir, 0755); err != nil {
		return errors.Errorf("Failed to mkdir %v: %v", se.secretDir, err)
	}

	for key, val := range secret.Data {
		if err := writeData(key, val.(string), se.secretDir); err != nil {
			return err
		}
	}
	return nil
}

func writeData(key, value, dir string) error {
	keyPath := path.Join(dir, key)
	return ioutil.WriteFile(keyPath, []byte(strings.TrimSpace(value)), 0644)
}