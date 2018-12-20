package vaultlib

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// VaultAuth holds the Vault Auth response from server
type VaultAuth struct {
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

// VaultSecretMount dfbnsfdbkjsdn
type VaultSecretMount struct {
	Auth   json.RawMessage `json:"auth"`
	Secret struct {
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
	} `json:"secret"`
}

//SetTokenFromAppRole get the token from Vault and set it in the client
func (c *VaultClient) SetTokenFromAppRole() error {
	if c.Config.AppRoleCredentials.RoleID == "" {
		return errors.New("No credentials provided")
	}

	req := new(request)
	req.Method = "POST"
	req.URL = c.Address
	req.URL.Path = "/v1/auth/approle/login"
	err := req.prepareRequest()
	if err != nil {
		return err
	}
	req.setJSONBody(c.Config.AppRoleCredentials)
	resp, err := req.execute(c.HTTPClient)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	var vaultData VaultSecretMount
	jsonErr := json.Unmarshal([]byte(resp.Data), &vaultData)
	if jsonErr != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}

	return nil
}

func (c *VaultClient) getKVVersion(kvName string) error {

	req := new(request)
	req.Method = "GET"
	req.URL = c.Address
	req.URL.Path = "/v1/sys/internal/ui/mounts"
	err := req.prepareRequest()
	if err != nil {
		return err
	}

	rsp, err := req.execute(c.HTTPClient)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}

	var vaultAuth VaultAuth
	jsonErr := json.Unmarshal([]byte(rsp.Auth), &vaultAuth)
	if jsonErr != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	var mountInfo VaultSecretMounts = make(map[string]*VaultSecretMount)

	err = mountInfo.UnmarshalJSON([]byte(vaultRsp.Data.Secret))
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	var version string
	for _, v := range mountInfo {
		if v.Name == name {
			log.Debugf("Selected Vault mount: %+v", v)
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
	a.VaultKvVersion = version
	if len(version) == 0 {
		return errors.New("Could not get kv version")
	}
	log.Debugf("%v is vault kv v %v", name, version)
	return nil
}
