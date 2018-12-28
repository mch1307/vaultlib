/*
Package vaultlib is a lightweight Go library for reading Vault KV secrets.
Interacts with Vault server using HTTP API only.
First create a new *Config object using NewConfig()
Then create you Vault client using NewClient(*Config)
*/
package vaultlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/pkg/errors"
)

type request struct {
	Req        *http.Request
	HTTPClient *http.Client
	Headers    http.Header
	Token      string
}

// Returns a ready to execute request
func newRequest(method, token string, url *url.URL, htCli *http.Client) (*request, error) {
	var err error
	req := new(request)
	req.HTTPClient = cleanhttp.DefaultPooledClient()
	req.Token = token

	req.Req, err = http.NewRequest(method, url.String(), nil)
	if err != nil {
		return req, err
	}
	req.Req.Header.Set("Content-Type", "application/json")
	if req.Token != "" {
		req.Req.Header.Set("X-Vault-Token", req.Token)
	}
	return req, err

}

// Adds JSON formatted body to request
func (r *request) setJSONBody(val interface{}) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	r.Req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return nil

}

// Executes the request
func (r *request) execute() (vaultResponse, error) {
	var vaultRsp vaultResponse
	res, err := r.HTTPClient.Do(r.Req)
	if err != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	if res.StatusCode != http.StatusOK {
		httpErr := fmt.Sprintf("Vault http call %v returned %v. Body: %v", r.Req.URL.String(), res.Status, string(body))
		return vaultRsp, errors.New(httpErr)
	}

	jsonErr := json.Unmarshal(body, &vaultRsp)
	if jsonErr != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	return vaultRsp, nil

}
