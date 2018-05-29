package main

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"os"
	"fmt"
	"strings"
	"runtime"
	"log"
)

const IS_WINDOWS = runtime.GOOS == "windows"

func showHelp()  {
	fmt.Println(`
Usage: 
	go run run.go [command]

The commands are:
	test      runs all tests in project
	build     builds applications to 'out' directory
	server    builds and runs server app
	client    builds and runs client app
	web       builds and runs client app with integrated website
`)
}

func build(binaryName, packageName string) {
	if IS_WINDOWS {
		binaryName += ".exe"
	}
	core.MustRunCommand("go", "build", "-i", "-o", "./out/" +binaryName, packageName)
	fmt.Printf("Build successful: out/%s\n", binaryName)
}

func buildServer() {
	build("server", "./server")
}

func buildServerImage() {
	buildDocker("docker/server/Dockerfile", "mysteriumnetwork/winter-server")
}

func buildClient() {
	build("client", "./client")
}

func buildWeb() {
	build("web", "./client_web")
	if IS_WINDOWS {
		core.Execute("xcopy /f /y client_web\\html out\\html")
	} else {
		core.Execute("cp -a ./client_web/html ./out")
	}
}

func buildDocker(dockerfile, image string) {
	fmt.Printf("Building image '%s'..\n", image)

	core.MustRunCommand("docker", "build", "--file", dockerfile, "--tag", image, ".")
	fmt.Print("Docker image building process complete!")
}

func doBuild(args []string) {
	switch strings.ToLower(args[0]) {
	case "server":
		buildServer()
	case "client":
		buildClient()
	case "web":
		buildWeb()
	case "server-image":
		buildServerImage()
	default:
		buildServer()
		buildClient()
		buildWeb()
		buildServerImage()
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need subcommand as first argument")
	}

	switch strings.ToLower(os.Args[1]) {
	case "test":
		core.Execute("go test ./...")
	case "build":
		doBuild(os.Args[2:])
	case "server":
		buildServer()
		core.ExecuteIn("./server", "./out")
	case "client":
		buildClient()
		core.ExecuteIn("./client", "./out")
	case "web":
		buildWeb()
		core.ExecuteIn("./web", "./out")
	case "help":
		fallthrough
	default:
		showHelp()
	}
}
