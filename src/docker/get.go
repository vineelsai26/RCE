package docker

import (
	"strings"
)

// getDockerImage returns the docker image required to run code for the given language
func getDockerImage(language string) string {
	switch language {
	case "python":
		return "vineelsai/python"
	case "c":
		return "vineelsai/gcc"
	case "cpp":
		return "vineelsai/gcc"
	case "javascript":
		return "vineelsai/nodejs"
	default:
		return "vineelsai/python"
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
