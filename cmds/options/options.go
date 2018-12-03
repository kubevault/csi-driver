package options

import (
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	Endpoint string
	NodeName string
}

func NewConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		Endpoint: "unix:/tmp/csi.sock",
		NodeName: hostname,
	}
}

func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Endpoint, "endpoint", c.Endpoint, "CSI endpoint")
	fs.StringVar(&c.NodeName, "node", c.NodeName, "Hostname")

}
