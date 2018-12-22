# vaultlib (work in progress)

Lightweight, simple Go library for Vault secret reading (http API).

Connect to Vault through app role or token.

Reads kv secret values

## Features

* Connect to Vault through app role
* Read Vault secret, `kv` type (v1 or v2 "versioned")

## Config

Reads the following env variables:

`VAULT_ADDR` : Vault server URL (including scheme and port)

`VAULT_CAPATH` : CA path

`VAULT_TOKEN` : Vault Token

`VAULT_ROLEID` : Vault app role id

`VAULT_SECRETID` : Vault app role secret id

`VAULT_CLIENT_TIMEOUT` : Client timeout

`VAULT_SKIP_VERIFY` : Do not check SSL


If not set, fallback to default values.

Can also be set programmatically through the `Config` object

## Getting Started

> For a simple, working example, check the sample folder.

```golang
package main

import (
	"fmt"
	"log"
	"os"

	vault "github.com/mch1307/vaultlib"
)

func main() {
	// Config can be set through ENV before invoking NewConfig
	os.Setenv("VAULT_ADDR", "http://localhost:8200")

	// Create a new config. Reads env variables, fallback to default value if needed
	vcConf := vault.NewConfig()

	// Config can also be done programmtically
	vcConf.Address = "http://localhost:8200"

	// set app role credentials (ie after reading from docker secret)
	// vcConf.AppRoleCredentials.RoleID = "myRoleID"
	// vcConf.AppRoleCredentials.SecretID = "mySecretID"

	// Create new client
	vaultCli, err := vault.NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}

	// Get the Vault secret data
	kv := make(map[string]string)
	kv, err = vaultCli.GetVaultSecret("my_kv/my_org/my_secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range kv {
		fmt.Printf("secret %v: %v\n", k, v)
	}
}
```