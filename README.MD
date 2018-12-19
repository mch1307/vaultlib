# vaultlib (work in progress)

Light Go library for simple Vault (http) reading

## Features

 [X] Connect to Vault through app role

 [] Read secret(s)

## Example

```golang
package main

import (
	"log"
	"os"

	vault "github.com/mch1307/vaultlib"
)

func main() {
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	// Create a new config. Reads env variables, fallback to default value if needed
	conf := vault.NewConfig()

	conf.CAPath = "/etc/myapp/ssl"
	// Create new client
	cli, err := vault.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	// set app role credentials (ie from docker secret)
	conf.AppRoleCredentials.RoleID = "acalcasdavfdvsdfv3d"
	conf.AppRoleCredentials.SecretID = "sdavs3vafdv4afdsv"
	
	// Authenticate through app role
	err = cli.SetTokenFromAppRole()
	if err != nil {
		log.Fatal(err)
	}
	
}
`` 