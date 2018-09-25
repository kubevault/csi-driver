package engines

import (
	vaultapi "github.com/hashicorp/vault/api"
	"fmt"
	"os"
	"github.com/pkg/errors"
	"path"
	"io/ioutil"
	"strings"
)
type kv struct{
	client *vaultapi.Client
	secretName string
}

//var _ Vault = kv{}

func (r *kv) ReadData() error  {
	path := fmt.Sprintf("/v1/kv/%s", r.secretName)
	req := r.client.NewRequest("GET", path)
	resp, err := r.client.RawRequest(req)
	if err != nil {
		return err
	}
	secret, err := vaultapi.ParseSecret(resp.Body)
	if err != nil {
		return err
	}
	for key, val := range secret.Data {
		writeData(key, val.(string), "")
	}
	return nil
}

func writeData(key, value, dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.Errorf("Failed to mkdir %v: %v", dir, err)
	}

	keyPath := path.Join(dir, key)

	err = ioutil.WriteFile(keyPath, []byte(strings.TrimSpace(value)), 0644)
	if err != nil {
		return err
	}

	return nil
}