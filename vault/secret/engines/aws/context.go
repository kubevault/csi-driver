package aws

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

}

var _ SecretEngine = &EngineInfo{}

const(
	UID = "AWS"
)
func init()  {
	RegisterSecretEngine(UID, func(ctx context.Context) (SecretEngine, error) {
		return New(ctx), nil
	})
}

func New(ctx context.Context) SecretEngine  {
	return &EngineInfo{ctx:ctx}
}


