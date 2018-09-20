package driver

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/juju/errors"
	"fmt"
)

// Mounter is responsible for formatting and mounting volumes
type Mounter interface {
	// Format formats the source with the given filesystem type
	Format(source, fsType string) error

	// Mount mounts source to target with the given fstype and options.
	Mount(target, fsType string, options map[string]string) error

	// Unmount unmounts the given target
	Unmount(target string) error

	// IsFormatted checks whether the source device is formatted or not. It
	// returns true if the source device is already formatted.
	IsFormatted(source string) (bool, error)

	// IsMounted checks whether the source device is mounted to the target
	// path. Source can be empty. In that case it only checks whether the
	// device is mounted or not.
	// It returns true if it's mounted.
	IsMounted(source, target string) (bool, error)
}

type mounter struct{
	vaultUrl string
	token string
}

func (m *mounter) Format(source, fsType string) error {

	return nil
}

//https://medium.com/@gmaliar/dynamic-secrets-on-kubernetes-pods-using-vault-35d9094d169
func (m *mounter) Mount(target, fsType string, opts map[string]string) error {
	stringPolicies, ok := opts["policy"]
	if !ok {
		return errors.Errorf("Missing policies")
	}
	policies := strings.Split(strings.Replace(stringPolicies, " ", "", -1), ",")
	if len(policies) == 0 {
		return errors.Errorf("Empty policies")
	}

	poduidreg := regexp.MustCompile("[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{8}")
	poduid := poduidreg.FindString(target)
	fmt.Println(poduid)
	if poduid == "" {
		//poduid = "poduid123456"
		return errors.Errorf("Couldn't extract poduid from path %v", target)
	}

	client, err := NewVaultClient(m.vaultUrl, m.token, &vaultapi.TLSConfig{Insecure: true})
	if err != nil {
		return err
	}
	token, metadata, err := client.GetTokenData(policies, poduid, true)
	if err != nil {
		return err
	}
	err = writeTokenData(token, metadata, target, "vault-token")
	if err != nil {
		return err
	}
	return nil
}

func (m *mounter) Unmount(target string) error {
	return cleanup(target)
}

func (m *mounter) IsFormatted(source string) (bool, error) {
	return false, nil
}

func (m *mounter) IsMounted(source, target string) (bool, error) {
	return false, nil
}

func writeTokenData(token string, metadata []byte, dir, tokenfilename string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.Errorf("Failed to mkdir %v: %v", dir, err)
	}

	tokenpath := path.Join(dir, tokenfilename)
	fulljsonpath := path.Join(dir, strings.Join([]string{tokenfilename, ".json"}, ""))

	err = ioutil.WriteFile(tokenpath, []byte(strings.TrimSpace(token)), 0644)
	if err != nil {
		return err
	}
	err = os.Chmod(tokenpath, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fulljsonpath, metadata, 0644)
	if err != nil {
		return err
	}
	err = os.Chmod(fulljsonpath, 0644)
	return err
}

func cleanup(dir string) error {
	// Good Guy RemoveAll does nothing is path doesn't exist and returns nil error :)
	err := os.RemoveAll(dir)
	if err != nil {
		return errors.Errorf("Failed to remove the directory %v: %v", dir, err)
	}
	return nil
}
