package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

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

	//Get the Vault secret kv_v1/path/my-secret
	resV1, err := vaultCli.GetSecret("kv_v1/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range resV1.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	fmt.Println(vaultCli.GetTokenInfo().NumUses)
	fmt.Printf("Client status: %v\n", vaultCli.GetStatus())

	// Get the Vault secret kv_v2/path/my-secret
	resV2, err := vaultCli.GetSecret("kv_v2/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vaultCli.GetTokenInfo().NumUses)
	for k, v := range resV2.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("# goroutines at the end %v\n", runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	fmt.Printf("# goroutines at the end %v\n", runtime.NumGoroutine())
}
