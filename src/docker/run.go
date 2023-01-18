package docker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Run(RUNS_DIR string, code string, language string) []byte {
	ctx := context.Background()

	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

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

	return output
}
