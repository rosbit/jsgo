package mod_http

/**
 * http module implementation, which is invoked via `var http = require('http')` in js.
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 6, 2018
 */

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"fmt"
)

var (
	jsEnv *js.JSEnv
)

type HttpModule struct {}

func NewHttpModule(ctx *js.JSEnv) interface{} {
	jsEnv = ctx
	return &HttpModule{}
}

func (m *HttpModule) CreateServer(jsCallback *js.EcmaObject, isFastCGI bool) (*HttpServer, error) {
	if jsCallback == nil {
		return nil, fmt.Errorf("argument must be function(request, respsone)")
	}
	server := CreateHttpServer(jsCallback, isFastCGI)
	return server, nil
}

