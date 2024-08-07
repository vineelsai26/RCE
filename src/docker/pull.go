package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
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

	dockerImages := []string{
		"vineelsai/python",
		"vineelsai/gcc",
		"vineelsai/nodejs",
	}

	// pull the images
	for _, container_image := range dockerImages {
		if reader, err := cli.ImagePull(ctx, container_image, image.PullOptions{}); reader != nil {
			if err != nil {
				panic(err)
			}

			defer reader.Close()
			io.Copy(os.Stdout, reader)
		} else {
			panic(err)
		}
	}
}
