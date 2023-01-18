package main

import (
	"fmt"
	"os"
	"strings"

	"rce/src/api"
	"rce/src/docker"
)

func help() {
	fmt.Println("Usage: rce [OPTIONS]")
	fmt.Println("--version - prints the version")
	fmt.Println("--pull-images - pulls the docker images required for the code executions to run")
	fmt.Println("--runs-dir=RUNS_DIR - sets the directory where the code files will be stored (default: /usr/src/app/runs)")
	fmt.Println("--port=PORT - sets the port where the server will run on (default: 3000)")
}

func main() {
	PORT := "3000"
	RUNS_DIR := "/usr/src/app/runs"
	VERSION := "1.1.1"

	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]

	// parse the command line arguments
	for _, arg := range args {
		if arg == "--help" {
			help()
			return
		}
		if arg == "--version" {
			fmt.Println(VERSION)
			return
		}
		if arg == "--pull-images" {
			docker.PullImages()
			return
		}
		if strings.Contains(arg, "--runs-dir=") {
			RUNS_DIR = strings.Split(arg, "=")[1]
		}
		if strings.Contains(arg, "--port=") {
			PORT = strings.Split(arg, "=")[1]
		}
	}

	if args[0] == "server" {
		api.Serve(PORT, RUNS_DIR)
	} else {
		help()
	}
}
