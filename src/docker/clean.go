package docker

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func clean(cli *client.Client, response container.CreateResponse, ctx context.Context, fileName string) {
	// remove the container
	if err := cli.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	// remove the file
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
}
