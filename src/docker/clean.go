package docker

import (
	"context"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func clean(cli *client.Client, response container.CreateResponse, ctx context.Context, fileName string) {
	timeout := 600

	// stop the container
	if err := cli.ContainerStop(ctx, response.ID, container.StopOptions{
		Signal:  "SIGTERM",
		Timeout: &timeout,
	}); err != nil {
		panic(err)
	}

	// // remove the container
	// if err := cli.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{}); err != nil {
	// 	panic(err)
	// }

	// remove the file
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
}
