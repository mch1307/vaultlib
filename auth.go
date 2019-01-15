package vaultlib

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

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
func (c *Client) renewToken() {
	var vaultData vaultAuth
	jsonToken := make(map[string]string)

	for {
		duration := c.token.TTL - 2
		time.Sleep(time.Second * time.Duration(duration))

		url := c.address.String() + "/v1/auth/token/renew"

		jsonToken["token"] = c.getTokenID()

		req, _ := c.newRequest("POST", url)

		req.setJSONBody(jsonToken)

		resp, err := req.execute()
		if err != nil {
			c.setStatus("Error renewing token " + err.Error())
			continue
		}

		jsonErr := json.Unmarshal([]byte(resp.Auth), &vaultData)
		if jsonErr != nil {
			c.setStatus("Error renewing token " + err.Error())
			continue
		}

		if err := c.setTokenInfo(); err != nil {
			c.setStatus("Error renewing token " + err.Error())
			continue
		}
		c.setStatus("token renewed")
	}
}

// setTokenFromAppRole get the token from Vault and set it in the client
func (c *Client) setTokenFromAppRole() error {
	var vaultData vaultAuth
	if c.config.AppRoleCredentials.RoleID == "" {
		return errors.New("No credentials provided")
	}

	url := c.address.String() + "/v1/auth/approle/login"

	req, _ := c.newRequest("POST", url)

	req.setJSONBody(c.config.AppRoleCredentials)

	resp, err := req.execute()
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}

	jsonErr := json.Unmarshal([]byte(resp.Auth), &vaultData)
	if jsonErr != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	c.withLockContext(func() {
		c.token.ID = vaultData.ClientToken
	})

	if err = c.setTokenInfo(); err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	if c.token.Renewable {
		go c.renewToken()
	}

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

func (c *Client) setTokenInfo() error {
	url := c.address.String() + "/v1/auth/token/lookup-self"
	var tokenInfo VaultTokenInfo

	req, _ := c.newRequest("GET", url)

	res, err := req.execute()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res.Data, &tokenInfo); err != nil {
		return err
	}
	c.withLockContext(func() {
		c.token = &tokenInfo
		c.isAuthenticated = true

	})
	return nil
}
