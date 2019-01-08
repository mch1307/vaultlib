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
func newRequest(method, token string, url *url.URL) (*request, error) {
	var err error
	req := new(request)

	req.Req, err = http.NewRequest(method, url.String(), nil)
	if err != nil {
		return req, errors.Wrap(errors.WithStack(err), errInfo())
	}

	req.HTTPClient = cleanhttp.DefaultPooledClient()
	req.Token = token

	req.Req.Header.Set("Content-Type", "application/json")
	if req.Token != "" {
		req.Req.Header.Set("X-Vault-Token", req.Token)
	}
	return req, errors.Wrap(errors.WithStack(err), errInfo())

}

// Adds JSON formatted body to request
func (r *request) setJSONBody(val interface{}) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(errors.WithStack(err), errInfo())
	}
	r.Req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return nil

}

// Executes the request, parse the result to vaultResponse
func (r *request) execute() (vaultResponse, error) {
	var vaultRsp vaultResponse
	res, err := r.executeRaw()
	if err != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	jsonErr := json.Unmarshal(res, &vaultRsp)
	if jsonErr != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	return vaultRsp, nil

}

// Executes the raw request, does not parse Vault response
func (r *request) executeRaw() ([]byte, error) {
	res, err := r.HTTPClient.Do(r.Req)
	if err != nil {
		return nil, errors.Wrap(errors.WithStack(err), errInfo())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return body, errors.Wrap(errors.WithStack(err), errInfo())
	}

	if res.StatusCode != http.StatusOK {
		httpErr := fmt.Sprintf("Vault http call %v returned %v. Body: %v", r.Req.URL.String(), res.Status, string(body))
		return body, errors.New(httpErr)
	}

	return body, nil

}
