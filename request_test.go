package vaultlib

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func Test_newRequest(t *testing.T) {
	testURL, _ := url.Parse("http://localhost:8200")
	emptyReq := new(request)
	type args struct {
		method string
		token  string
		url    *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    *request
		wantErr bool
	}{
		{"err", args{"@", "", testURL}, emptyReq, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRequest(tt.args.method, tt.args.token, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_execute(t *testing.T) {
	nReq, _ := http.NewRequest("GET", "hppT:/localhost.18200", nil)
	htcli := cleanhttp.DefaultPooledClient()
	var rsp vaultResponse
	type fields struct {
		Req        *http.Request
		HTTPClient *http.Client
		Headers    http.Header
		Token      string
	}
	tests := []struct {
		name    string
		fields  fields
		want    vaultResponse
		wantErr bool
	}{
		{"err", fields{nReq, htcli, nil, ""}, rsp, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				Req:        tt.fields.Req,
				HTTPClient: tt.fields.HTTPClient,
				Headers:    tt.fields.Headers,
				Token:      tt.fields.Token,
			}
			got, err := r.execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("request.execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("request.execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_request_setJSONBody(t *testing.T) {
	var cred AppRoleCredentials
	cred.RoleID = "aa"
	cred.SecretID = "bb"
	htCli := new(http.Client)
	url, _ := url.Parse("http://localhot:8200")
	req, _ := newRequest("GET", "", url)
	ch := make(chan int)

	type fields struct {
		Req        *http.Request
		HTTPClient *http.Client
		Headers    http.Header
		Token      string
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"noJSON", fields{req.Req, htCli, nil, ""}, args{ch}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{
				Req:        tt.fields.Req,
				HTTPClient: tt.fields.HTTPClient,
				Headers:    tt.fields.Headers,
				Token:      tt.fields.Token,
			}
			if err := r.setJSONBody(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("request.setJSONBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
