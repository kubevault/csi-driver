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
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

const (
	driverName        = "secrets.csi.kubevault.com"
	vendorVersion     = "1.0.0"
	podName           = "csi.storage.k8s.io/pod.name"
	podNamespace      = "csi.storage.k8s.io/pod.namespace"
	podUID            = "csi.storage.k8s.io/pod.uid"
	podServiceAccount = "csi.storage.k8s.io/serviceAccount.name"
)

// Driver implements the following CSI interfaces:
//
//   csi.IdentityServer
//   csi.ControllerServer
//   csi.NodeServer
//
type Driver struct {
	config

	kubeClient kubernetes.Interface
	appClient  appcat_cs.AppcatalogV1alpha1Interface

	srv         *grpc.Server
	vaultClient *vaultapi.Client
	mounter     Mounter
	log         *logrus.Entry

	ch map[string]*vaultapi.Renewer

	// ready defines whether the driver is ready to function. This value will
	// be used by the `Identity` service via the `Probe()` method.
	readyMu sync.Mutex // protects ready
	ready   bool
}

// Run starts the CSI plugin by communication over the given CSIAddress
func (d *Driver) Run() error {
	u, err := url.Parse(d.Endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse address: %q", err)
	}

	var addr string
	if u.Scheme == "unix" {
		addr = u.Path
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s, error: %s", addr, err.Error())
		}

		listenDir := filepath.Dir(addr)
		if _, err := os.Stat(listenDir); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("expected Kubelet plugin watcher to create parent dir %s but did not find such a dir", listenDir)
			} else {
				return fmt.Errorf("failed to stat %s, error: %s", listenDir, err.Error())
			}
		}
	} else {
		return fmt.Errorf("%v endpoint scheme not supported", u.Scheme)
	}

	d.log.Infof("Start listening with scheme %v, addr %v", u.Scheme, addr)
	listener, err := net.Listen(u.Scheme, addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// log response errors for better observability
	errHandler := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			d.log.WithError(err).WithField("method", info.FullMethod).Error("method failed")
		}
		return resp, err
	}

	d.srv = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			errHandler,
		)),
	)
	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)
	d.ready = true
	grpc_prometheus.Register(d.srv)

	d.log.WithField("addr", addr).Info("server started")
	return d.srv.Serve(listener)
}

// Stop stops the plugin
func (d *Driver) Stop() {
	d.readyMu.Lock()
	d.ready = false
	d.readyMu.Unlock()

	d.log.Info("server stopped")
	d.srv.Stop()
}
