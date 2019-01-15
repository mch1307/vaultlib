//Package vaultlib is a lightweight Go library for reading Vault KV secrets.
//Interacts with Vault server using HTTP API only.
//
//First create a new *config object using NewConfig().
//
//Then create you Vault client using NewClient(*config).
package vaultlib

import (
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// Client holds the vault client
type Client struct {
	sync.RWMutex
	address *url.URL
	//httpClient      *http.Client
	config          *Config
	token           *VaultTokenInfo
	status          string
	isAuthenticated bool
}

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

// AppRoleCredentials holds the app role secret and role ids
type AppRoleCredentials struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// config holds the vault client config
type Config struct {
	Address            string
	MaxRetries         int
	Timeout            time.Duration
	CAPath             string
	InsecureSSL        bool
	AppRoleCredentials *AppRoleCredentials
	Token              string
}

// NewConfig returns a new configuration based on env vars or default value.
// Modify the returned config object to make proper configuration.
func NewConfig() *Config {
	var cfg Config
	appRoleCredentials := new(AppRoleCredentials)
	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.Address = v
	} else {
		cfg.Address = "http://localhost:8200"
	}

	if v := os.Getenv("VAULT_CAPATH"); v != "" {
		cfg.CAPath = v
	}

	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.Token = v
	}

	if v := os.Getenv("VAULT_ROLEID"); v != "" {
		appRoleCredentials.RoleID = v
	}

	if v := os.Getenv("VAULT_SECRETID"); v != "" {
		appRoleCredentials.SecretID = v
	}

	if t := os.Getenv("VAULT_CLIENT_TIMEOUT"); t != "" {
		to, err := strconv.Atoi(t)
		if err != nil {
			cfg.Timeout = time.Duration(30) * time.Second
		}
		clientTimeout := time.Duration(to) * time.Second
		cfg.Timeout = clientTimeout
	} else {
		cfg.Timeout = time.Duration(30 * time.Second)
	}

	if v := os.Getenv("VAULT_SKIP_VERIFY"); v != "" {
		var err error
		cfg.InsecureSSL, err = strconv.ParseBool(v)
		if err != nil {
			cfg.InsecureSSL = true
		}
	} else {
		cfg.InsecureSSL = true
	}
	cfg.AppRoleCredentials = appRoleCredentials
	return &cfg
}

// NewClient returns a new client based on the provided config
func NewClient(c *Config) (*Client, error) {
	// If no config provided, use a new one based on default values and env vars
	if c == nil {
		c = NewConfig()
	}
	var cli Client
	cli.status = "New"
	cli.config = c
	cli.config.Address = c.Address
	cli.config.CAPath = c.CAPath
	cli.config.InsecureSSL = c.InsecureSSL
	cli.config.MaxRetries = c.MaxRetries
	cli.config.Timeout = c.Timeout
	cli.config.Token = c.Token
	cli.config.AppRoleCredentials.RoleID = c.AppRoleCredentials.RoleID
	cli.config.AppRoleCredentials.SecretID = c.AppRoleCredentials.SecretID
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	cli.address = u
	//cli.httpClient = cleanhttp.DefaultPooledClient()
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

	cli.status = "token ready"
	return &cli, nil
}

func (c *Client) getTokenID() string {
	var tk string
	c.withLockContext(func() {
		tk = c.token.ID
	})
	return tk
}

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

func (c *Client) GetStatus() string {
	var status string
	c.withLockContext(func() {
		status = c.status
	})
	return status
}

func (c *Client) IsAuthenticated() bool {
	var authOK bool
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
