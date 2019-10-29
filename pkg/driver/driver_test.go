/*
Copyright The KubeVault Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package driver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/appscode/pat"
	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	"github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	cr "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appFake "kmodules.xyz/custom-resources/client/clientset/versioned/fake"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

const testNamespace = "default"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestDriverSuite(t *testing.T) {
	socket := "/tmp/csi.sock"
	endpoint := "unix://" + socket
	if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove unix domain socket file %s, error: %s", socket, err)
	}

	ts := NewFakeVaultServer()
	defer ts.Close()

	fakeAppClient, err := getAppBindingWithFakeClient(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	fakeKubeClient := getFakeKubeClient()
	if err = setupKubernetes(fakeKubeClient); err != nil {
		t.Fatal(err)
	}

	driver := &Driver{
		config: config{
			Endpoint: endpoint,
			NodeId:   "1234567879",
		},
		vaultClient: nil,
		mounter:     &fakeMounter{},
		log:         logrus.New().WithField("test_enabed", true),

		kubeClient: fakeKubeClient,
		appClient:  fakeAppClient,
	}
	defer driver.Stop()

	go utilruntime.Must(driver.Run())

	tp := os.TempDir() + "/csi-target"
	sp := os.TempDir() + "/csi-staging"
	defer utilruntime.Must(os.RemoveAll(tp))
	defer utilruntime.Must(os.RemoveAll(sp))

	cfg := &sanity.Config{
		TargetPath:  tp,
		StagingPath: sp,
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

func NewFakeVaultServer() *httptest.Server {
	authResp := `
{
  "auth": {
    "client_token": "1234"
  }
}
`
	secResp := `
{
  "auth": null,
  "data": {
    "foo": "bar"
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}`
	m := pat.New()
	m.Post("/v1/auth/kubernetes/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v map[string]interface{}
		defer utilruntime.Must(r.Body.Close())
		utilruntime.Must(json.NewDecoder(r.Body).Decode(&v))
		if val, ok := v["jwt"]; ok {
			if val.(string) == "sanity-token" {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(authResp))
				utilruntime.Must(err)
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	m.Get("/v1/kv/:secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer utilruntime.Must(r.Body.Close())
		if params, found := pat.FromContext(r.Context()); found {
			if got, want := params.Get(":secret"), "my-key"; got == want {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(secResp))
				utilruntime.Must(err)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
	}))

	return httptest.NewServer(m)
}

func getFakeKubeClient() *fake.Clientset {
	kubeClient := fake.NewSimpleClientset()
	return kubeClient
}

func setupKubernetes(kc kubernetes.Interface) error {
	svc := core.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sanity-service",
			Namespace: testNamespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		Secrets: []core.ObjectReference{
			{
				Name:      "sanity-service-secret",
				Namespace: testNamespace,
			},
		},
	}
	_, err := kc.CoreV1().ServiceAccounts(testNamespace).Create(&svc)
	if err != nil {
		return err
	}

	secret := core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sanity-service-secret",
			Namespace: testNamespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		Data: map[string][]byte{
			"token": []byte("sanity-token"),
		},
	}
	if _, err = kc.CoreV1().Secrets(testNamespace).Create(&secret); err != nil {
		return err
	}

	return nil
}

func getAppBindingWithFakeClient(vaultUrl string) (appcat_cs.AppcatalogV1alpha1Interface, error) {
	data := `{
      "apiVersion": "kubevault.com/v1alpha1",
      "kind": "VaultServerConfiguration",
      "usePodServiceAccountForCSIDriver": true,
      "authPath": "kubernetes",
	  "policyControllerRole": "testrole"
    }`

	app := cr.AppBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sanity-app",
			Namespace: "default",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "appcatalog.appscode.com/v1alpha1",
			Kind:       "AppBinding",
		},
		Spec: cr.AppBindingSpec{
			ClientConfig: cr.ClientConfig{
				URL:                   &vaultUrl,
				InsecureSkipTLSVerify: true,
			},
			Parameters: &runtime.RawExtension{
				Raw: []byte(data),
			},
		},
	}
	client := appFake.NewSimpleClientset().AppcatalogV1alpha1()
	_, err := client.AppBindings(testNamespace).Create(&app)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestRaw(t *testing.T) {
	data := `{
      "apiVersion": "kubevault.com/v1alpha1",
      "kind": "VaultServerConfiguration",
      "usePodServiceAccountForCSIDriver": "true",
      "authPath": "kubernetes"
    }`
	x, e := json.Marshal(data)
	fmt.Println(e)
	d := cr.AppBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "appcatalog.appscode.com/v1alpha1",
			Kind:       "AppBinding",
		},
		Spec: cr.AppBindingSpec{
			Parameters: &runtime.RawExtension{
				Raw: x,
			},
		},
	}

	y, e := json.Marshal(d)
	fmt.Println(string(y), e)
}
