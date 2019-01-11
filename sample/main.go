package main

import (
	"fmt"
	"log"
	"time"

	vault "github.com/mch1307/vaultlib"
)

func main() {
	// Create a new config. Reads env variables, fallback to default value if needed
	// trace.Start(os.Stderr)
	// defer trace.Stop()
	vcConf := vault.NewConfig()

	// Create new client
	vaultCli, err := vault.NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("AppRole RoleID: %v\n", vaultCli.Config.AppRoleCredentials.RoleID)
	// fmt.Printf("AppRole SecretID: %v\n", vaultCli.Config.AppRoleCredentials.SecretID)
	// fmt.Printf("AppRole token: %v\n", vaultCli.Token)
	// fmt.Printf("Client status: %v\n", vaultCli.Status)
	// Get the Vault secret kv_v1/path/my-secret
	// resV1, err := vaultCli.GetSecret("kv_v1/path/my-secret")
	// _, err = vaultCli.GetSecret("kv_v1/path/my-secret")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for k, v := range resV1.KV {
	// 	fmt.Printf("Secret %v: %v\n", k, v)
	// }
	fmt.Printf("Sleeping: %v\n", vaultCli.Status)
	time.Sleep(10 * time.Second)
	fmt.Printf("Client status: %v\n", vaultCli.Status)

	// Get the Vault secret kv_v2/path/my-secret
	resV2, err := vaultCli.GetSecret("kv_v2/path/json-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range resV2.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
}
