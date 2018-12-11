package main

/**
 * a module loader.
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 6, 2018
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

// ------------- implement interface GoModuleLoader -----------------
type MinNodeModuleLoader struct {}

func (loader *MinNodeModuleLoader) GetExtName() string {
	return ".loader_for_only_http_module"
}

func (loader *MinNodeModuleLoader) LoadModule(modHome string, modName string) interface{} {
	if fn, ok := mods[modName]; !ok {
		return nil
	} else {
		return fn(jsEnv)
	}
}

func (loader *MinNodeModuleLoader) FinalizeModule(modName string, modHandler interface{}) {
}
