package executer

import (
	"bytes"
	"orange-judge/log"
	"os/exec"
	"path"
)

func Compile(inputFile string, outputFile string) (*bytes.Buffer, error) {
	var sourceName = path.Join(uploadedFilesDir, inputFile+".cpp")
	var compiledName = path.Join(compiledFilesDir, outputFile)
	log.DebugFmt("Compile %s to %s", sourceName, compiledName)

	cmd := exec.Command("g++", "-std=c++14", sourceName, "-o", path.Join(compiledFilesDir, outputFile))

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	return &out, err
}

func CompileFrom(inputFileName string) (*bytes.Buffer, error) {
	return Compile(inputFileName, inputFileName)
}
