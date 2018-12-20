package vaultlib

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

var (
	vaultMount  VaultSecretMount
	vaultMounts []VaultSecretMount
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain: starting Vault stub")
	InitVaultStub()
	ret := m.Run()
	os.Exit(ret)
}

func InitVaultStub() {
	//http.HandleFunc("/v1/auth/approle/login", loginHandler)
	http.HandleFunc("/v1/sys/internal/ui/mounts", mountsHandler)
	http.HandleFunc("/v1/secret/data/test-test", vaultSecretV2Handler)
	http.HandleFunc("/v1/kv1/test-test", vaultSecretV1Handler)

	go http.ListenAndServe(":8500", nil)
}

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	var vaultAuthResponse VaultAuthResponse
// 	var vaultCredentials VaultAppRoleCredentials
// 	err := json.NewDecoder(r.Body).Decode(&vaultCredentials)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 	} else {
// 		if vaultCredentials.RoleID == "ok" {
// 			vaultAuthResponse.Auth.ClientToken = "ok"
// 		} else if vaultCredentials.RoleID == "ko" {
// 			vaultAuthResponse.Auth.ClientToken = "ko"
// 		} else if vaultCredentials.RoleID == "5xx" {
// 			vaultAuthResponse.Auth.ClientToken = "5xx"
// 			w.WriteHeader(http.StatusInternalServerError)
// 		} else if vaultCredentials.RoleID == "jsonKO" {
// 			var dummy = `{"invalid":"json"
// 			`
// 			w.Write([]byte(dummy))
// 		}
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(vaultAuthResponse)
// }

func mountsHandler(w http.ResponseWriter, r *http.Request) {
	var okResp = `{
		"request_id": "2db2f477-9f5e-6d95-c749-fa13f6faa98e",
		"lease_id": "",
		"renewable": false,
		"lease_duration": 0,
		"data": {
			"secret": {
				"cubbyhole/": {
					"accessor": "cubbyhole_c7e85579",
					"config": {
						"default_lease_ttl": 0,
						"force_no_cache": false,
						"max_lease_ttl": 0,
						"plugin_name": ""
					},
					"description": "per-token private secret storage",
					"local": true,
					"options": null,
					"seal_wrap": false,
					"type": "cubbyhole"
				},
				"identity/": {
					"accessor": "identity_47ff67dd",
					"config": {
						"default_lease_ttl": 0,
						"force_no_cache": false,
						"max_lease_ttl": 0,
						"plugin_name": ""
					},
					"description": "identity store",
					"local": false,
					"options": null,
					"seal_wrap": false,
					"type": "identity"
				},
				"kvV1/": {
					"accessor": "kv_8bf51dd1",
					"config": {
						"default_lease_ttl": 0,
						"force_no_cache": false,
						"max_lease_ttl": 0,
						"plugin_name": ""
					},
					"description": "",
					"local": false,
					"options": {
						"version": "1"
					},
					"seal_wrap": false,
					"type": "kv"
				},
				"kvV2/": {
					"accessor": "kv_0cb1bd73",
					"config": {
						"default_lease_ttl": 0,
						"force_no_cache": false,
						"max_lease_ttl": 0,
						"plugin_name": ""
					},
					"description": "key/value secret storage",
					"local": false,
					"options": {
						"version": "2"
					},
					"seal_wrap": false,
					"type": "kv"
				},
				"sys/": {
					"accessor": "system_f177f6e8",
					"config": {
						"default_lease_ttl": 0,
						"force_no_cache": false,
						"max_lease_ttl": 0,
						"plugin_name": ""
					},
					"description": "system endpoints used for control, policy and debugging",
					"local": false,
					"options": null,
					"seal_wrap": false,
					"type": "system"
				}
			}
		},
		"wrap_info": null,
		"warnings": null,
		"auth": null
	}`
	var koResp = `{
		"request_id": "660de5d6-92ff-a614-0e31-0460ee0c49c2",
		"lease_id": "",
		"renewable": false,
		"lease_duration": 0,
		"data": {
			"auth": {},
			"secret": {}
		},
		"wrap_info": null,
		"warnings": null,
		"auth": null
	}`
	if r.Header.Get("X-Vault-Token") == "goodToken" {
		w.Write([]byte(okResp))
	} else {
		w.Write([]byte(koResp))
	}

}

func vaultSecretV2Handler(w http.ResponseWriter, r *http.Request) {
	var v2Secret = `{
    "request_id": "af823ad6-da9f-29ab-b0dd-13222f0d4033",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "data": {
            "test-key": "test-value"
        },
        "metadata": {
            "created_time": "2018-11-07T08:14:23.912677243Z",
            "deletion_time": "",
            "destroyed": false,
            "version": 1
        }
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`
	var koResp = `{
    "errors": [
        "permission denied"
    ]
}`
	if r.Header.Get("X-Vault-Token") == "goodToken" {
		w.Write([]byte(v2Secret))
	} else {
		w.Write([]byte(koResp))
	}
}

func vaultSecretV1Handler(w http.ResponseWriter, r *http.Request) {
	var v1Secret = `{
    "request_id": "d915438f-8aeb-0241-9927-d0b7b0c42166",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 2764800,
    "data": {
        "test-key": "test-value"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}`
	var koResp = `{
    "errors": [
        "permission denied"
    ]
}`
	if r.Header.Get("X-Vault-Token") == "goodToken" {
		w.Write([]byte(v1Secret))
	} else {
		w.Write([]byte(koResp))
	}
}
