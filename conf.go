package vaultlib

import (
	"os"
	"strconv"
	"time"
)

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

// NewConfig returns a new configuration based on env vars or default value.
//
// Reads ENV:
//	VAULT_ADDR            Vault server URL (default http://localhost:8200)
//	VAULT_CAPATH          CA path
//	VAULT_TOKEN           Vault Token
//	VAULT_ROLEID          Vault app role id
//	VAULT_SECRETID        Vault app role secret id
//	VAULT_CLIENT_TIMEOUT  Client timeout
//	VAULT_SKIP_VERIFY     Do not check SSL
//
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
