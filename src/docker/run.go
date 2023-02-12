package docker

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Run(filePath string, language string) []byte {
	ctx := context.Background()

	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Negotiate to use compatable Docker API version
	cli.NegotiateAPIVersion(ctx)

	runsDir := ""

	for _, folder := range strings.Split(filePath, "/")[0 : len(strings.Split(filePath, "/"))-1] {
		runsDir += folder + "/"
	}

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
				runsDir + ":" + runsDir,
			},
			RestartPolicy: container.RestartPolicy{
				Name: "no",
			},
			Resources: container.Resources{
				Memory: 1024 * 1024 * 512, // 512 MB
			},
			// AutoRemove: true,
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

	timeout := 600

	// stop the container
	if err := cli.ContainerStop(ctx, response.ID, container.StopOptions{
		Signal:  "SIGTERM",
		Timeout: &timeout,
	}); err != nil {
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
	defer clean(cli, response, ctx, filePath)

	// ignore first 8 bits of nonsense
	ignore := make([]byte, 8)
	out.Read(ignore)

	// read the rest of the output
	output, err := io.ReadAll(out)
	if err != nil {
		panic(err)
	}

	return output
}
