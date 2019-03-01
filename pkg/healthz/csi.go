package healthz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/kubernetes-csi/livenessprobe/pkg/connection"
	lib "k8s.io/apiserver/pkg/server/healthz"
)

// probe implements the simplest possible healthz checker.
type probe struct {
	csiAddress        string
	connectionTimeout time.Duration
}

func NewCSIProbe(csiAddress string, connectionTimeout time.Duration) lib.HealthzChecker {
	return &probe{
		csiAddress:        csiAddress,
		connectionTimeout: connectionTimeout,
	}
}

func (probe) Name() string {
	return "csi-probe"
}

// CSIProbe is a health check that returns true if CSI probe is successful.
func (p probe) Check(_ *http.Request) error {
	csiConn, err := p.getCSIConnection()
	if err != nil {
		return fmt.Errorf("Failed to get connection to CSI  with error: %v.", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.connectionTimeout)
	defer cancel()
	return p.runProbe(ctx, csiConn)
}

// ref: https://github.com/kubernetes-csi/livenessprobe/blob/v1.0.1/cmd/main.go#L44

func (p probe) getCSIConnection() (connection.CSIConnection, error) {
	// Connect to CSI.
	glog.Infof("Attempting to open a gRPC connection with: %s", p.csiAddress)
	csiConn, err := connection.NewConnection(p.csiAddress, p.connectionTimeout)
	if err != nil {
		return nil, err
	}
	return csiConn, nil
}

func (p probe) runProbe(ctx context.Context, csiConn connection.CSIConnection) error {
	// Get CSI driver name.
	glog.Infof("Calling CSI driver to discover driver name.")
	csiDriverName, err := csiConn.GetDriverName(ctx)
	if err != nil {
		return err
	}
	glog.Infof("CSI driver name: %q", csiDriverName)

	// Sending Probe request
	glog.Infof("Sending probe request to CSI driver.")
	if err := csiConn.LivenessProbe(ctx); err != nil {
		return err
	}
	return nil
}
