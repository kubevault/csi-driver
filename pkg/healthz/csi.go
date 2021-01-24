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

package healthz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	connlib "github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/metrics"
	"github.com/kubernetes-csi/csi-lib-utils/rpc"
	"google.golang.org/grpc"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	lib "k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/klog"
)

// healthProbe implements the simplest possible healthz checker.
type healthProbe struct {
	addr         string
	probeTimeout time.Duration

	conn       *grpc.ClientConn
	driverName string
	once       sync.Once
}

func NewCSIProbe(csiAddress string, probeTimeout time.Duration) (lib.HealthChecker, error) {
	u, err := url.Parse(csiAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to parse address: %q", err)
	}
	if u.Scheme != "unix" {
		return nil, fmt.Errorf("%v endpoint scheme not supported", u.Scheme)
	}

	return &healthProbe{
		addr:         u.Path,
		probeTimeout: probeTimeout,
	}, nil
}

func (h *healthProbe) Name() string {
	return "csi-vault-health-probe"
}

func (h *healthProbe) init() error {
	csiConn, err := connlib.Connect(h.addr, metrics.NewCSIMetricsManager("secrets.csi.kubevault.com"))
	if err != nil {
		// connlib should retry forever so a returned error should mean
		// the grpc client is misconfigured rather than an error on the network
		return fmt.Errorf("failed to establish connection to CSI driver: %v", err)
	}

	klog.Infof("calling CSI driver to discover driver name")
	csiDriverName, err := rpc.GetDriverName(context.Background(), csiConn)
	if err != nil {
		return fmt.Errorf("failed to get CSI driver name: %v", err)
	}
	klog.Infof("CSI driver name: %q", csiDriverName)

	h.conn = csiConn
	h.driverName = csiDriverName
	return nil
}

func (h *healthProbe) Check(req *http.Request) error {
	h.once.Do(func() {
		utilruntime.Must(h.init())
	})

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
