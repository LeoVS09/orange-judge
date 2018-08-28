package executer

import (
	"bytes"
	"orange-judge/log"
	"os/exec"
	"runtime"
	"strings"
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
	var name = "./"
	if runtime.GOOS == "windows" {
		name += inputFileName + ".exe"
	} else {
		name += inputFileName
	}
	cmd := exec.Command(name)
	cmd.Stdin = strings.NewReader("stdin input read and write stdout")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return &out, err
}

func RunWithOAR(inputFileName string, input string, output string, errorFile string) (*bytes.Buffer, error) {

	//cmd := exec.Command("./oar", "-i", input, "-o", output, "-e", errorFile, "-D", inputFileName)
	cmd := exec.Command("./oar", "-D", "~DEBUG", inputFileName)
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

func RunFromSource(inputFileName string) (*bytes.Buffer, error) {
	var out, err = CompileFrom(inputFileName)
	if err != nil {
		return out, err
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", inputFileName, out.String())

	out, err = Run(inputFileName)
	if err != nil {
		return out, err
	}
	log.DebugFmt("Was run file: %s\nOutput of program:\n%s", inputFileName, out.String())

	return out, nil
}

func RunFromSourceWithOAR(inputFileName string, input string, output string, errorFile string) (*bytes.Buffer, error) {
	var out, err = CompileFrom(inputFileName)
	if err != nil {
		return out, err
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", inputFileName, out.String())

	out, err = RunWithOAR(inputFileName, input, output, errorFile)
	if err != nil {
		return out, err
	}
	log.DebugFmt("Was run with oar file: %s\nOutput of program:\n%s", inputFileName, out.String())

	return out, nil
}
