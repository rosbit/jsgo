package mod_http

/**
 * http server implementation, which will be invoked via `var server = http.createServer(function (req, resp){})`
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 6, 2018
 */

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"fmt"
	"net"
	"net/http"
	"strings"
	"encoding/json"
	sc "github.com/rosbit/gojs/server_counter"
)

type httpdHandlerParams struct {
	w http.ResponseWriter
	r *http.Request
	done chan bool
}

type HttpServer struct {
	jsCallback *js.EcmaObject // function callback(request, response)
	jsMod  *js.EcmaObject

	listener net.Listener
	serverStarted bool

	httpdHandlers chan httpdHandlerParams
}

func CreateHttpServer(jsCallback *js.EcmaObject) *HttpServer {
	httpdHandlers := make(chan httpdHandlerParams, 5) // if more than 5 requests at same time, they will blocked.
	return &HttpServer{jsCallback:jsCallback, httpdHandlers:httpdHandlers}
}

func (s *HttpServer) saveModule(jsMod *js.EcmaObject) {
	s.jsMod = jsMod
}

func (s *HttpServer) Listen(port int, hostname string) error {
	if s.serverStarted {
		return fmt.Errorf("I am running already, don't call listen() many times")
	}

	var e error
	if strings.HasPrefix(hostname, "unix:") {
		fn := hostname[5:]
		s.listener, e = net.Listen("unix", fn)
	} else {
		server := fmt.Sprintf("%s:%d", hostname, port)
		s.listener, e = net.Listen("tcp", server)
	}
	if e != nil {
		return fmt.Errorf("%v", e)
	}
	s.serverStarted = true
	sc.IncServerCount()

	go s.accept()
	return nil
}

func (s *HttpServer) Close() {
	if s.serverStarted {
		s.serverStarted = false
		s.listener.Close()
		sc.DecServerCount()
	}
}

func writeJson(w http.ResponseWriter, v interface{}) {
	b, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeError(w http.ResponseWriter, err string, code int) {
	v := map[string]interface{} {"code": code, "message": err}
	w.WriteHeader(code)
	writeJson(w, v)
}

func writeResult(w http.ResponseWriter, res interface{}) {
	switch res.(type) {
	case string:
		w.Write([]byte(res.(string)))
	case []byte:
		w.Write(res.([]byte))
	case []interface{}:
		writeJson(w, res)
	case map[string]interface{}:
		writeJson(w, res)
	default:
		b, _ := json.Marshal(res)
		w.Write(b)
	}
}

func (s *HttpServer) handleHttp(w http.ResponseWriter, r *http.Request) {
	req := newJSRequest(r)
	if req.Error != "" {
		writeError(w, req.Error, http.StatusInternalServerError)
		return
	}

	reqMod, err := jsEnv.CreateEcmascriptModule(req)
	if err != nil {
		writeError(w, "Failed to create request module", http.StatusInternalServerError)
		return
	}
	defer jsEnv.DestroyEcmascriptModule(reqMod)

	resp := &JSResponse{w}
	respMod, err := jsEnv.CreateEcmascriptModule(resp)
	if err != nil {
		writeError(w, "Failed to create response module", http.StatusInternalServerError)
		return
	}
	defer jsEnv.DestroyEcmascriptModule(respMod)

	res, err := jsEnv.CallEcmascriptFunc(s.jsCallback, reqMod, respMod)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res != nil {
		writeResult(w, res)
	}
}

func (s *HttpServer) handleHttpProxy(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	httpdParam := httpdHandlerParams{w, r, done}
	s.httpdHandlers <- httpdParam
	<-done
	close(done)
}

func (s *HttpServer) accept() {
	sm := http.NewServeMux()
	sm.HandleFunc("/", s.handleHttpProxy)

	go func() {
		for {
			httpdParam := <-s.httpdHandlers
			w, r, done := httpdParam.w, httpdParam.r, httpdParam.done
			if r == nil {
				break
			}
			s.handleHttp(w, r)
			done <- true
		}
	}()

	http.Serve(s.listener, sm)
	s.httpdHandlers <- httpdHandlerParams{nil, nil, nil} // let the go-routine done

	jsEnv.DestroyEcmascriptModule(s.jsMod)
	jsEnv.DestroyEcmascriptFunc(s.jsCallback)

	s.Close()
}
