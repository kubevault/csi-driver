package options

import (
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	Endpoint string
	Token    string
	Url      string
	NodeName string
}

func NewConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		Endpoint: "unix:///var/lib/kubelet/plugins/com.vault.csi.vaultdbs/csi.sock",
		Url:      "https://api.vault.com/",
		Token:    "",
		NodeName: hostname,
	}
}

func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Endpoint, "endpoint", c.Endpoint, "CSI endpoint")
	fs.StringVar(&c.Token, "token", c.Token, "Vault access token")
	fs.StringVar(&c.Url, "url", c.Url, "Vault API URL")
	fs.StringVar(&c.NodeName, "node", c.NodeName, "Linode Hostname")

}
