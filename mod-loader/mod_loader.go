package ml

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"fmt"
)

/**
 * a module loader.
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 6, 2018
 */

func ListModules() {
	fmt.Printf("Built-in moduels in jsgo\n")
	i := 0
	for m, _ := range mods {
		i++
		fmt.Printf(" #%d: %s\n", i, m)
	}
}

// ------------- implement interface GoModuleLoader -----------------
type MinNodeModuleLoader struct {
	jsEnv *js.JSEnv
}

func (loader *MinNodeModuleLoader) SetJSEnv(jsEnv *js.JSEnv) {
	loader.jsEnv = jsEnv
}

func (loader *MinNodeModuleLoader) GetExtName() string {
	return ".loader_for_builtin_modules"
}

func (loader *MinNodeModuleLoader) LoadModule(modHome string, modName string) interface{} {
	if fn, ok := mods[modName]; !ok {
		return nil
	} else {
		return fn(loader.jsEnv)
	}
}

func (loader *MinNodeModuleLoader) FinalizeModule(modName string, modHandler interface{}) {
}
