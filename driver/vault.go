package driver

import (
	"encoding/json"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

type Client struct {
	vc   *vaultapi.Client
	role string
}

const defaultRoleName = "nginx"

func NewVaultClient(url, token string, tlsConfig *vaultapi.TLSConfig) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	if url != "" {
		cfg.Address = url
	}
	cfg.ConfigureTLS(tlsConfig)
	vc, err := vaultapi.NewClient(cfg)
	vc.SetToken(token)

	if err != nil {
		return nil, err
	}
	return &Client{vc: vc, role: defaultRoleName}, nil
}

// Gets token data
func (c *Client) GetTokenData(policies []string, poduid string, unwrap bool) (string, []byte, error) {

	if unwrap {
		// We override the default WrappingLookupFunction which honors the VAULT_WRAP_TTL env variable
		c.vc.SetWrappingLookupFunc(func(_, _ string) string { return "" })
	}

	secret, err := c.getTokenForPolicy(policies, poduid)
	if err != nil {
		return "", []byte{}, err
	}
	if secret == nil {
		return "", []byte{}, errors.Errorf("Got nil secret when getting token")
	}

	if unwrap {
		metadata, err := json.Marshal(secret)
		if err != nil {
			return "", []byte{}, errors.Errorf("Cloudn't marshall metadata: %v", err)
		}
		return secret.Auth.ClientToken, metadata, nil
	}
	// else we want a wrapped token :
	if secret.WrapInfo == nil {
		return "", []byte{}, errors.Errorf("got unwrapped token ! Set VAULT_WRAP_TTL in kubelet environment")
	}

	metadata, err := json.Marshal(secret.WrapInfo)
	if err != nil {
		return "", []byte{}, errors.Errorf("Couldn't marshal vault response: %v", err)
	}
	return secret.WrapInfo.Token, metadata, nil
}

// GetTokenForPolicy gets a wrapped token from Vault scoped with given policy
func (c *Client) getTokenForPolicy(policies []string, poduid string) (*vaultapi.Secret, error) {

	metadata := map[string]string{
		"poduid":  poduid,
		"creator": "csi-driver",
	}
	req := vaultapi.TokenCreateRequest{
		Policies: policies,
		Metadata: metadata,
	}

	secret, err := c.vc.Auth().Token().Create(&req)
	if err != nil {
		return nil, errors.Errorf("Couldn't create scoped token for policies %v : %v", req.Policies, err)
	}
	return secret, nil

}
