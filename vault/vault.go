package vault

import (
	vaultapi "github.com/hashicorp/vault/api"
)



func NewVaultClient(url, token string, tlsConfig *vaultapi.TLSConfig) (*vaultapi.Client, error) {
	cfg := vaultapi.DefaultConfig()
	if url != "" {
		cfg.Address = url
	}
	if tlsConfig == nil {
		tlsConfig = &vaultapi.TLSConfig{Insecure: false}
	}
	cfg.ConfigureTLS(tlsConfig)
	vc, err := vaultapi.NewClient(cfg)
	vc.SetToken(token)

	if err != nil {
		return nil, err
	}
	return vc, nil
}

/*
func (c *Client) GetPolicyToken(policies []string, unwrap bool) (string, error)  {

	if unwrap {
		// We override the default WrappingLookupFunction which honors the VAULT_WRAP_TTL env variable
		c.Vc.SetWrappingLookupFunc(func(_, _ string) string { return "" })
	}

	metadata := map[string]string{
		//"creator": "csi-driver",
	}
	req := vaultapi.TokenCreateRequest{
		Policies: policies,
		Metadata: metadata,
		Period: "24h",
	}

	secret, err := c.Vc.Auth().Token().Create(&req)
	if err != nil {
		return "", errors.Errorf("Couldn't create scoped token for policies %v : %v", req.Policies, err)
	}
	if secret == nil {
		return "", errors.Errorf("Got nil secret when getting token")
	}

	if unwrap {
		return secret.Auth.ClientToken, nil
	}
	if secret.WrapInfo == nil {
		return "",  errors.Errorf("got unwrapped token ! Set VAULT_WRAP_TTL in kubelet environment")
	}

	return secret.WrapInfo.Token, nil
}
*/