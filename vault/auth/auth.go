package auth

import (
	"github.com/golang/glog"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"sync"
)

type PodInfo struct {
	Name           string
	Namespace      string
	UID            string
	ServiceAccount string
}

type Authentication interface {
	GetLoginToken() (string, error)
	SetRole(role string)
	SetVaultUrl(url string)
}

type Factory func(info PodInfo, client *vaultapi.Client) (Authentication, error)

var (
	authMutex   sync.Mutex
	authMethods = make(map[string]Factory)
)

func RegisterAuthMethod(name string, method Factory) {
	authMutex.Lock()
	defer authMutex.Unlock()
	if _, found := authMethods[name]; found {
		glog.Fatalf("Auth method %s was registered twice", name)
	}
	authMethods[name] = method
}

func GetAuthMethod(name string, info PodInfo, client *vaultapi.Client) (Authentication, error) {
	authMutex.Lock()
	defer authMutex.Unlock()
	f, found := authMethods[name]
	if !found {
		return nil, errors.Errorf("%s auth engine not found", name)
	}
	return f(info, client)
}
