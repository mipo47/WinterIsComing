package main

import (
	"github.com/mysteriumnetwork/winter-server/core"
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

func buildServer() string {
	return build("server", "./server")
}

func buildClient() string {
	return build("client", "./client")
}

func buildWeb() string {
	name := build("web", "./client_web")
	if IS_WINDOWS {
		core.Execute("xcopy /f /y client_web\\html out\\html")
	} else {
		core.Execute("cp -a ./client_web/html ./out")
	}

	return name
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
		buildServer()
		buildClient()
		buildWeb()

	case "server":
		run = buildServer()
	case "client":
		run = buildClient()
	case "web":
		run = buildWeb()

	case "help":
		fallthrough
	default: showHelp()
	}

	if run != "" {
		core.ExecuteIn("./" + run, "./out")
	}
}
