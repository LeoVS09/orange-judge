package executer

import (
	"bytes"
	"orange-judge/configuration"
	"orange-judge/log"
	"os/exec"
	"path"
)

func Compile(inputFile string, outputFile string) (*bytes.Buffer, error) {
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	var sourceName = path.Join(config.Directories.Uploaded, inputFile+".cpp")
	var compiledName = path.Join(config.Directories.Compiled, outputFile)
	log.DebugFmt("Compile %s to %s", sourceName, compiledName)

	cmd := exec.Command("g++", "-std=c++14", sourceName, "-o", path.Join(config.Directories.Compiled, outputFile))

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	return &out, err
}

func CompileFrom(inputFileName string) (*bytes.Buffer, error) {
	return Compile(inputFileName, inputFileName)
}
