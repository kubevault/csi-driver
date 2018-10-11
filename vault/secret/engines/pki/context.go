package pki

import (
	vaultapi "github.com/hashicorp/vault/api"
	. "github.com/kubevault/csi-driver/vault/secret"
	"context"
	"time"
)
type EngineInfo struct {
	ctx context.Context
	vc  *vaultapi.Client

	secretName string
	secretDir string

	certificate *certificate
	stopCh chan struct{}
	renewTime time.Duration
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


