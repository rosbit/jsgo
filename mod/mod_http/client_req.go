package mod_http

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

type reqResult struct {
	StatusCode int
	Headers map[string][]string
	Data []byte
}

func (m *HttpModule) Request(options map[string]interface{}, jsCallback *js.EcmaObject) (map[string]interface{}, error) {
	defer func() {
		if jsCallback != nil {
			jsEnv.DestroyEcmascriptFunc(jsCallback)
		}
	}()

	if options == nil || len(options) == 0 {
		return nil, fmt.Errorf("no options specified to request")
	}

	var method string
	var url string

	if m, ok := options["method"]; ok {
		method = m.(string)
	}
	if method == "" {
		method = "GET"
	} else {
		method = strings.ToUpper(method)
	}
	if u, ok := options["url"]; ok {
		url = u.(string)
	} else {
		return nil, fmt.Errorf("please specify url")
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, nil)
	if err != nil {
		return nil, err
	}

	if reqCookies, ok := options["cookies"]; ok {
		for key, value := range reqCookies.(map[string]interface{}) {
			req.AddCookie(&http.Cookie{Name: key, Value: fmt.Sprintf("%v", value)})
		}
	}

	if reqQuery, ok := options["query"]; ok {
		req.URL.RawQuery = reqQuery.(string)
	}

	if method != "GET" && method != "HEAD" {
		if body, ok := options["body"]; ok {
			// "form" is deprecated.
			if body, ok = options["form"]; ok {
				// Only set the Content-Type to application/x-www-form-urlencoded
				// when someone uses "form", not for "body".
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			b := body.(string)
			req.ContentLength = int64(len(b))
			req.Body = ioutil.NopCloser(strings.NewReader(b))
		}
	}

	// Set these last. That way the code above doesn't overwrite them.
	if reqHeaders, ok := options["headers"]; ok {
		for key, value := range reqHeaders.(map[string]interface{}) {
			req.Header.Set(key, fmt.Sprintf("%v", value))
		}
	}

	client := http.Client{}
	if res, err := client.Do(req); err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		if body, err := ioutil.ReadAll(res.Body); err != nil {
			return nil, err
		} else {
			resp := map[string]interface{}{"statusCode":res.StatusCode, "headers": res.Header, "data": string(body)}
			if jsCallback != nil {
				jsEnv.CallEcmascriptFunc(jsCallback, resp)
				return nil, nil
			}
			return resp, nil
		}
	}
}
