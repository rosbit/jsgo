package main

import (
	"fmt"
	"os"
	"github.com/jessevdk/go-flags"
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	sc "github.com/rosbit/gojs/server_counter"
)

var (
	jsEnv *js.JSEnv
	exit chan bool

	buildTime string
	osInfo    string
	goInfo    string
)

var options struct {
	Version  bool   `short:"v" long:"version" description:"Print gojs version"`
	Check    string `short:"c" long:"check" description:"Syntax check script without executing"`
	Eval     string `short:"e" long:"eval" description:"Evaluate script"`
}

func main() {
	parser := flags.NewParser(&options, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	if options.Version {
		showVersion()
		os.Exit(0)
	}

	if options.Check != "" {
		os.Exit(checkSyntax(options.Check))
	}

	if options.Eval != "" {
		os.Exit(evalCode(options.Eval, args))
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stderr)
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(2)
	}

	os.Exit(evalFile(args))
}

func evalCode(jsCode string, args []string) int {
	jsEnv = js.NewEnv(&MinNodeModuleLoader{})
	defer jsEnv.Destroy()

	setEnv(jsEnv, args)
	if res, err := jsEnv.Eval(jsCode); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	} else if res != nil {
		fmt.Fprintf(os.Stderr, "==> %v\n", res)
	}
	return 0
}

func setEnv(jsEnv *js.JSEnv, args []string) {
	envs := os.Environ()
	jsEnv.RegisterVar("ENVIRON", envs)
	jsEnv.RegisterGoFunc("getenv", os.Getenv)
	jsEnv.RegisterVar("ARGV", args)
}

func evalFile(args []string) int {
	exit = make(chan bool)
	defer close(exit)

	go sc.CountServer(exit)

	jsEnv = js.NewEnv(&MinNodeModuleLoader{})
	defer jsEnv.Destroy()

	setEnv(jsEnv, args)
	if res, err := jsEnv.EvalFile(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "error to run %s: %v\n", args[0], err)
		return 1
	} else if res != nil {
		fmt.Fprintf(os.Stderr, "==> %v\n", res)
	}

	if sc.ServerRunning() {
		<-exit
	}
	return 0
}

func showVersion() {
	for prompt, info := range map[string]string{
		"os name": osInfo,
		"compiler": goInfo,
		"buildTime": buildTime,
	} {
		if info != "" {
			fmt.Printf("%11s: %s\n", prompt, info)
		}
	}
}

func checkSyntax(jsFile string) int {
	jsEnv = js.NewEnv(&MinNodeModuleLoader{})
	defer jsEnv.Destroy()

	err := jsEnv.SyntaxCheckFile(jsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "syntax ok\n")
	return 0
}
