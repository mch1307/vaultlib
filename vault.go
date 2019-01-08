package vaultlib

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Secret holds the returned secret
// KV contains data in case of KV secret
// JSONSecret contains data in case of JSON raw secret
type Secret struct {
	KV         map[string]string
	JSONSecret json.RawMessage
}

type rawSecretData struct {
	Data json.RawMessage
}

// vaultResponse holds the generic json response from Vault server
type vaultResponse struct {
	RequestID     string          `json:"request_id"`
	LeaseID       string          `json:"lease_id"`
	Renewable     bool            `json:"renewable"`
	LeaseDuration int             `json:"lease_duration"`
	Data          json.RawMessage `json:"data"`
	WrapInfo      json.RawMessage `json:"wrap_info"`
	Warnings      json.RawMessage `json:"warnings"`
	Auth          json.RawMessage `json:"auth"`
}

// vaultMountResponse holds the Vault Mount list response (used to unmarshall the global vault response)
type vaultMountResponse struct {
	Auth   json.RawMessage `json:"auth"`
	Secret json.RawMessage `json:"secret"`
}

// vaultSecretMounts hodls the vault secret engine def
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

func (c *VaultClient) getKVInfo(path string) (version, name string, err error) {
	var mountResponse vaultMountResponse
	var vaultSecretMount = make(map[string]vaultSecretMounts)
	url := c.Address
	url.Path = "/v1/sys/internal/ui/mounts"

	req, err := newRequest("GET", c.Token, url)
	if err != nil {
		return "", "", errors.Wrap(errors.WithStack(err), errInfo())
	}

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
		}
	}
	if len(version) == 0 {
		return "", "", errors.New("Could not get kv version")
	}
	return version, name, nil

}

// vaultAuth holds the Vault Auth response from server
type vaultAuth struct {
	ClientToken string   `json:"client_token"`
	Accessor    string   `json:"accessor"`
	Policies    []string `json:"policies"`
	Metadata    struct {
		RoleName string `json:"role_name"`
	} `json:"metadata"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	EntityID      string `json:"entity_id"`
}

// renew the client's token, launched at client creation time as a go routine
func (c *VaultClient) renewToken() {
	var vaultData vaultAuth
	jsonToken := make(map[string]string)

	for {
		duration := c.Lease - 1
		time.Sleep(time.Second * time.Duration(duration))

		url := c.Address
		url.Path = "v1/auth/token/renew"
		jsonToken["token"] = c.Token

		req, err := newRequest("POST", c.Token, url)
		if err != nil {
			c.Status = "Error renewing token " + err.Error()
			continue
		}
		err = req.setJSONBody(jsonToken)
		if err != nil {
			c.Status = "Error renewing token " + err.Error()
			continue
		}
		resp, err := req.execute()
		if err != nil {
			c.Status = "Error renewing token " + err.Error()
			continue
		}

		jsonErr := json.Unmarshal([]byte(resp.Auth), &vaultData)
		if jsonErr != nil {
			c.Status = "Error renewing token " + err.Error()
			continue
		}
		c.Lease = vaultData.LeaseDuration
		c.Status = "Token renewed"
	}
}

// setTokenFromAppRole get the token from Vault and set it in the client
func (c *VaultClient) setTokenFromAppRole() error {
	var vaultData vaultAuth
	if c.Config.AppRoleCredentials.RoleID == "" {
		return errors.New("No credentials provided")
	}

	url := c.Address
	url.Path = "/v1/auth/approle/login"

	req, err := newRequest("POST", c.Token, url)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	err = req.setJSONBody(c.Config.AppRoleCredentials)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}

	resp, err := req.execute()
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}

	jsonErr := json.Unmarshal([]byte(resp.Auth), &vaultData)
	if jsonErr != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	if vaultData.Renewable {
		go c.renewToken()
	}
	c.Token = vaultData.ClientToken

	return nil
}

// vaultSecretKV2 holds the Vault secret (kv v2)
type vaultSecretKV2 struct {
	Data     map[string]interface{} `json:"data"`
	Metadata struct {
		CreatedTime  time.Time `json:"created_time"`
		DeletionTime string    `json:"deletion_time"`
		Destroyed    bool      `json:"destroyed"`
		Version      int       `json:"version"`
	} `json:"metadata"`
}

// GetVaultSecret returns the Vault secret object
// KV: map[string]string if the secret is a KV
// JSONSecret: json.RawMessage if the secret is a json
func (c *VaultClient) GetVaultSecret(path string) (secret Secret, err error) {
	var v2Secret vaultSecretKV2
	var vaultRsp rawSecretData
	secret.KV = make(map[string]string)

	kvVersion, kvName, err := c.getKVInfo(path)
	if err != nil {
		return secret, errors.Wrap(errors.WithStack(err), errInfo())
	}
	url := c.Address

	if kvVersion == "2" {
		url.Path = "/v1/" + kvName + "data/" + strings.TrimPrefix(path, kvName)
	} else {
		url.Path = "/v1/" + path
	}

	req, err := newRequest("GET", c.Token, url)
	if err != nil {
		return secret, errors.Wrap(errors.WithStack(err), errInfo())
	}

	rsp, err := req.execute()
	if err != nil {
		return secret, errors.Wrap(errors.WithStack(err), errInfo())
	}

	// parse to Vx and get a simple kv map back
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
