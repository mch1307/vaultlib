package main

import (
	"fmt"
	"log"

	vault "github.com/mch1307/vaultlib"
)

func main() {
	// Create a new config. Reads env variables, fallback to default value if needed
	vcConf := vault.NewConfig()

	// Create new client
	vaultCli, err := vault.NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}

	// Get the Vault secret kv_v1/path/my-secret
	kv, err := vaultCli.GetVaultSecret("kv_v1/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range kv {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	// Get the Vault secret kv_v2/path/my-secret
	kv2, err := vaultCli.GetVaultSecret("kv_v2/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range kv2 {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
}
