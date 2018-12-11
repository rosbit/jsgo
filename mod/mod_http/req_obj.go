package mod_http

import (
	"net/http"
	"io/ioutil"
)

type JSRequest struct {
	Method string
	Uri string
	HttpVersion string
	Headers map[string][]string
	Params map[string][]string
	Host string
	RemoteAddr string
	Body []byte
	Error string
}

func newJSRequest(r *http.Request) *JSRequest {
	req := &JSRequest{}
	req.Method = r.Method
	req.Uri = r.RequestURI
	req.HttpVersion = r.Proto
	req.Headers = r.Header
	if err := r.ParseForm(); err != nil {
		req.Error = err.Error()
	} else {
		req.Params = r.Form
	}
	req.Host = r.Host
	req.RemoteAddr = r.RemoteAddr
	if r.Method == "POST" || r.Method == "PUT" {
		if b, err := ioutil.ReadAll(r.Body); err != nil {
			req.Error = err.Error()
		} else {
			req.Body = b
		}
	}
	return req
}
