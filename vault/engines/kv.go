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
	dir string
}

//var _ Vault = kv{}

func NewKVEngine(client *vaultapi.Client, secretName, dir string) *kv   {
	return &kv{client:client, secretName:secretName, dir:dir}
}

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
		if err := writeData(key, val.(string), r.dir); err != nil {
			return err
		}
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