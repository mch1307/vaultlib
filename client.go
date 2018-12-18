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

// Config holds the vault client config
type Config struct {
	Address            string
	MaxRetries         int
	Timeout            time.Duration
	CAPath             string
	InsecureSSL        bool
	AppRoleCredentials AppRoleCredentials
	Token              string
}

// AppRoleCredentials holds the app role secret and role ids
type AppRoleCredentials struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// VaultConfig returns a configuration based on ENV vars and default value
// Modify the returned Config
func VaultConfig() *Config {
	var cfg *Config
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
		cfg.AppRoleCredentials.RoleID = v
	}

	if v := os.Getenv("VAULT_SECRETID"); v != "" {
		cfg.AppRoleCredentials.SecretID = v
	}

	if t := os.Getenv("VAULT_CLIENT_TIMEOUT"); t != "" {
		to, err := strconv.Atoi(t)
		if err != nil {
			//fmt.Errorf("Could not parse VAULT_CLIENT_TIMEOUT")
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
			//fmt.Errorf("Error parsing VAULT_SKIP_VERIFY")
			cfg.InsecureSSL = true
		}
	}
	return cfg
}

// SetAppRole sets the app role role_id and secret_id in config
func (c *Config) SetAppRole(cred AppRoleCredentials) error {
	c.AppRoleCredentials = cred
	return nil
}

// NewClient returns a new client based on the provided config
func NewClient(c *Config) (*VaultClient, error) {
	// If no config provided, use a new one based on default values and env vars
	if c == nil {
		c = VaultConfig()
	}
	var cli VaultClient
	cli.Config = c
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	cli.Address = u
	cli.HTTPClient = cleanhttp.DefaultPooledClient()
	cli.Token = c.Token

	return &cli, nil
}
