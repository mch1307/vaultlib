package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	vault "github.com/TrueTickets/vaultlib"
)

func main() {
	// Create a new config. Reads env variables, fallback to default value if needed
	//trace.Start(os.Stderr)
	//defer trace.Stop()
	vcConf := vault.NewConfig()

	// Create new client
	fmt.Printf("# goroutines before new cli %v\n", runtime.NumGoroutine())
	vaultCli, err := vault.NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AppRole token: %v\n", vaultCli.GetTokenInfo().ID)
	fmt.Printf("Client status: %v\n", vaultCli.GetStatus())
	//Get the Vault secret kv_v1/path/my-secret
	fmt.Printf("# goroutines before getsecret %v\n", runtime.NumGoroutine())
	resV1, err := vaultCli.GetSecret("kv_v1/path/my-secret")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range resV1.KV {
		fmt.Printf("Secret %v: %v\n", k, v)
	}
	fmt.Printf("# goroutines after getsecret v1 %v", runtime.NumGoroutine())
	fmt.Println(vaultCli.GetTokenInfo().NumUses)
	fmt.Printf("Sleeping: %v\n", vaultCli.GetStatus())
	fmt.Printf("# goroutines before sleep %v\n", runtime.NumGoroutine())
	time.Sleep(90 * time.Second)
	fmt.Printf("# goroutines after sleep %v\n", runtime.NumGoroutine())
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
