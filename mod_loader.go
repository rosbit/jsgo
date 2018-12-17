package main

import (
	"fmt"
)

/**
 * a module loader.
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 6, 2018
 */

func listModules() {
	fmt.Printf("Built-in moduels in gojs\n")
	i := 0
	for m, _ := range mods {
		i++
		fmt.Printf(" #%d: %s\n", i, m)
	}
}

// ------------- implement interface GoModuleLoader -----------------
type MinNodeModuleLoader struct {}

func (loader *MinNodeModuleLoader) GetExtName() string {
	return ".loader_for_builtin_modules"
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
