# vaultlib

[![Build Status](https://travis-ci.org/mch1307/vaultlib.svg?branch=master)](https://travis-ci.org/mch1307/vaultlib)
[![Coverage Status](https://coveralls.io/repos/github/mch1307/vaultlib/badge.svg?branch=master)](https://coveralls.io/github/mch1307/vaultlib?branch=master) [![GoDoc](https://godoc.org/github.com/mch1307/vaultlib?status.svg)](https://godoc.org/github.com/mch1307/vaultlib) [![Go Report Card](https://goreportcard.com/badge/github.com/mch1307/vaultlib)](https://goreportcard.com/report/github.com/mch1307/vaultlib)

Lightweight, simple Go library for Vault secret reading (http API).

Connect to Vault through app role or token.

Reads kv secret values

## Features

* Connect to Vault through app role
* Read Vault secret, `kv` type (v1 or v2 "versioned")
* Automatically renew token
* Execute any HTTP request on Vault (RawRequest)

## Config

Configuration can be done through env variables or programmatically through the `Config` object
The following env variables are supported:

```bash
VAULT_ADDR            Vault server URL (default http://localhost:8200)
VAULT_CAPATH          CA path
VAULT_TOKEN           Vault Token
VAULT_ROLEID          Vault app role id
VAULT_SECRETID        Vault app role secret id
VAULT_CLIENT_TIMEOUT  Client timeout
VAULT_SKIP_VERIFY     Do not check SSL
```

If not set, `vaultlib` will fallback to safe default values.

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
    kv, err := vaultCli.GetSecret("my_kv/my_org/my_secret")
    if err != nil {
        fmt.Println(err)
    }
    for k, v := range kv {
        fmt.Printf("secret %v: %v\n", k, v)
    }
}
```