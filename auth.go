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
		duration := c.Token.TTL - 2
		time.Sleep(time.Second * time.Duration(duration))
		url := c.Address
		url.Path = "v1/auth/token/renew"

		jsonToken["token"] = c.Token.ID

		req, err := newRequest("POST", c.Token.ID, url)
		if err != nil {
			c.Status = "Error renewing token " + err.Error()
			continue
		}
		req.setJSONBody(jsonToken)

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

		if err := c.setTokenInfo(); err != nil {
			c.Status = "Error renewing token " + err.Error()
			continue
		}
		c.Lock()
		c.Status = "Token renewed"
		c.Unlock()
	}
}

// setTokenFromAppRole get the token from Vault and set it in the client
func (c *Client) setTokenFromAppRole() error {
	var vaultData vaultAuth
	if c.Config.AppRoleCredentials.RoleID == "" {
		return errors.New("No credentials provided")
	}

	url := c.Address
	url.Path = "/v1/auth/approle/login"

	req, err := newRequest("POST", c.Token.ID, url)
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

	c.Lock()
	c.Token.ID = vaultData.ClientToken
	c.Unlock()

	if err = c.setTokenInfo(); err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	if c.Token.Renewable {
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

type tokenRenewable struct {
	Renewable bool `json:"renewable"`
}

func (c *Client) setTokenInfo() error {
	c.Lock()
	defer c.Unlock()
	url := c.Address
	url.Path = "/v1/auth/token/lookup-self"
	var tokenInfo vaultTokenInfo
	req, err := newRequest("GET", c.Token.ID, url)
	if err != nil {
		return err
	}
	res, err := req.execute()
	if err != nil {
		c.Status = err.Error()
		c.isAuthenticated = false
		return err
	}
	if err := json.Unmarshal(res.Data, &tokenInfo); err != nil {
		c.Status = err.Error()
		return err
	}

	c.Token = &tokenInfo
	c.isAuthenticated = true

	return nil
}
