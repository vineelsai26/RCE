package docker

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func CreateFile(code string, language string, RUNS_DIR string) string {
	// create the runs directory if it doesn't exist
	if _, err := os.Stat(RUNS_DIR); os.IsNotExist(err) {
		os.Mkdir(RUNS_DIR, 0777)
	}

	// create the file name
	fileName := fmt.Sprintf("%s%d%s", "run", rand.Intn(100000), getFileExtension(language))

	// create the file path
	filePath := filepath.Join(RUNS_DIR, fileName)

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

	return filePath
}
