# Sample

## Setup vault server

> You will need vault installed. The `initVaultDev.sh` will first do a `killall vault` so make sure you don't have another vault process running.

Open a terminal and launch `initVaultDev.sh`

```shell
$ source initVaultDev.sh
```

It will:

* start a new Vault server in dev mode (root token = )
* create and enable two secret engines: `kv_v1/path`  and `kv_v2/path`
* write demo secrets to each engine
* enable app role auth
* create a policy
* create a role_id and secret_id linked to rhe policy
* export the role_id and secret_id

## Run go code

From the same session, you can run the sample `main.go` file

```shell
$ go run main.go
```

If everything goes well, you should have this output:

```shell
Secret my-v1-secret: my-v1-secret-value
Secret my-second-secret: my-second-secret-value
Secret my-first-secret: my-first-secret-value
```