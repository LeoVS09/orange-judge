package executer

import (
	"bytes"
	"orange-judge/log"
	"os/exec"
)

func Compile(inputFile string, outputFile string) (*bytes.Buffer, error) {
	cmd := exec.Command("g++", "-std=c++14", inputFile+".cpp", "-o", outputFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return &out, err
}

func CompileFrom(inputFileName string) (*bytes.Buffer, error) {
	return Compile(inputFileName, inputFileName)
}

func Run(inputFileName string) (*bytes.Buffer, error) {
	cmd := exec.Command("./" + inputFileName + ".exe")
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return &out, err
}

func RunFromSource(inputFileName string) (*bytes.Buffer, error) {
	var out, err = CompileFrom(inputFileName)
	if err != nil {
		return out, err
	}
	log.DebugFmt("Compiled file: %s\nOuput of compiler:\n%s", inputFileName, out.String())

	out, err = Run(inputFileName)
	if err != nil {
		return out, err
	}
	log.DebugFmt("Was run file: %s\nOuput of program:\n%s", inputFileName, out.String())

	return out, nil
}
