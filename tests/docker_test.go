package tests

import (
	"fmt"
	"testing"

	"vineelsai.com/rce/src/docker"
	"vineelsai.com/rce/src/utils"
)

var RUNS_DIR = "runs"

func TestPullImage(t *testing.T) {
	docker.PullImages()
}

func TestRunPython(t *testing.T) {
	code := "print('Hello World')"
	language := "python"

	filePath, runId := utils.CreateFile(code, language, RUNS_DIR)

	output := docker.Run(filePath, language, runId)

	if string(output) != string("Hello World\n") {
		t.Error("Output not Matched")
	}
}

func TestRunPythonSleep(t *testing.T) {
	code := "import time\ntime.sleep(10)\nprint(\"Hello World\")"
	language := "python"

	filePath, runId := utils.CreateFile(code, language, RUNS_DIR)

	output := docker.Run(filePath, language, runId)

	if string(output) != string("Hello World\n") {
		t.Error("Output not Matched")
	}
}

func TestRunC(t *testing.T) {
	code := "#include <stdio.h>\n int main() {\n	printf(\"Hello World\");\n    return 0;\n}"
	language := "c"

	filePath, runId := utils.CreateFile(code, language, RUNS_DIR)

	output := docker.Run(filePath, language, runId)

	if string(output) != string("Hello World") {
		fmt.Println(string(output))
		t.Error("Output not Matched")
	}
}

func TestRunCpp(t *testing.T) {
	code := "#include <iostream>\n int main() {\n    std::cout << \"Hello World\";\n    return 0;\n}"
	language := "cpp"

	filePath, runId := utils.CreateFile(code, language, RUNS_DIR)

	output := docker.Run(filePath, language, runId)

	if string(output) != string("Hello World") {
		fmt.Println(string(output))
		t.Error("Output not Matched")
	}
}

func TestRunJavaScript(t *testing.T) {
	code := "console.log('Hello World')"
	language := "javascript"

	filePath, runId := utils.CreateFile(code, language, RUNS_DIR)

	output := docker.Run(filePath, language, runId)

	if string(output) != string("Hello World\n") {
		t.Error("Output not Matched")
	}
}
