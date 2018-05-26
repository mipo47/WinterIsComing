package main

import (
	"./core"
	"os"
	"fmt"
	"strings"
	"runtime"
)

const IS_WINDOWS = runtime.GOOS == "windows"

func showHelp()  {
	fmt.Println(`
Usage: 
	go run run.go [command]

The commands are:
	test      runs all tests in project
	build     builds server and client applications to 'out' directory
	server    builds and runs server app
	client    builds and runs client app
	web       builds and runs client app with integrated website
`)
}

func build(name, folder string) string {
	if IS_WINDOWS {
		name += ".exe"
	}
	core.Execute("go build -i -o ./out/" + name + " " + folder)
	fmt.Println("Build successful: out/" + name)
	return name
}

func copyResources()  {
	if IS_WINDOWS {
		core.Execute("xcopy /f /y client_web\\html out\\html")
	} else {
		core.Execute("cp -a ./client_web/html ./out")
	}
}

func main() {
	if len(os.Args) != 2 {
		showHelp()
		return
	}

	var run string
	switch strings.ToLower(os.Args[1]) {
	case "test":
		core.Execute("go test ./...")
	case "build":
		build("server", "./server")
		build("client", "./client")
		build("web", "./client_web")
		copyResources()
	case "server":
		run = build("server", "./server")
	case "client":
		run = build("client", "./client")
	case "web":
		run = build("web", "./client_web")
		copyResources()
	case "help":
		fallthrough
	default: showHelp()
	}

	if run != "" {
		core.ExecuteIn("./" + run, "./out")
	}
}