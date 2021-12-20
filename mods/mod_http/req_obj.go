package mod_http

import (
	"net/http"
	"io/ioutil"
)

type JSRequest struct {
	Method string
	Uri string
	Path string
	HttpVersion string
	Headers map[string][]string
	Host string
	RemoteAddr string
	Auth string
	RawQuery string
	Hash string

	r *http.Request
	log *http_access_log
}

func newJSRequest(r *http.Request, log *http_access_log) *JSRequest {
	req := &JSRequest{r:r, log:log}
	req.Method = r.Method
	req.Uri = r.RequestURI
	req.Path = r.URL.Path
	req.HttpVersion = r.Proto
	req.Headers = r.Header
	req.Host = r.Host
	req.RemoteAddr = r.RemoteAddr
	if r.URL.User != nil {
		req.Auth = r.URL.User.String()
	}
	req.RawQuery = r.URL.RawQuery
	req.Hash = r.URL.Fragment

	return req
}

func (req *JSRequest) ParseParams() (map[string][]string, error) {
	if err := req.r.ParseForm(); err != nil {
		return nil, err
	} else {
		return map[string][]string(req.r.Form), nil
	}
}

func (req *JSRequest) Param(key string) string {
	return req.r.FormValue(key)
}

func (req *JSRequest) ReadBody() ([]byte, error) {
	if req.r.Method == "POST" || req.r.Method == "PUT" {
		if b, err := ioutil.ReadAll(req.r.Body); err != nil {
			return nil, err
		} else {
			return b, nil
		}
	}
	return nil, nil
}

func (req *JSRequest) GetLocalTime() string {
	return req.log.formatTime("")
}
