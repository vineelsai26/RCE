package docker

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// getDockerImage returns the docker image required to run code for the given language
func getDockerImage(language string, cli *client.Client, ctx context.Context) string {
	images, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	var tags []string = []string{}

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			tags = append(tags, repoTag)
		}
	}

	var container_image string

	switch language {
	case "python":
		container_image = "vineelsai/python"
	case "c":
		container_image = "vineelsai/gcc"
	case "cpp":
		container_image = "vineelsai/gcc"
	case "javascript":
		container_image = "vineelsai/nodejs"
	default:
		container_image = "vineelsai/python"
	}

	for _, tag := range tags {
		if tag == container_image {
			return container_image
		}
	}

	if reader, err := cli.ImagePull(ctx, container_image, image.PullOptions{}); reader != nil {
		if err != nil {
			panic(err)
		}

		defer reader.Close()
		io.Copy(os.Stdout, reader)

		return container_image
	} else {
		panic(err)
	}
}

// getRunCommand returns the command to run the code for the given language
func getRunCommand(language string, filePath string) []string {
	// split the file path into the file name including path and the extension
	path := strings.Split(filePath, ".")

	// return the command to run the code
	switch language {
	case "python":
		return []string{"python", filePath}
	case "c":
		return []string{"bash", "-c", "gcc " + filePath + " -o " + path[0] + " && " + path[0]}
	case "cpp":
		return []string{"bash", "-c", "g++ " + filePath + " -o " + path[0] + " && " + path[0]}
	case "javascript":
		return []string{"node", filePath}
	default:
		return []string{"python", filePath}
	}
}
