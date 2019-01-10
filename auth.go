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
func (c *Client) setTokenFromAppRole() error {
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

type tokenRenewable struct {
	Renewable bool `json:"renewable"`
}

func (c *Client) isCurrentTokenRenewable() bool {
	var tokenRenew tokenRenewable
	url := c.Address
	url.Path = "/v1/auth/token/lookup-self"
	req, err := newRequest("GET", c.Token, url)
	if err != nil {
		return false
	}
	res, err := req.execute()
	if err != nil {
		c.Status = err.Error()
		return false
	}
	if err := json.Unmarshal(res.Data, &tokenRenew); err != nil {
		c.Status = err.Error()
		return false
	}
	return tokenRenew.Renewable
}
