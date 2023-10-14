package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func getFileExtension(language string) string {
	switch language {
	case "python":
		return ".py"
	case "c":
		return ".c"
	case "cpp":
		return ".cpp"
	case "javascript":
		return ".js"
	default:
		return ".py"
	}
}

func CreateFile(code string, language string, dir string) (string, string) {
	// create the file name
	runId := fmt.Sprintf("%s%d", "run", rand.Intn(100000))
	folderPath := filepath.Join(dir)
	filePath := filepath.Join(folderPath, "main"+getFileExtension(language))

	// create the runs directory if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0777)
	}

	// create the file
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	// close the file when the function returns
	defer file.Close()

	// write the code to the file
	_, err = file.WriteString(code)
	if err != nil {
		panic(err)
	}

	return filePath, runId
}
