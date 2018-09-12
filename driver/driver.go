package driver

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"

	csi "github.com/container-storage-interface/spec/lib/go/csi/v0"
	vaultapi "github.com/hashicorp/vault/api"
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
	region   string

	srv         *grpc.Server
	vaultClient *vaultapi.Client
	mounter     Mounter
	log         *logrus.Entry
}

func NewDriver(ep, token string) (*Driver, error) {
	config := vaultapi.DefaultConfig()
	config.Address = ep
	if err := config.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("Cannot get config from env: %v", err)
	}

	// By default this added the system's CAs
	err := config.ConfigureTLS(&vaultapi.TLSConfig{Insecure: false})
	if err != nil {
		return nil, fmt.Errorf("Failed to configureTLS: %v", err)
	}

	// Create the client
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %s", err)
	}
	client.SetToken(token)

	// The generator token is periodic so we can set the increment to 0
	// and it will default to the period.
	if _, err = client.Auth().Token().RenewSelf(0); err != nil {
		return nil, fmt.Errorf("Couldn't renew generator token: %v", err)
	}
	return &Driver{
		endpoint:    ep,
		vaultClient: client,
	}, nil

}

// Run starts the CSI plugin by communication over the given endpoint
func (d *Driver) Run() error {
	u, err := url.Parse(d.endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse address: %q", err)
	}

	addr := path.Join(u.Host, filepath.FromSlash(u.Path))
	if u.Host == "" {
		addr = filepath.FromSlash(u.Path)
	}

	// CSI plugins talk only over UNIX sockets currently
	if u.Scheme != "unix" {
		return fmt.Errorf("currently only unix domain sockets are supported, have: %s", u.Scheme)
	} else {
		// remove the socket if it's already there. This can happen if we
		// deploy a new version and the socket was created from the old running
		// plugin.
		d.log.WithField("socket", addr).Info("removing socket")
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove unix domain socket file %s, error: %s", addr, err)
		}
	}

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

	d.srv = grpc.NewServer(grpc.UnaryInterceptor(errHandler))
	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)

	d.log.WithField("addr", addr).Info("server started")
	return d.srv.Serve(listener)
	return nil
}

// Stop stops the plugin
func (d *Driver) Stop() {
	d.log.Info("server stopped")
	d.srv.Stop()
}
