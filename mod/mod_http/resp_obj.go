package mod_http

import (
	"net/http"
	"fmt"
)

type JSResponse struct {
	w http.ResponseWriter
}

func (resp *JSResponse) WriteHead(statusCode int, headers map[string]interface{}) {
	if headers != nil {
		header := resp.w.Header()
		for k, v := range headers {
			header.Add(k, fmt.Sprintf("%v",v))
		}
	}
	resp.w.WriteHeader(statusCode)
}

func (resp *JSResponse) SetHeader(key string, value interface{}) {
	resp.w.Header().Set(key, fmt.Sprintf("%v", value))
}

func (resp *JSResponse) Write(chunk []byte, encoding string) {
	resp.w.Write(chunk)
}

func (resp *JSResponse) End(data []byte, encoding string) {
	resp.w.Write(data)
}

