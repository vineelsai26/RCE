package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
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

func createFile(code string, language string) string {
	if _, err := os.Stat("/usr/src/app/runs"); os.IsNotExist(err) {
		os.Mkdir("/usr/src/app/runs", 0777)
	}

	fileName := fmt.Sprintf("/usr/src/app/runs/%s%d%s", "run", rand.Intn(100000), getFileExtension(language))
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		panic(err)
	}

	return fileName
}

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

func getRunCommand(language string, fileName string) []string {
	path := strings.Split(fileName, ".")

	switch language {
	case "python":
		return []string{"python", fileName}
	case "c":
		return []string{"bash", "-c", "gcc " + fileName + " -o " + path[0] + " && " + path[0]}
	case "cpp":
		return []string{"bash", "-c", "g++ " + fileName + " -o " + path[0] + " && " + path[0]}
	default:
		return []string{"python", fileName}
	}
}

func clean(cli *client.Client, response container.ContainerCreateCreatedBody, ctx context.Context, fileName string) {
	// stop the container
	if err := cli.ContainerStop(ctx, response.ID, nil); err != nil {
		panic(err)
	}

	// remove the container
	if err := cli.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	// remove the file
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
}

func run(res http.ResponseWriter, req *http.Request, ctx context.Context) {
	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	code := req.FormValue("code")
	language := req.FormValue("language")
	fileName := createFile(code, language)

	// Create the container
	response, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           getDockerImage(language),
			Cmd:             getRunCommand(language, fileName),
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Binds: []string{
				"/usr/src/app/runs:/usr/src/app/runs",
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

	// get the logs from the container
	out, err := cli.ContainerLogs(ctx, response.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		panic(err)
	}

	// ignore first 8 bits of nonsense
	ignore := make([]byte, 8)
	out.Read(ignore)

	buf, err := io.ReadAll(out)
	if err != nil {
		panic(err)
	}

	output, err := json.Marshal(map[string]string{
		"output": string(buf),
	})

	if err != nil {
		panic(err)
	}

	go clean(cli, response, ctx, fileName)

	res.Header().Add("Content-Type", "application/json")
	res.Write(output)
}

func main() {
	godotenv.Load()
	ctx := context.Background()
	PORT := os.Getenv("PORT")

	http.HandleFunc("/run", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			run(res, req, ctx)
		case http.MethodGet:
			res.Write([]byte("Hello World"))
		}
	})

	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		panic(err)
	}
}
