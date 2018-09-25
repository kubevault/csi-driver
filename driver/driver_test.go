package driver

import (
	"testing"
	"github.com/kubernetes-csi/csi-test/pkg/sanity"

	"math/rand"
	"time"
	"os"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/kubevault/csi-driver/vault"
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}


func TestDriverSuite(t *testing.T) {
	socket := "/tmp/csi.sock"
	endpoint := "unix://" + socket
	if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove unix domain socket file %s, error: %s", socket, err)
	}



	driver := &Driver{
		endpoint:     endpoint,
		nodeId:       "1234567879",
		vaultClient:  nil,
		mounter:      &fakeMounter{},
		log:          logrus.New().WithField("test_enabed", true),
	}
	defer driver.Stop()

	go driver.Run()

	mntDir, err := ioutil.TempDir("", "mnt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(mntDir)

	mntStageDir, err := ioutil.TempDir("", "mnt-stage")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(mntStageDir)

	cfg := &sanity.Config{
		StagingPath: mntStageDir,
		TargetPath:  mntDir,
		Address:     endpoint,
	}

	sanity.Test(t, cfg)
}

type fakeMounter struct{}

func (f *fakeMounter) Format(source string, fsType string) error {
	return nil
}

func (f *fakeMounter) VaultMount(target string, fsType string, options map[string]string) error  {
	return nil
}

func (f *fakeMounter) Mount(source string, target string, fsType string, options ...string) error {
	return nil
}

func (f *fakeMounter) VaultUnmount(target string) error {
	return nil
}

func (f *fakeMounter) Unmount(target string) error {
	return nil
}

func (f *fakeMounter) IsFormatted(source string) (bool, error) {
	return true, nil
}
func (f *fakeMounter) IsMounted(source, target string) (bool, error) {
	return true, nil
}

func TestKVPolicy(t *testing.T) {
	client, err := vault.NewVaultClient("http://159.65.253.198:30001", "root", nil)
	fmt.Println(client, err)
	token, err := client.GetPolicyToken([]string{"nginx"}, true)
	fmt.Println(err)
	fmt.Println(token)

	c, err := vault.NewVaultClient("http://159.65.253.198:30001", "root", nil)

	path := fmt.Sprintf("/v1/kv/%s", "my-secret")
	req := c.Vc.NewRequest("GET", path)
	resp, err := c.Vc.RawRequest(req)
	fmt.Println(err)
	secret, err := vaultapi.ParseSecret(resp.Body)
	fmt.Println(secret.Data["my-value"])
}
