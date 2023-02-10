package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func PullImages() {
	ctx := context.Background()

	// Connect to the Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Negotiate to use compatable Docker API version
	cli.NegotiateAPIVersion(ctx)

	// pull the images
	if reader, err := cli.ImagePull(ctx, "vineelsai/python", types.ImagePullOptions{}); reader != nil {
		if err != nil {
			panic(err)
		}

		defer reader.Close()
		io.Copy(os.Stdout, reader)
	} else {
		panic(err)
	}

	if reader, err := cli.ImagePull(ctx, "vineelsai/gcc", types.ImagePullOptions{}); reader != nil {
		if err != nil {
			panic(err)
		}

		defer reader.Close()
		io.Copy(os.Stdout, reader)
	} else {
		panic(err)
	}

	if reader, err := cli.ImagePull(ctx, "vineelsai/nodejs", types.ImagePullOptions{}); reader != nil {
		if err != nil {
			panic(err)
		}

		defer reader.Close()
		io.Copy(os.Stdout, reader)
	} else {
		panic(err)
	}
}
