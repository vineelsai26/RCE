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
	if _, err := os.Stat("runs"); os.IsNotExist(err) {
		os.Mkdir("runs", 0777)
	}

	fileName := fmt.Sprintf("runs/%s%d.py", "run", rand.Intn(100000))
	file, err := os.Create(fileName)
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
