package main

import (
	"fmt"
	"os"
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	sc "github.com/rosbit/gojs/server_counter"
)

var (
	jsEnv *js.JSEnv
	exit chan bool
)

func main() {
	nArgs := len(os.Args)
	if nArgs == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <js-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s -e <js-code>\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-e" {
		if nArgs == 2 {
			fmt.Fprintf(os.Stderr, "Usage: %s -e <js-code>\n", os.Args[0])
			os.Exit(2)
		}
		evalCode()
		return
	}

	evalFile()
}

func evalCode() {
	jsEnv = js.NewEnv(&MinNodeModuleLoader{})
	defer jsEnv.Destroy()

	if res, err := jsEnv.Eval(os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	} else if res != nil {
		fmt.Fprintf(os.Stderr, "==> %v\n", res)
	}
}

func evalFile() {
	exit = make(chan bool)
	go sc.CountServer(exit)

	jsEnv = js.NewEnv(&MinNodeModuleLoader{})
	defer jsEnv.Destroy()

	if res, err := jsEnv.EvalFile(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	} else if res != nil {
		fmt.Fprintf(os.Stderr, "==> %v\n", res)
	}

	if sc.ServerRunning() {
		<-exit
	}
}
