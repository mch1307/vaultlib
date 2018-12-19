package vaultlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"bytes"

	"github.com/pkg/errors"
)

type request struct {
	Method  string
	URL     *url.URL
	Path    string
	Req     *http.Request
	Headers http.Header
	Token   string
	Body    []byte
}

func (r *request) setJSONBody(val interface{}) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	r.Req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	return nil

}

func (r *request) prepareRequest() error {
	var err error
	r.Req, err = http.NewRequest(r.Method, r.URL.String(), nil)
	if err != nil {
		return err
	}
	r.Req.Header.Set("Content-Type", "application/json")
	if r.Token != "" {
		r.Req.Header.Set("X-Vault-Token", r.Token)
	}
	return nil
}

// execute executes the request using the provided client
func (r *request) execute(c *http.Client) (VaultResponse, error) {
	var vaultRsp VaultResponse
	res, err := c.Do(r.Req)
	if err != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	if res.StatusCode != http.StatusOK {
		httpErr := fmt.Sprintf("Vault http call %v returned %v. Body: %v", r.URL.String(), res.Status, string(body))
		return vaultRsp, errors.New(httpErr)
	}

	jsonErr := json.Unmarshal(body, &vaultRsp)
	if jsonErr != nil {
		return vaultRsp, errors.Wrap(errors.WithStack(err), errInfo())
	}

	return vaultRsp, nil

}
type VaultResponse struct {
	RequestID     string          `json:"request_id"`
	LeaseID       string          `json:"lease_id"`
	Renewable     bool            `json:"renewable"`
	LeaseDuration int             `json:"lease_duration"`
	Data          json.RawMessage `json:"data"`
	WrapInfo      json.RawMessage `json:"wrap_info"`
	Warnings      json.RawMessage `json:"warnings"`
	Auth          json.RawMessage `json:"auth"`
}

func errInfo() (info string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function + ":" + strconv.Itoa(frame.Line)
}
