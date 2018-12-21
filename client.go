package vaultlib

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// VaultClient holds the vault client
type VaultClient struct {
	Address    *url.URL
	HTTPClient *http.Client
	Config     *Config
	Token      string
}

// AppRoleCredentials holds the app role secret and role ids
type AppRoleCredentials struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// Config holds the vault client config
type Config struct {
	Address            string
	MaxRetries         int
	Timeout            time.Duration
	CAPath             string
	InsecureSSL        bool
	AppRoleCredentials *AppRoleCredentials
	Token              string
}

// NewConfig returns a configuration based on ENV vars and default value
// Modify the returned Config
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

// SetAppRole sets the app role role_id and secret_id in config
func (c *Config) SetAppRole(cred AppRoleCredentials) error {
	c.AppRoleCredentials = &cred
	return nil
}

// NewClient returns a new client based on the provided config
func NewClient(c *Config) (*VaultClient, error) {
	// If no config provided, use a new one based on default values and env vars
	if c == nil {
		c = NewConfig()

	}
	var cli VaultClient
	cli.Config = c
	cli.Config.Address = c.Address
	cli.Config.CAPath = c.CAPath
	cli.Config.InsecureSSL = c.InsecureSSL
	cli.Config.MaxRetries = c.MaxRetries
	cli.Config.Timeout = c.Timeout
	cli.Config.Token = c.Token
	cli.Config.AppRoleCredentials.RoleID = c.AppRoleCredentials.RoleID
	cli.Config.AppRoleCredentials.SecretID = c.AppRoleCredentials.SecretID
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	cli.Address = u
	cli.HTTPClient = cleanhttp.DefaultPooledClient()
	cli.Token = c.Token

	return &cli, nil
}
