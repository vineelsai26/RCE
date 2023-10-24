package main

import (
	"fmt"
	"os"
	"strings"

	"vineelsai.com/rce/src"
	"vineelsai.com/rce/src/api"
	"vineelsai.com/rce/src/docker"
)

func help() {
	fmt.Println("Usage: rce [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("server - starts the server")
	fmt.Println("-h or --help - prints this help")
	fmt.Println("-v or --version - prints the version")
	fmt.Println("-i or --pull-images - pulls the docker images required for the code executions to run")
	fmt.Println("-p PORT or --port=PORT - sets the port where the server will run on (default: 3000)")
}

func main() {
	PORT := "3000"

	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]

	// parse the command line arguments
	for index, arg := range args {
		if arg == "--help" || arg == "-h" {
			help()
			return
		}
		if arg == "--version" || arg == "-v" {
			src.GetVersion()
			return
		}
		if arg == "--pull-images" || arg == "-i" {
			docker.PullImages()
			return
		}
		if strings.Contains(arg, "--port=") {
			PORT = strings.Split(arg, "=")[1]
		}
		if arg == "-p" {
			PORT = args[index+1]
		}
	}

	if args[0] == "server" {
		api.Serve(PORT)
	} else {
		help()
	}
}
