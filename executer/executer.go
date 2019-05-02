package executer

import (
	"bytes"
	"fmt"
	"io"
	"orange-judge/database"
	"orange-judge/log"
	"orange-judge/utils"
	"strings"
)

func RunFromSourceWithOAR(inputFileName string, input string) (*bytes.Buffer, int, error) {
	var reader = *strings.NewReader(input)
	var in io.Reader = &reader
	return compileAndRun(inputFileName, func() (*bytes.Buffer, int, error) {
		return RunWithOAR(inputFileName, &in)
	})
}

func TestProgram(inputFileName string, input string, outputResult string) (bool, int, error) {
	var reader = *strings.NewReader(input)
	var inputData io.Reader = &reader
	var outputData, status, err = RunWithOAR(inputFileName, &inputData)
	if err != nil {
		log.DebugFmt("Cannot run file %s, with input: %s", inputFileName, input)
		return false, status, err
	}

	var output = utils.RemoveUnnecessarySymbols(outputData.String())
	var result = output == outputResult
	log.DebugFmt("Result of compare test (%s) with real result (%s): %v", outputData, outputResult, result)
	return result, status, nil
}

type ExecutionStage int

const (
	Compilation ExecutionStage = 0
	Testing     ExecutionStage = 1
)

type ExecutionError struct {
	Stage     ExecutionStage
	TestIndex int
	RealError error
	Status    int
}

func (e ExecutionError) Error() string {
	switch e.Stage {
	case Compilation:
		return fmt.Sprintf("Error on compiltaion: %v", e.RealError)
	case Testing:
		return fmt.Sprintf("Error on testing: %v", e.RealError)
	}

	return e.RealError.Error()
}

func CompileAndTest(fileName string, tests []database.Test) (bool, int, error) {
	var out, err = CompileFrom(fileName)
	if err != nil {
		log.DebugFmt("Cannot compile %s, compiler return:\n%s", fileName, err.Error())
		return false, 0, ExecutionError{
			Stage:     Compilation,
			TestIndex: 0,
			RealError: err,
			Status:    -1,
		}
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", fileName, out.String())

	for i, test := range tests {
		log.DebugFmt("Test %v:\n%v", i, test)

		var input, output = test.Input, test.Output
		resultOfTest, status, err := TestProgram(fileName, input, output)
		if err != nil || status != 0 {
			log.DebugFmt("Error when test %s, on test:%s", fileName, test.Id)
			return false, 0, ExecutionError{
				Stage:     Testing,
				TestIndex: i,
				RealError: err,
				Status:    status,
			}
		}
		log.DebugFmt("Result of test %v: %v", i, resultOfTest)

		if resultOfTest == false {
			return false, i, nil
		}
	}

	return true, 0, nil
}

func compileAndRun(fileName string, runner func() (*bytes.Buffer, int, error)) (*bytes.Buffer, int, error) {
	var out, err = CompileFrom(fileName)
	if err != nil {
		log.DebugFmt("Cannot compile %s\nOutput of compiler:\n%s", fileName, out.String())
		return out, 0, err
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", fileName, out.String())

	out, status, err := runner()
	if err != nil {
		log.DebugFmt("Cannot run file: %s\nOutput of program:\n%s", fileName, out.String())
		return out, status, err
	}
	log.DebugFmt("Was run file: %s\nOutput of program:\n%s", fileName, out.String())

	return out, status, nil
}
