//Package vaultlib is a lightweight Go library for reading Vault KV secrets.
//Interacts with Vault server using HTTP API only.
//
//First create a new *config object using NewConfig().
//
//Then create you Vault client using NewClient(*config).
//
// See Also
//
// https://github.com/mch1307/vaultlib#vaultlib
package vaultlib

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// Client holds the vault client
type Client struct {
	sync.RWMutex
	address            *url.URL
	httpClient         *http.Client
	appRoleCredentials *AppRoleCredentials
	token              *VaultTokenInfo
	namespace          string
	status             string
	isAuthenticated    bool
}

// VaultTokenInfo holds the Vault token information
type VaultTokenInfo struct {
	Accessor       string      `json:"accessor"`
	CreationTime   int         `json:"creation_time"`
	CreationTTL    int         `json:"creation_ttl"`
	DisplayName    string      `json:"display_name"`
	EntityID       string      `json:"entity_id"`
	ExpireTime     interface{} `json:"expire_time"`
	ExplicitMaxTTL int         `json:"explicit_max_ttl"`
	ID             string      `json:"id"`
	IssueTime      time.Time   `json:"issue_time"`
	Meta           interface{} `json:"meta"`
	NumUses        int         `json:"num_uses"`
	Orphan         bool        `json:"orphan"`
	Path           string      `json:"path"`
	Policies       []string    `json:"policies"`
	Renewable      bool        `json:"renewable"`
	TTL            int         `json:"ttl"`
	Type           string      `json:"type"`
}

// NewClient returns a new client based on the provided config
func NewClient(c *Config) (*Client, error) {
	var caPool *x509.CertPool
	// If no config provided, use a new one based on default values and env vars
	if c == nil {
		c = NewConfig()
	}

	var cli Client
	cli.status = "New"
	cli.appRoleCredentials = new(AppRoleCredentials)
	cli.appRoleCredentials.RoleID = c.AppRoleCredentials.RoleID
	cli.appRoleCredentials.SecretID = c.AppRoleCredentials.SecretID

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	cli.address = u

	if c.CACert != "" {
		caCert, err := ioutil.ReadFile(c.CACert)
		if err != nil {
			c.CACert = ""
		}
		caPool = x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
	}

	cli.httpClient = &http.Client{}

	tsp := cleanhttp.DefaultPooledTransport()

	tsp.TLSClientConfig = &tls.Config{
		RootCAs:            caPool,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: c.InsecureSSL,
	}

	cli.httpClient.Transport = tsp

	cli.token = new(VaultTokenInfo)
	cli.token.ID = c.Token

	if cli.token.ID == "" {
		err = cli.setTokenFromAppRole()
		if err != nil {
			cli.status = "Authentication Error: " + err.Error()
			return &cli, err
		}
	} else {
		if err = cli.setTokenInfo(); err != nil {
			cli.status = "Authentication Error: " + err.Error()
			return &cli, err
		}
		if cli.token.Renewable {
			go cli.renewToken()
		}

	}
	cli.status = "Token ready"
	return &cli, nil
}

func (c *Client) getTokenID() string {
	var tk string
	c.withLockContext(func() {
		tk = c.token.ID
	})
	return tk
}

// GetTokenInfo returns the current token information
func (c *Client) GetTokenInfo() *VaultTokenInfo {
	vt := new(VaultTokenInfo)
	c.withLockContext(func() {
		vt = c.token
	})
	return vt

}

func (c *Client) setStatus(status string) {
	c.withLockContext(func() {
		c.status = status
	})
}

// GetStatus return the last action status/log
func (c *Client) GetStatus() string {
	var status string
	c.withLockContext(func() {
		status = c.status
	})
	return status
}

// IsAuthenticated returns bool if last call to vault was ok
func (c *Client) IsAuthenticated() bool {
	authOK := false
	c.withLockContext(func() {
		authOK = c.isAuthenticated
	})
	return authOK
}

func (c *Client) withLockContext(fn func()) {
	c.Lock()
	defer c.Unlock()
	fn()
}
