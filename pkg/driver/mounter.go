/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/juju/errors"
)

// Mounter is responsible for formatting and mounting volumes
type Mounter interface {
	// Format formats the source with the given filesystem type
	Format(source, fsType string) error

	//VaultMount(target, fsType string, options map[string]string ) error
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

type mounter struct{}

func (m *mounter) Format(source, fsType string) error {
	return nil
}

//https://medium.com/@gmaliar/dynamic-secrets-on-kubernetes-pods-using-vault-35d9094d169

func (m *mounter) Mount(source, target, fsType string, opts ...string) error {
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

	// create target, os.Mkdirall is noop if it exists

	err := os.MkdirAll(target, 0755)
	if err != nil {
		return err
	}

	mountArgs = append(mountArgs, "-t", fsType)

	if len(opts) > 0 {
		mountArgs = append(mountArgs, "-o", strings.Join(opts, ","))
	}
	mountArgs = append(mountArgs, source, target)

	//mountArgs = append(mountArgs, target)

	fmt.Println(mountArgs)
	out, err := exec.Command(mountCmd, mountArgs...).CombinedOutput()
	if err != nil {
		return errors.Errorf("mounting failed: %v cmd: '%s %s' output: %q",
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
		return errors.Errorf("unmounting failed: %v cmd: '%s %s' output: %q",
			err, umountCmd, target, string(out))
	}

	return nil
}

func (m *mounter) IsFormatted(source string) (bool, error) {
	return true, nil
}

func (m *mounter) IsMounted(source, target string) (bool, error) {
	findmntCmd := "findmnt"
	_, err := exec.LookPath(findmntCmd)
	if err != nil {
		if err == exec.ErrNotFound {
			return false, errors.Errorf("%q executable not found in $PATH", findmntCmd)
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

		return false, errors.Errorf("checking mounted failed: %v cmd: %q output: %q",
			err, findmntCmd, string(out))
	}

	if strings.TrimSpace(string(out)) == "" {
		return false, nil
	}

	return true, nil
}

func cleanup(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return errors.Errorf("Failed to remove the directory %v: %v", dir, err)
	}
	return nil
}
