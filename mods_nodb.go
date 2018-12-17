// +build !db

package main

/**
 * list of builtin modules
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 17, 2018
 */

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"github.com/rosbit/gojs/mod/mod_http"
	"github.com/rosbit/gojs/mod/mod_fs"
	"github.com/rosbit/gojs/mod/mod_url"
)

type fnNewModule func (*js.JSEnv) interface{}

var (
	mods = map[string]fnNewModule {
		"http": mod_http.NewHttpModule,
		"fs":   mod_fs.NewFsModule,
		"url":  mod_url.NewUrlModule,
	}
)

