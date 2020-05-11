package vaultlib

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

var vaultRoleID, vaultSecretID, noKVRoleID, noKVSecretID string

var vaultVersion string

func init() {
	flag.StringVar(&vaultVersion, "vaultVersion", "1.4.1", "provide vault version to be tested against")
	flag.Parse()
}
func TestMain(m *testing.M) {

	fmt.Println("Testing with Vault version", vaultVersion)
	fmt.Println("TestMain: Preparing Vault server")
	prepareVault()
	ret := m.Run()
	os.Exit(ret)
}

func prepareVault() {
	err := startVault(vaultVersion)
	if err != nil {
		log.Fatalf("Error in initVaultDev.sh %v", err)
	}
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

	cmd = exec.Command("./vault", "read", "-field=role_id", "auth/approle/role/no-kv/role-id")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_TOKEN=my-dev-root-vault-token")
	cmd.Env = append(cmd.Env, "VAULT_ADDR=http://localhost:8200")

	out, err = cmd.Output()
	if err != nil {
		log.Fatalf("error getting role id %v %v", err, out)
	}
	noKVRoleID = string(out)

	cmd = exec.Command("./vault", "write", "-field=secret_id", "-f", "auth/approle/role/no-kv/secret-id")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_TOKEN=my-dev-root-vault-token")
	cmd.Env = append(cmd.Env, "VAULT_ADDR=http://localhost:8200")
	out, err = cmd.Output()
	if err != nil {
		log.Fatalf("error getting secret id %v", err)
	}
	noKVSecretID = string(out)
	os.Unsetenv("VAULT_TOKEN")

}

func startVault(version string) error {
	cmd := exec.Command("bash", "./test-files/initVaultDev.sh", version)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil

}
