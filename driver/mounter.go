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
	"os/exec"
)

// Mounter is responsible for formatting and mounting volumes
type Mounter interface {
	// Format formats the source with the given filesystem type
	Format(source, fsType string) error

	VaultMount(target, fsType string, options map[string]string ) error
	// Mount mounts source to target with the given fstype and options.
	Mount(source, target, fsType string, options ...string) error

	// Unmount unmounts the given target
	Unmount(target string) error
	VaultUnmount(target string) error

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

func (m *mounter) VaultMount(target, fsType string, opts map[string]string) error {
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

func (m *mounter) Mount(source, target, fsType string, opts ...string) error  {
	mountCmd := "mount"
	mountArgs := []string{}

	if fsType == "" {
		return errors.New("fs type is not specified for mounting the volume")
	}

	if source == "" {
		return errors.New("source is not specified for mounting the volume")
	}

	if target == "" {
		return errors.New("target is not specified for mounting the volume")
	}

	mountArgs = append(mountArgs, "-t", fsType)

	if len(opts) > 0 {
		mountArgs = append(mountArgs, "-o", strings.Join(opts, ","))
	}

	mountArgs = append(mountArgs, source)
	mountArgs = append(mountArgs, target)

	// create target, os.Mkdirall is noop if it exists
	err := os.MkdirAll(target, 0750)
	if err != nil {
		return err
	}

	out, err := exec.Command(mountCmd, mountArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mounting failed: %v cmd: '%s %s' output: %q",
			err, mountCmd, strings.Join(mountArgs, " "), string(out))
	}

	return nil
}


func (m *mounter) VaultUnmount(target string) error {
	return cleanup(target)
}

func (m *mounter) Unmount(target string) error {
	umountCmd := "umount"
	if target == "" {
		return errors.New("target is not specified for unmounting the volume")
	}

	out, err := exec.Command("umount", target).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unmounting failed: %v cmd: '%s %s' output: %q",
			err, umountCmd, target, string(out))
	}

	return nil
}

func (m *mounter) IsFormatted(source string) (bool, error) {
	if source == "" {
		return false, errors.New("source is not specified")
	}

	blkidCmd := "blkid"
	_, err := exec.LookPath(blkidCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			return false, fmt.Errorf("%q executable not found in $PATH", blkidCmd)
		}
		return false, err
	}

	out, err := exec.Command(blkidCmd, source).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("checking formatting failed: %v cmd: %q output: %q",
			err, blkidCmd, string(out))
	}

	if strings.TrimSpace(string(out)) == "" {
		return false, nil
	}

	return true, nil
}

func (m *mounter) IsMounted(source, target string) (bool, error) {
	findmntCmd := "findmnt"
	_, err := exec.LookPath(findmntCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			return false, fmt.Errorf("%q executable not found in $PATH", findmntCmd)
		}
		return false, err
	}

	findmntArgs := []string{"--mountpoint", target}
	if source != "" {
		findmntArgs = append(findmntArgs, "--source", source)
	}

	out, err := exec.Command(findmntCmd, findmntArgs...).CombinedOutput()
	if err != nil {
		// findmnt exits with non zero exit status if it couldn't find anything
		if strings.TrimSpace(string(out)) == "" {
			return false, nil
		}

		return false, fmt.Errorf("checking mounted failed: %v cmd: %q output: %q",
			err, findmntCmd, string(out))
	}

	if strings.TrimSpace(string(out)) == "" {
		return false, nil
	}

	return true, nil
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
	err := os.RemoveAll(dir)
	if err != nil {
		return errors.Errorf("Failed to remove the directory %v: %v", dir, err)
	}
	return nil
}
