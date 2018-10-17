package auth

import (
	"github.com/kubevault/operator/apis"
	awsauth "github.com/kubevault/operator/pkg/vault/auth/aws"
	certauth "github.com/kubevault/operator/pkg/vault/auth/cert"
	k8sauth "github.com/kubevault/operator/pkg/vault/auth/kubernetes"
	tokenauth "github.com/kubevault/operator/pkg/vault/auth/token"
	basicauth "github.com/kubevault/operator/pkg/vault/auth/userpass"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

type AuthInterface interface {
	// successful login will return token
	// unsuccessful login will return err
	Login() (string, error)
}

func NewAuth(kc kubernetes.Interface, vApp *appcat.AppBinding) (AuthInterface, error) {
	if vApp == nil {
		return nil, errors.New("vault AppBinding is not provided")
	}
	if vApp.Spec.Secret == nil {
		return nil, errors.New("secret is not provided")
	}

	secret, err := kc.CoreV1().Secrets(vApp.Namespace).Get(vApp.Spec.Secret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	switch secret.Type {
	case core.SecretTypeBasicAuth:
		return basicauth.New(vApp, secret)
	case core.SecretTypeTLS:
		return certauth.New(vApp, secret)
	case core.SecretTypeServiceAccountToken:
		return k8sauth.New(vApp, secret)
	case apis.SecretTypeTokenAuth:
		return tokenauth.New(secret)
	case apis.SecretTypeAWSAuth:
		return awsauth.New(vApp, secret)
	default:
		return nil, errors.New("Invalid/Unsupported secret type")
	}
}
