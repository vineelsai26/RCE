package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func PullImages() {
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
