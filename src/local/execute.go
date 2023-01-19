package local

import (
	"path/filepath"
	"rce/src/docker"
)

func Execute(fileName string, language string) string {
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		panic(err)
	}

	output := docker.Run(filePath, language)

	return string(output)
}
