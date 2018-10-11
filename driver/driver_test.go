package driver

import (
	"testing"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"

	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"os/exec"
	"path/filepath"

	"bytes"
	"context"
	"encoding/json"
	"net/http"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/kubevault/csi-driver/vault"
	"github.com/sirupsen/logrus"
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
		endpoint:    endpoint,
		nodeId:      "1234567879",
		vaultClient: nil,
		mounter:     &fakeMounter{},
		log:         logrus.New().WithField("test_enabed", true),
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

func (f *fakeMounter) VaultMount(target string, fsType string, options map[string]string) error {
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

	client, err := vault.NewVaultClient("http://142.93.77.58:30001", "root", nil)
	fmt.Println(client.Headers(), err)

	p := "v1/auth/kubernetes/login"

	r := client.NewRequest("POST", p)
	body := map[string]interface{}{
		"role": "testrole",
		"jwt":  "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InBvc3RncmVzLXZhdWx0LXRva2VuLXg5djRyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InBvc3RncmVzLXZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiOTVjM2FmNzAtY2FiZS0xMWU4LWExMzQtYTZmNTM5NDhkMzQ0Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6cG9zdGdyZXMtdmF1bHQifQ.Q9xGPbPt_cNwRrzlX-kFJsN7eJPceAYP7P7rkU8VZPeEuIip2jFoSbF8LZsL6TfB7XhUnvQT4NTXry-XVQA_Rvfrs61IPtvw0HOwkCxtd0PglW1p53B_onH6NknofRT0ZThoC9Jhs8NYa4FkTyyK1Wo46_aZQ2XbCny9UZzOjBxYo8iv_OL3crIytQV6UjrA2q-XkJuGCRc_vvXpPS4KO3ke7dsjrCwOTTz8QRGiljyscHzCJmN733VxvGSDuDoonxty894DhqsL6iRHKS5X8UVaq3MGNyndfQBSJfUnT75dFYD12Cr_BZRONBF66iGSXbaa-_Ft-eTgCEq0o_j2Nw",
	}
	d, e := json.Marshal(body)
	fmt.Println(string(d), e)
	if err := r.SetJSONBody(body); err != nil {
		fmt.Println(err, "***************")
	}
	r.Headers = make(map[string][]string)
	r.Headers.Set("Content-Type", "application/json")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := client.RawRequestWithContext(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)
	//secret, err :=  vaultapi.ParseSecret(resp.Body)
	//fmt.Println(secret, err)
	/*return



	token, err := client.GetPolicyToken([]string{"nginx"}, true)
	fmt.Println(err)
	fmt.Println(token)

	c, err := vault.NewVaultClient("http://159.65.253.198:30001", "root", nil)

	path := fmt.Sprintf("/v1/kv/%s", "my-secret")
	req := c.Vc.NewRequest("GET", path)
	resp, err := c.Vc.RawRequest(req)
	fmt.Println(err)
	secret, err := vaultapi.ParseSecret(resp.Body)
	fmt.Println(secret.Data["my-value"])*/
}

func TestVault(t *testing.T) {
	c, err := vault.NewVaultClient("http://142.93.77.58:30001", "root", nil)

	path := fmt.Sprintf("/v1/pki/roles/%s", "my-pki-role")
	req := c.NewRequest("GET", path)
	resp, err := c.RawRequest(req)
	fmt.Println(err)
	secret, err := vaultapi.ParseSecret(resp.Body)
	fmt.Println(secret.Data)
}

func TestAT(t *testing.T) {
	l := map[string]interface{}{
		"max_ttl": 259200,
	}
	fmt.Println(l["max_ttl"].(int))

}

func TestPKI(t *testing.T) {
	c, err := vault.NewVaultClient("http://142.93.77.58:30001", "root", nil)

	r := c.NewRequest("POST", "/v1/pki/issue/my-pki-role")
	if err := r.SetJSONBody(map[string]string{
		"common_name": "www.my-website.com",
	}); err != nil {
		fmt.Println(err, "**********")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.RawRequestWithContext(ctx, r)
	if err != nil {
		fmt.Println(err, ")))")
	}
	defer resp.Body.Close()

	secret, err := vaultapi.ParseSecret(resp.Body)
	fmt.Println(secret.Data)
}

func TestHttp(t *testing.T) {
	url := "http://142.93.77.58:30001/v1/pki/issue/my-pki-role"
	body := map[string]interface{}{
		"common_name": "www.my-website.com",
	}
	d, e := json.Marshal(body)

	fmt.Println(e)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	fmt.Println(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	bdy, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bdy))
}

func TestPath(t *testing.T) {
	path := "/var/www"
	fmt.Println(filepath.Join(path, "*"))

	//fmt.Println(os.Link("/home/sanjid/test/a", "/home/sanjid/test/b"))
	//return

	args := []string{
		"-s",
		"/home/sanjid/test/a/*",
		"/home/sanjid/test/b",
		"-v",
	}
	fmt.Println(args)
	err := exec.Command("ln", args...).Run()
	fmt.Println(err)
}

/*


vault write pki/config/urls \
    issuing_certificates="http://142.93.77.58:30001/v1/pki/ca" \
    crl_distribution_points="http://142.93.77.58:30001/v1/pki/crl"

vault write pki/roles/my-pki-role \
    allowed_domains=my-website.com \
    allow_subdomains=true \
    max_ttl=72h

vault write pki/issue/my-pki-role \
    common_name=www.my-website.com
*/
