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
package server

import (
	"flag"
	"os"
	"time"

	"kubevault.dev/csi-driver/pkg/driver"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

type ExtraOptions struct {
	Endpoint          string
	NodeName          string
	ConnectionTimeout time.Duration
	QPS               float64
	Burst             int
}

func NewExtraOptions() *ExtraOptions {
	hostname, _ := os.Hostname()
	return &ExtraOptions{
		Endpoint:          "unix:///var/lib/kubelet/plugins/com.kubevault.csi.secrets/csi.sock",
		NodeName:          hostname,
		ConnectionTimeout: 30 * time.Second,
		QPS:               100,
		Burst:             100,
	}
}

func (s *ExtraOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.Endpoint, "endpoint", s.Endpoint, "CSI endpoint")
	fs.StringVar(&s.NodeName, "node", s.NodeName, "Hostname")
	fs.DurationVar(&s.ConnectionTimeout, "connection-timeout", s.ConnectionTimeout, "Timeout for waiting for CSI driver socket in seconds.")

	fs.Float64Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
}

func (s *ExtraOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("vault-server", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *ExtraOptions) ApplyTo(cfg *driver.Config) error {
	var err error

	cfg.ClientConfig.QPS = float32(s.QPS)
	cfg.ClientConfig.Burst = s.Burst
	cfg.Endpoint = s.Endpoint
	cfg.NodeId = s.NodeName

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.AppClient, err = appcat_cs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	return nil
}
