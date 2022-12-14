package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func createFile(req *http.Request) string {
	fileName := fmt.Sprintf("%d.py", rand.Intn(100000))
	file, err := os.Create("runs/" + fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(req.FormValue("code"))
	if err != nil {
		panic(err)
	}

	return fileName
}

func run(res http.ResponseWriter, req *http.Request, ctx context.Context) {
	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	fileName := createFile(req)

	// Create the container
	response, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           "vineelsai/python",
			Cmd:             []string{"python", "/usr/src/app/" + fileName},
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
	})
	if err != nil {
		panic(err)
	}

	buf, err := io.ReadAll(out)

	// remove the container
	if err := cli.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	output, err := json.Marshal(map[string]string{
		"output": string(buf),
	})

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
