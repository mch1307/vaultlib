package vaultlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type request struct {
	Req        *http.Request
	HTTPClient *http.Client
	Headers    http.Header
	Token      string
}

// RawRequest create and execute http request against Vault HTTP API for client.
// Use the client's token for authentication.
//
// Specify http method, Vault path (ie /v1/auth/token/lookup) and optional json payload.
// Return the Vault JSON response .
func (c *Client) RawRequest(method, path string, payload interface{}) (result json.RawMessage, err error) {

	if len(method) == 0 || len(path) == 0 {
		return result, errors.New("Both method and path must be specified")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := c.address.String() + path

	req, err := c.newRequest(method, url)
	if err != nil {
		return result, errors.Wrap(errors.WithStack(err), errInfo())
	}

	if payload != nil {
		if err = req.setJSONBody(payload); err != nil {
			return result, errors.Wrap(errors.WithStack(err), errInfo())
		}
	}

	rsp, err := req.executeRaw()
	if err != nil {
		return rsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	return rsp, nil
}

// Returns a ready to execute request
func (c *Client) newRequest(method, url string) (*request, error) {
	var err error
	req := new(request)

	req.Req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return req, errors.Wrap(errors.WithStack(err), errInfo())
	}

	req.HTTPClient = c.httpClient
	token := c.getTokenID()

	req.Req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Req.Header.Set("X-Vault-token", token)
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

// vaultResponse holds the generic json response from Vault server
type vaultResponse struct {
	RequestID     string          `json:"request_id"`
	LeaseID       string          `json:"lease_id"`
	Renewable     bool            `json:"renewable"`
	LeaseDuration int             `json:"lease_duration"`
	Data          json.RawMessage `json:"data"`
	WrapInfo      json.RawMessage `json:"wrap_info"`
	Warnings      json.RawMessage `json:"warnings"`
	Auth          json.RawMessage `json:"auth"`
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
	defer res.Body.Close()

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
