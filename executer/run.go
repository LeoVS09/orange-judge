package executer

import (
	"bytes"
	"io"
	"orange-judge/log"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func Run(inputFileName string, in *io.Reader) (*bytes.Buffer, error) {
	var name = "./"
	var out bytes.Buffer

	if runtime.GOOS == "windows" {
		name += inputFileName + ".exe"
	} else {
		name += inputFileName
	}

	cmd := exec.Command(path.Join(compiledFilesDir, name))

	cmd.Stdin = *in
	cmd.Stdout = &out

	err := cmd.Run()

	return &out, err
}

func testRun(inputFileName string) (*bytes.Buffer, error) {
	var reader = *strings.NewReader("stdin input read and write stdout")
	var result io.Reader = &reader
	return Run(inputFileName, &result)
}

func RunWithOAR(inputFileName string, input string, output string, errorFile string) (*bytes.Buffer, error) {

	//cmd := exec.Command("./oar", "-i", input, "-o", output, "-e", errorFile, "-D", inputFileName)
	cmd := exec.Command("./oar", "-D", "~DEBUG", path.Join(compiledFilesDir, inputFileName))
	cmd.Stdin = strings.NewReader("stdin input read and write stdout")
	var out bytes.Buffer
	cmd.Stdout = &out
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		log.DebugFmt("oar output:\n%s", out.String())
		log.DebugFmt("oar error:\n%s", errOut.String())
	}
	return &out, err
}
