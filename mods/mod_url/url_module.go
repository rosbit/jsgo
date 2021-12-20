package mod_url

/**
 * url module implementation, which is invoked via `var url = require('url')` in js.
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 11, 2018
 */

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"fmt"
	"net/url"
)

type UrlModule struct {}

func NewUrlModule(ctx *js.JSEnv) interface{} {
	return &UrlModule{}
}

func (m *UrlModule) toResult(urlStr string, u *url.URL) (map[string]interface{}, error) {
	q, e := m.ParseQuery(u.RawQuery)
	if e != nil {
		q = nil
	}

	return map[string]interface{}{
		"href":     urlStr,
		"protocol": u.Scheme,
		"host":     u.Host,
		"auth":     u.User.String(),
		"path":     u.Path,
		"search":   u.RawQuery,
		"query":    q,
		"hash":     u.Fragment,
	}, nil
}

func (m *UrlModule) Parse(urlStr string) (map[string]interface{}, error) {
	if urlStr == "" {
		return nil, fmt.Errorf("please specify a url to parse")
	}
	u, e := url.Parse(urlStr)
	if e != nil {
		return nil, e
	}
	return m.toResult(urlStr, u)
}

func (m *UrlModule) ParseRequestURI(uri string) (map[string]interface{}, error) {
	u, e := url.ParseRequestURI(uri)
	if e != nil {
		return nil, e
	}
	return m.toResult(uri, u)
}

func (m *UrlModule) ParseQuery(query string) (map[string]interface{}, error) {
	vv, e := url.ParseQuery(query)
	if e != nil {
		return nil, e
	}
	r := make(map[string]interface{})
	for k, v := range vv {
		if v == nil || len (v) == 0 {
			r[k] = nil
		} else if len(v) == 1 {
			r[k] = v[0]
		} else {
			r[k] = v
		}
	}
	return r, nil
}

