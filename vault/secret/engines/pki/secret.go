package pki

import (
	vaultapi "github.com/hashicorp/vault/api"
	"fmt"
	"os"
	"github.com/pkg/errors"
	"path"
	"io/ioutil"
	"strings"
	"time"
	"encoding/json"
	"context"
	"github.com/hashicorp/vault/helper/parseutil"
)

type certificate struct {
	CommonName string `json:"common_name,omitempty"`
	AltName string `json:"alt_name,omitempty"`
	IpSans string `json:"ip_sans,omitempty"`
	UriSans string `json:"uri_sans,omitempty"`
	OtherSans string `json:"other_sans,omitempty"`
	Ttl string `json:"ttl,omitempty"`
	Format string `json:"format,omitempty"`
	PrivateKeyFormat string `json:"private_key_format,omitempty"`
	ExcludeCnFromSans bool `json:"exclude_cn_from_sans,omitempty"`
}


func (se *EngineInfo) InitializeEngine(vc *vaultapi.Client, opts map[string]string) error {
	se.vc = vc
	se.secretName = opts["secretName"]
	se.secretDir = opts["targetDir"]
	se.stopCh = make(chan struct{})

	data, err := json.Marshal(opts)
	if err != nil {
		return err
	}
	cert := &certificate{}
	if err = json.Unmarshal(data, cert); err != nil {
		return err
	}
	se.certificate = cert

	if cert.Ttl != "" {
		if se.renewTime, err = time.ParseDuration(cert.Ttl); err != nil {
			return err
		}
	} else {
		if se.renewTime, err = se.getRoleTTL(); err != nil {
			return  err
		}
	}
	se.renewTime = se.renewTime - time.Minute

	return nil
}


func (se *EngineInfo) ReadSecret() error {
	path := fmt.Sprintf("/v1/pki/issue/%s", se.secretName)

	r := se.vc.NewRequest("POST", path)
	if err := r.SetJSONBody(se.certificate); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(se.ctx)
	defer cancelFunc()
	resp, err := se.vc.RawRequestWithContext(ctx, r)
	if err != nil {
		return  err
	}
	defer resp.Body.Close()

	secret, err := vaultapi.ParseSecret(resp.Body)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(se.secretDir, 0755); err != nil {
		return errors.Errorf("Failed to mkdir %v: %v", se.secretDir, err)
	}

	for key, val := range secret.Data {
		if val != nil {
			if err := writeData(key, val.(string), se.secretDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeData(key, value, dir string) error {
	keyPath := path.Join(dir, key)
	return ioutil.WriteFile(keyPath, []byte(strings.TrimSpace(value)), 0644)
}

func (se *EngineInfo) RenewSecret(vol string) error {
	for {
		select {
		case <-se.stopCh:
			return nil
		case <-time.After(se.renewTime):
			fmt.Println("Renewing pki data...")
			if err := se.ReadSecret(); err != nil {
				return err
			}

		}
	}
}

func (se *EngineInfo) StopSync()  {
	close(se.stopCh)
}

func (se *EngineInfo) getRoleTTL() (time.Duration, error) {
	path := fmt.Sprintf("/v1/pki/roles/%s", se.secretName)
	req := se.vc.NewRequest("GET", path)
	resp, err := se.vc.RawRequest(req)
	if err != nil {
		return 0, err
	}
	secret, err := vaultapi.ParseSecret(resp.Body)
	if err != nil {
		return 0, err
	}
	ttl, err := parseutil.ParseDurationSecond(secret.Data["ttl"])
	if err != nil {
		return 0, err
	}
	if ttl == 0{
		return parseutil.ParseDurationSecond(secret.Data["max_ttl"])
	}
	return ttl, nil
}