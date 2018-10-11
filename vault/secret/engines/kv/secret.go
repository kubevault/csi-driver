package kv

import (
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func (se *EngineInfo) InitializeEngine(vc *vaultapi.Client, opts map[string]string) error {
	se.vc = vc
	se.secretName = opts["secretName"]
	se.secretDir = opts["targetDir"]
	return nil
}

func (se *EngineInfo) ReadSecret() error {
	path := fmt.Sprintf("/v1/kv/%s", se.secretName)
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

func (se *EngineInfo) RenewSecret(vol string) error {
	return nil
}

func (se *EngineInfo) StopSync() {

}
