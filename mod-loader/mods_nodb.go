// +build !db

package ml

/**
 * list of builtin modules
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 17, 2018
 */

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"github.com/rosbit/jsgo/mods/mod_http"
	"github.com/rosbit/jsgo/mods/mod_fs"
	"github.com/rosbit/jsgo/mods/mod_url"
)

type fnNewModule func (*js.JSEnv) interface{}

var (
	mods = map[string]fnNewModule {
		"http": mod_http.NewHttpModule,
		"fs":   mod_fs.NewFsModule,
		"url":  mod_url.NewUrlModule,
	}
)

