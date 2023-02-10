package docker

import (
	"fmt"
	"testing"
)

func TestPullImage(t *testing.T) {
	PullImages()
}

func TestRunPython(t *testing.T) {
	code := "print('Hello World')"
	language := "python"
	RUNS_DIR := "/usr/src/app/runs"

	filePath := CreateFile(code, language, RUNS_DIR)

	output := Run(filePath, language)

	if string(output) != string("Hello World\n") {
		t.Error("Output not Matched")
	}
}

func TestRunPythonSleep(t *testing.T) {
	code := "import time\ntime.sleep(10)\nprint(\"Hello World\")"
	language := "python"
	RUNS_DIR := "/usr/src/app/runs"

	filePath := CreateFile(code, language, RUNS_DIR)

	output := Run(filePath, language)

	if string(output) != string("Hello World\n") {
		t.Error("Output not Matched")
	}
}

func TestRunC(t *testing.T) {
	code := "#include <stdio.h>\n int main() {\n	printf(\"Hello World\");\n    return 0;\n}"
	language := "c"
	RUNS_DIR := "/usr/src/app/runs"

	filePath := CreateFile(code, language, RUNS_DIR)

	output := Run(filePath, language)

	if string(output) != string("Hello World") {
		fmt.Println(string(output))
		t.Error("Output not Matched")
	}
}

func TestRunCpp(t *testing.T) {
	code := "#include <iostream>\n int main() {\n    std::cout << \"Hello World\";\n    return 0;\n}"
	language := "cpp"
	RUNS_DIR := "/usr/src/app/runs"

	filePath := CreateFile(code, language, RUNS_DIR)

	output := Run(filePath, language)

	if string(output) != string("Hello World") {
		fmt.Println(string(output))
		t.Error("Output not Matched")
	}
}
