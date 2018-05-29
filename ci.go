package main

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"os"
	"fmt"
	"strings"
	"runtime"
)

const IS_WINDOWS = runtime.GOOS == "windows"

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

func doRun(args []string) {
	help := `
Usage: 
	go run ci.go run <artifact>

The artifacts are:
	server    builds and runs server app
	client    builds and runs client app
	web       builds and runs client app with integrated website
	help
`

	if len(args) < 1 {
		fmt.Println("Need <artifact> as first argument")
		fmt.Print(help)
		os.Exit(1)
	}

	switch strings.ToLower(args[0]) {
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
		fmt.Print(help)
	default:
		buildServer()
		buildClient()
		buildWeb()
		buildServerImage()
	}
}

func doBuild(args []string) {
	help := `
Usage: 
	go run ci.go build <artifact>

The artifacts are:
	server        builds server to 'out' directory
	client        builds client to 'out' directory
	web           builds web to 'out' directory
	server-image  builds Docker image of server
	help
`

	switch strings.ToLower(args[0]) {
	case "server":
		buildServer()
	case "client":
		buildClient()
	case "web":
		buildWeb()
	case "server-image":
		buildServerImage()
	case "help":
		fmt.Print(help)
	default:
		buildServer()
		buildClient()
		buildWeb()
		buildServerImage()
	}
}

func do(args []string) {
	help := `
Usage: 
	go run ci.go <command>

The commands are:
	test      runs all tests in project
	run       run application artifacts
	build     builds application artifacts
	help
`

	if len(args) < 2 {
		fmt.Println("Need subcommand as first argument")
		fmt.Print(help)
		os.Exit(1)
	}

	switch strings.ToLower(args[0]) {
	case "test":
		core.Execute("go test ./...")
	case "run":
		doRun(args[1:])
	case "build":
		doBuild(args[1:])
	case "help":
		fallthrough
	default:
		fmt.Print(help)
	}
}

func main() {
	do(os.Args[1:])
}
