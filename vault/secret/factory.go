package secret

import (
	"context"
	"sync"
	"github.com/golang/glog"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

type SecretEngine interface {
	InitializeEngine(vc *vaultapi.Client, opts map[string]string) error
	ReadSecret() error
	RenewSecret(vol string) error
	StopSync()

}

type Factory func(ctx context.Context) (SecretEngine, error)

var (
	engineMutex sync.Mutex
	secretEngines = make(map[string]Factory)
)

func RegisterSecretEngine(name string, engine Factory)  {
	engineMutex.Lock()
	defer engineMutex.Unlock()

	if _, found := secretEngines[name]; found {
		glog.Fatalf("Secret engine %s was registered twice", name)
	}
	secretEngines[name] = engine
}

func GetSecretEngine(name string, ctx context.Context) (SecretEngine, error)  {
	engineMutex.Lock()
	defer engineMutex.Unlock()

	f, found := secretEngines[name]
	if !found{
		return nil, errors.Errorf("No secret engine found with name %s", name)
	}
	return f(ctx)
}