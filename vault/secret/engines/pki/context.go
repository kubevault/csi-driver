package pki

import (
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/kubevault/csi-driver/vault/secret"
	"context"
)
type EngineInfo struct {
	ctx context.Context
	vc  *vaultapi.Client

	secretName string
	secretDir string

	stopCh chan struct{}
}

var _ SecretEngine = &EngineInfo{}

const(
	UID = "PKI"
)
func init()  {
	RegisterSecretEngine(UID, func(ctx context.Context) (SecretEngine, error) {
		return New(ctx), nil
	})
}

func New(ctx context.Context) SecretEngine  {
	return &EngineInfo{ctx:ctx}
}


