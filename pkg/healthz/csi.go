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
package healthz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	connlib "github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/rpc"
	"google.golang.org/grpc"
	lib "k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/klog"
)

// healthProbe implements the simplest possible healthz checker.
type healthProbe struct {
	conn         *grpc.ClientConn
	driverName   string
	probeTimeout time.Duration
}

func NewCSIProbe(csiAddress string, probeTimeout time.Duration) (lib.HealthChecker, error) {
	csiConn, err := connlib.Connect(csiAddress)
	if err != nil {
		// connlib should retry forever so a returned error should mean
		// the grpc client is misconfigured rather than an error on the network
		return nil, fmt.Errorf("failed to establish connection to CSI driver: %v", err)
	}

	klog.Infof("calling CSI driver to discover driver name")
	csiDriverName, err := rpc.GetDriverName(context.Background(), csiConn)
	if err != nil {
		return nil, fmt.Errorf("failed to get CSI driver name: %v", err)
	}
	klog.Infof("CSI driver name: %q", csiDriverName)

	return &healthProbe{
		conn:         csiConn,
		driverName:   csiDriverName,
		probeTimeout: probeTimeout,
	}, nil
}

func (h healthProbe) Name() string {
	return "csi-healthProbe"
}

func (h healthProbe) Check(req *http.Request) error {
	ctx, cancel := context.WithTimeout(req.Context(), h.probeTimeout)
	defer cancel()

	klog.Infof("Sending probe request to CSI driver %q", h.driverName)
	ready, err := rpc.Probe(ctx, h.conn)
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}

	if !ready {
		return errors.New("driver responded but is not ready")
	}

	klog.Infof("Health check succeeded")
	return nil
}
