package util

import (
	"fmt"
	"path/filepath"

	vaultapi "github.com/hashicorp/vault/api"
	"k8s.io/kubernetes/pkg/apis/core"
)

const (
	VaultContainerName         = "vault"
	VaultUnsealerContainerName = "vault-unsealer"
	VaultInitContainerName     = "vault-config"
	VaultExporterContainerName = "vault-exporter"
)

const (
	// VaultConfigFile is the file that vault pod uses to read config from
	VaultConfigFile = "/etc/vault/config/vault.hcl"

	// VaultTLSAssetDir is the dir where vault's server TLS sits
	VaultTLSAssetDir = "/etc/vault/tls/"
)

var listenerFmt = `
listener "tcp" {
  address = "0.0.0.0:8200"
  cluster_address = "0.0.0.0:8201"
  tls_cert_file = "%s"
  tls_key_file  = "%s"
}
`

// NewConfigWithDefaultParams appends to given config data some default params:
// - tcp listener
func NewConfigWithDefaultParams() string {
	return fmt.Sprintf(listenerFmt, filepath.Join(VaultTLSAssetDir, core.TLSCertKey), filepath.Join(VaultTLSAssetDir, core.TLSPrivateKeyKey))
}

// ListenerConfig creates tcp listener config
func GetListenerConfig() string {
	listenerCfg := fmt.Sprintf(listenerFmt,
		filepath.Join(VaultTLSAssetDir, core.TLSCertKey),
		filepath.Join(VaultTLSAssetDir, core.TLSPrivateKeyKey))

	return listenerCfg
}

func NewVaultClient(hostname string, port string, tlsConfig *vaultapi.TLSConfig) (*vaultapi.Client, error) {
	cfg := vaultapi.DefaultConfig()
	podURL := fmt.Sprintf("https://%s:%s", hostname, port)
	cfg.Address = podURL
	cfg.ConfigureTLS(tlsConfig)
	return vaultapi.NewClient(cfg)
}