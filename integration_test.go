package vaultlib

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

var vaultRoleID, vaultSecretID string

func TestMain(m *testing.M) {
	fmt.Println("TestMain: Preparing Vault server")
	prepareVault()
	ret := m.Run()
	os.Exit(ret)
}

func prepareVault() {
	go startVault()
	// wait 20 seconds vault
	time.Sleep(20 * time.Second)
	cmd := exec.Command("./vault", "read", "-field=role_id", "auth/approle/role/my-role/role-id")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_TOKEN=my-dev-root-vault-token")
	cmd.Env = append(cmd.Env, "VAULT_ADDR=http://localhost:8200")

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("error getting role id %v %v", err, out)
	}
	vaultRoleID = string(out)

	cmd = exec.Command("./vault", "write", "-field=secret_id", "-f", "auth/approle/role/my-role/secret-id")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_TOKEN=my-dev-root-vault-token")
	cmd.Env = append(cmd.Env, "VAULT_ADDR=http://localhost:8200")
	out, err = cmd.Output()
	if err != nil {
		log.Fatalf("error getting secret id %v", err)
	}
	vaultSecretID = string(out)
	os.Unsetenv("VAULT_TOKEN")

}

func startVault() {
	cmd := exec.Command("./initVaultDev.sh")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	_ = cmd.Wait()

}
