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
	"github.com/kubevault/csi-driver/vault"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	driverName    = "com.vault.csi.vaultdbs"
	vendorVersion = "0.0.1"
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
	url      string

	srv         *grpc.Server
	vaultClient *vault.Client
	mounter     Mounter
	log         *logrus.Entry
}

func NewDriver(ep, url, node, token string) (*Driver, error) {
	// Create the client
	client, err := vault.NewVaultClient(url, token, &vaultapi.TLSConfig{Insecure: false})
	if err != nil {
		return nil, errors.Errorf("failed to create vault client: %s", err)
	}
	//client.vc.SetToken(token)

	// The generator token is periodic so we can set the increment to 0
	// and it will default to the period.
	/*if _, err = client.vc.Auth().Token().RenewSelf(0); err != nil {
		return nil, errors.Errorf("Couldn't renew generator token: %v", err)
	}*/
	return &Driver{
		endpoint:    ep,
		url:         url,
		vaultClient: client,
		mounter:     &mounter{url, token},
		nodeId:      node,
		log: logrus.New().WithFields(logrus.Fields{
			"node-id": node,
		}),
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
