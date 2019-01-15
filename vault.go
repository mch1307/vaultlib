package vaultlib

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// Secret holds the secret.
//
// KV contains data in case of KV secret.
//
// JSONSecret contains data in case of JSON raw secret.
type Secret struct {
	KV         map[string]string
	JSONSecret json.RawMessage
}

type rawSecretData struct {
	Data json.RawMessage
}

// GetSecret returns the Vault secret object
//
// KV: map[string]string if the secret is a KV
//
// JSONSecret: json.RawMessage if the secret is a json
func (c *Client) GetSecret(path string) (secret Secret, err error) {
	var v2Secret vaultSecretKV2
	var vaultRsp rawSecretData
	secret.KV = make(map[string]string)

	kvVersion, kvName, err := c.getKVInfo(path)
	if err != nil {
		return secret, errors.Wrap(errors.WithStack(err), errInfo())
	}
	url := c.address.String()

	if kvVersion == "2" {
		url = url + "/v1/" + kvName + "data/" + strings.TrimPrefix(path, kvName)
	} else {
		url = url + "/v1/" + path
	}

	req, _ := c.newRequest("GET", url)

	rsp, err := req.execute()
	if err != nil {
		return secret, errors.Wrap(errors.WithStack(err), errInfo())
	}

	if kvVersion == "2" {
		err = json.Unmarshal([]byte(rsp.Data), &v2Secret)
		if err != nil {
			return secret, errors.Wrap(errors.WithStack(err), errInfo())
		}
		for k, v := range v2Secret.Data {
			switch t := v.(type) {
			case string:
				secret.KV[k] = t
			case interface{}:
				//secret.JSONSecret = rsp.Data
				//Parse twice to remove
				err = json.Unmarshal([]byte(rsp.Data), &vaultRsp)
				if err != nil {
					return secret, err
				}
				err := json.Unmarshal([]byte(vaultRsp.Data), &secret.JSONSecret)
				if err != nil {
					return secret, err
				}
				return secret, err
			}
		}
	} else if kvVersion == "1" {
		raw := make(map[string]interface{})
		err = json.Unmarshal([]byte(rsp.Data), &raw)
		if err != nil {
			return secret, errors.Wrap(errors.WithStack(err), errInfo())
		}
		for k, v := range raw {
			switch t := v.(type) {
			case string:
				secret.KV[k] = t
			case interface{}:
				secret.JSONSecret = rsp.Data
				return secret, err
			}
		}
	}
	return secret, nil
}

// vaultMountResponse holds the Vault Mount list response (used to unmarshall the global vault response)
type vaultMountResponse struct {
	Auth   json.RawMessage `json:"auth"`
	Secret json.RawMessage `json:"secret"`
}

// vaultSecretMounts holds the vault secret engine def
type vaultSecretMounts struct {
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

func (c *Client) getKVInfo(path string) (version, name string, err error) {
	var mountResponse vaultMountResponse
	var vaultSecretMount = make(map[string]vaultSecretMounts)
	url := c.address.String() + "/v1/sys/internal/ui/mounts"

	req, _ := c.newRequest("GET", url)

	rsp, err := req.execute()
	if err != nil {
		return "", "", errors.Wrap(errors.WithStack(err), errInfo())
	}

	jsonErr := json.Unmarshal([]byte(rsp.Data), &mountResponse)
	if jsonErr != nil {
		return "", "", errors.Wrap(errors.WithStack(err), errInfo())
	}

	jsonErr = json.Unmarshal([]byte(mountResponse.Secret), &vaultSecretMount)
	if jsonErr != nil {
		return "", "", errors.Wrap(errors.WithStack(err), errInfo())
	}

	for kvName, v := range vaultSecretMount {
		if strings.HasPrefix(path, kvName) {
			name = kvName
			if len(v.Options) > 0 {
				switch v.Options["version"].(type) {
				case string:
					version = v.Options["version"].(string)
				}
			} else {
				//kv v1
				version = "1"
			}
			break
		}
	}
	if len(version) == 0 {
		return "", "", errors.New("Could not get kv version")
	}
	return version, name, nil

}
