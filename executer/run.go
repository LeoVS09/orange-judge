package executer

import (
	"bytes"
	"fmt"
	"io"
	"orange-judge/configuration"
	"orange-judge/log"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"
)

func Run(inputFileName string, in *io.Reader) (*bytes.Buffer, error) {
	var name = "./"
	var out bytes.Buffer

	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	if runtime.GOOS == "windows" {
		name += inputFileName + ".exe"
	} else {
		name += inputFileName
	}

	cmd := exec.Command(path.Join(config.Directories.Compiled, name))

	cmd.Stdin = *in
	cmd.Stdout = &out

	err = cmd.Run()
	return &out, err
}

func testRun(input, inputFileName string) (*bytes.Buffer, error) {
	var reader = *strings.NewReader(input)
	var result io.Reader = &reader
	return Run(inputFileName, &result)
}

func RunWithOAR(inputFileName string, in *io.Reader) (*bytes.Buffer, int, error) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	var name = path.Join(config.Directories.Compiled, inputFileName)
	log.LogFmt("Run oar file: %s", name)

	cmd := exec.Command(
		"./oar",
		"--debug",
		fmt.Sprintf("-c=%v", 4000),
		fmt.Sprintf("-m=%v", 256000),
		fmt.Sprintf("-t=%v", 10000),
		name,
	)
	// Just magic
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}

	cmd.Stdin = *in
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	status, err := runCommand(cmd)

	if err != nil {
		log.DebugFmt("oar output:\n%s", out.String())
	}

	log.DebugFmt("oar error:\n%s", errOut.String())
	return &out, status, err
}

func runCommand(cmd *exec.Cmd) (int, error) {
	var err = cmd.Start()
	if err != nil {
		log.Error("Error cmd start", err)
		return 0, err
	}

	err = cmd.Wait()
	if err == nil {
		return 0, nil
	}

	exiterr, ok := err.(*exec.ExitError)

	if !ok {
		log.Error("Error cmd wait", err)
		return 0, err
	}

	status, ok := exiterr.Sys().(syscall.WaitStatus)

	if !ok {
		log.Error("Error cmd wait", err)
		return 0, err
	}

	return status.ExitStatus(), nil
}
