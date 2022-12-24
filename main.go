package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func getFileExtension(language string) string {
	switch language {
	case "python":
		return ".py"
	case "c":
		return ".c"
	case "cpp":
		return ".cpp"
	default:
		return ".py"
	}
}

func createFile(code string, language string, RUNS_DIR string) string {
	// create the runs directory if it doesn't exist
	if _, err := os.Stat(RUNS_DIR); os.IsNotExist(err) {
		os.Mkdir(RUNS_DIR, 0777)
	}

	// create the file name
	fileName := fmt.Sprintf("%s%d%s", "run", rand.Intn(100000), getFileExtension(language))

	// create the file path
	filePath := filepath.Join(RUNS_DIR, fileName)

	// create the file
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	// close the file when the function returns
	defer file.Close()

	// write the code to the file
	_, err = file.WriteString(code)
	if err != nil {
		panic(err)
	}

	return filePath
}

// getDockerImage returns the docker image required to run code for the given language
func getDockerImage(language string) string {
	switch language {
	case "python":
		return "vineelsai/python"
	case "c":
		return "vineelsai/gcc"
	case "cpp":
		return "vineelsai/gcc"
	default:
		return "vineelsai/python"
	}
}

// getRunCommand returns the command to run the code for the given language
func getRunCommand(language string, filePath string) []string {
	// split the file path into the file name including path and the extension
	path := strings.Split(filePath, ".")

	// return the command to run the code
	switch language {
	case "python":
		return []string{"python", filePath}
	case "c":
		return []string{"bash", "-c", "gcc " + filePath + " -o " + path[0] + " && " + path[0]}
	case "cpp":
		return []string{"bash", "-c", "g++ " + filePath + " -o " + path[0] + " && " + path[0]}
	default:
		return []string{"python", filePath}
	}
}

func clean(cli *client.Client, response container.ContainerCreateCreatedBody, ctx context.Context, fileName string) {
	// remove the container
	if err := cli.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	// remove the file
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
}

func run(res http.ResponseWriter, req *http.Request, ctx context.Context, RUNS_DIR string) {
	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	code := req.FormValue("code")
	language := req.FormValue("language")
	filePath := createFile(code, language, RUNS_DIR)

	// Create the container
	response, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           getDockerImage(language),
			Cmd:             getRunCommand(language, filePath),
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Binds: []string{
				RUNS_DIR + ":" + RUNS_DIR,
			},
			RestartPolicy: container.RestartPolicy{
				Name: "no",
			},
			Resources: container.Resources{
				Memory: 1024 * 1024 * 512, // 512 MB
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}

	// start the container, if it returns an error, print it
	if err := cli.ContainerStart(ctx, response.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// stop the container
	if err := cli.ContainerStop(ctx, response.ID, nil); err != nil {
		panic(err)
	}

	// get the logs from the container
	out, err := cli.ContainerLogs(ctx, response.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		panic(err)
	}

	// clean up the container and the file
	go clean(cli, response, ctx, filePath)

	// ignore first 8 bits of nonsense
	ignore := make([]byte, 8)
	out.Read(ignore)

	// read the rest of the output
	buf, err := io.ReadAll(out)
	if err != nil {
		panic(err)
	}

	// convert the output to json
	output, err := json.Marshal(map[string]string{
		"output": string(buf),
	})
	if err != nil {
		panic(err)
	}

	// set the content type to json and write the output
	res.Header().Add("Content-Type", "application/json")
	res.Write(output)
}

func pullImages() {
	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// pull the images
	if reader, _ := cli.ImagePull(ctx, "vineelsai/python", types.ImagePullOptions{}); reader != nil {
		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}

	if reader, _ := cli.ImagePull(ctx, "vineelsai/gcc", types.ImagePullOptions{}); reader != nil {
		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}
}

func main() {
	godotenv.Load()
	ctx := context.Background()
	PORT := os.Getenv("PORT")
	RUNS_DIR := "/usr/src/app/runs"

	args := os.Args[1:]

	// parse the command line arguments
	for _, arg := range args {
		if arg == "--help" {
			fmt.Println("Usage: rce [OPTIONS]")
			fmt.Println("--version - prints the version")
			fmt.Println("--pull-images - pulls the docker images required for the code executions to run")
			fmt.Println("--runs-dir=RUNS_DIR - sets the directory where the code files will be stored (default: /usr/src/app/runs)")
			fmt.Println("--port=PORT - sets the port where the server will run on (default: 3000)")
			return
		}
		if arg == "--version" {
			fmt.Println("Version: 1.0.0")
			return
		}
		if arg == "--pull-images" {
			pullImages()
			return
		}
		if strings.Contains(arg, "--runs-dir=") {
			RUNS_DIR = strings.Split(arg, "=")[1]
		}
		if strings.Contains(arg, "--port=") {
			PORT = strings.Split(arg, "=")[1]
		}
	}

	// handle the /run endpoint
	http.HandleFunc("/run", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			run(res, req, ctx, RUNS_DIR)
		case http.MethodGet:
			res.Write([]byte("Hello World"))
		}
	})

	// start the server
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		panic(err)
	}
}
