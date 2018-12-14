package driver

import (
	"context"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

const (
	driverName        = "com.kubevault.csi.secrets"
	vendorVersion     = "0.1.1"
	podName           = "csi.storage.k8s.io/pod.name"
	podNamespace      = "csi.storage.k8s.io/pod.namespace"
	podUID            = "csi.storage.k8s.io/pod.uid"
	podServiceAccount = "csi.storage.k8s.io/serviceAccount.name"

	TestEnvForCSIDriver = "VAULT_CSI_TEST"
)

// Driver implements the following CSI interfaces:
//
//   csi.IdentityServer
//   csi.ControllerServer
//   csi.NodeServer
//
type Driver struct {
	endpoint string
	nodeId   string

	kubeClient kubernetes.Interface
	appClient  appcat_cs.AppcatalogV1alpha1Interface

	srv         *grpc.Server
	vaultClient *vaultapi.Client
	mounter     Mounter
	log         *logrus.Entry

	ch map[string]*vaultapi.Renewer
}

func NewDriver(ep, node string) (*Driver, error) {
	var kubeClient kubernetes.Interface
	var appClient appcat_cs.AppcatalogV1alpha1Interface

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	if kubeClient, err = kubernetes.NewForConfig(config); err != nil {
		return nil, err
	}

	if appClient, err = appcat_cs.NewForConfig(config); err != nil {
		return nil, err
	}

	return &Driver{
		endpoint: ep,
		mounter:  &mounter{},
		nodeId:   node,

		kubeClient: kubeClient,
		appClient:  appClient,

		log: logrus.New().WithFields(logrus.Fields{
			"node-id": node,
		}),
		ch: make(map[string]*vaultapi.Renewer),
	}, nil

}

// Run starts the CSI plugin by communication over the given endpoint
func (d *Driver) Run() error {
	u, err := url.Parse(d.endpoint)
	if err != nil {
		return errors.Errorf("unable to parse address: %q", err)
	}

	addr := path.Join(u.Host, filepath.FromSlash(u.Path))
	if u.Host == "" {
		addr = filepath.FromSlash(u.Path)
	}

	// CSI plugins talk only over UNIX sockets currently
	if u.Scheme != "unix" {
		return errors.Errorf("currently only unix domain sockets are supported, have: %s", u.Scheme)
	} else {
		// remove the socket if it's already there. This can happen if we
		// deploy a new version and the socket was created from the old running
		// plugin.
		d.log.WithField("socket", addr).Info("removing socket")
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			return errors.Errorf("failed to remove unix domain socket file %s, error: %s", addr, err)
		}
	}

	listener, err := net.Listen(u.Scheme, addr)
	if err != nil {
		return errors.Errorf("failed to listen: %v", err)
	}

	// log response errors for better observability
	errHandler := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			d.log.WithError(err).WithField("method", info.FullMethod).Error("method failed")
		}
		return resp, err
	}

	d.srv = grpc.NewServer(grpc.UnaryInterceptor(errHandler))
	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)

	d.log.WithField("addr", addr).Info("server started")
	return d.srv.Serve(listener)
}

// Stop stops the plugin
func (d *Driver) Stop() {
	d.log.Info("server stopped")
	d.srv.Stop()
}
