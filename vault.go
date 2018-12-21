package vaultlib

import (
	"encoding/json"

	"github.com/pkg/errors"
)

//VaultResponse holds the generic json response from Vault server
type VaultResponse struct {
	RequestID     string          `json:"request_id"`
	LeaseID       string          `json:"lease_id"`
	Renewable     bool            `json:"renewable"`
	LeaseDuration int             `json:"lease_duration"`
	Data          json.RawMessage `json:"data"`
	WrapInfo      json.RawMessage `json:"wrap_info"`
	Warnings      json.RawMessage `json:"warnings"`
	Auth          json.RawMessage `json:"auth"`
}

//VaultMountResponse holds the Vault Mount list response (used to unmarshall the globa vault response)
type VaultMountResponse struct {
	Auth   json.RawMessage `json:"auth"`
	Secret json.RawMessage `json:"secret"`
}

// VaultSecretMounts hodls the vault secret engine def
type VaultSecretMounts struct {
	Name     string `json:"??,string"`
	Accessor string `json:"accessor"`
	Config   struct {
		DefaultLeaseTTL int    `json:"default_lease_ttl"`
		ForceNoCache    bool   `json:"force_no_cache"`
		MaxLeaseTTL     int    `json:"max_lease_ttl"`
		PluginName      string `json:"plugin_name"`
	} `json:"config"`
	Description string                 `json:"description"`
	Local       bool                   `json:"local"`
	Options     map[string]interface{} `json:"options"`
	SealWrap    bool                   `json:"seal_wrap"`
	Type        string                 `json:"type"`
}

func (c *VaultClient) getKVVersion(kvName string) (version string, err error) {
	req := new(request)
	req.Method = "GET"
	req.URL = c.Address
	req.URL.Path = "/v1/sys/internal/ui/mounts"
	req.Token = c.Token
	err = req.prepareRequest()
	if err != nil {
		return "", err
	}

	rsp, err := req.execute(c.HTTPClient)
	if err != nil {
		return "", errors.Wrap(errors.WithStack(err), errInfo())
	}
	var mountResponse VaultMountResponse
	var vaultSecretMount = make(map[string]VaultSecretMounts)
	jsonErr := json.Unmarshal([]byte(rsp.Data), &mountResponse)
	if jsonErr != nil {
		return "", errors.Wrap(errors.WithStack(err), errInfo())
	}

	jsonErr = json.Unmarshal([]byte(mountResponse.Secret), &vaultSecretMount)
	if jsonErr != nil {
		return "", errors.Wrap(errors.WithStack(err), errInfo())
	}

	for _, v := range vaultSecretMount {
		if v.Name == kvName {
			if len(v.Options) > 0 {
				switch v.Options["version"].(type) {
				case string:
					version = v.Options["version"].(string)
				default:
					version = "1"
				}
			} else {
				//kv v1
				version = "1"
			}
		}
	}
	if len(version) == 0 {
		return "", errors.New("Could not get kv version")
	}
	return version, nil

}
