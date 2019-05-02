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

func TestRunFromSource(input, inputFileName string) (*bytes.Buffer, error) {
	return compileAndRun(inputFileName, func() (*bytes.Buffer, error) {
		return testRun(input, inputFileName)
	})
}

func RunFromSourceWithOAR(inputFileName string, input string, output string, errorFile string) (*bytes.Buffer, error) {
	return compileAndRun(inputFileName, func() (*bytes.Buffer, error) {
		return RunWithOAR(inputFileName, input, output, errorFile)
	})
}

func RunFromSource(inputFileName string, input string) (*bytes.Buffer, error) {
	var reader = *strings.NewReader(input)
	var result io.Reader = &reader
	return compileAndRun(inputFileName, func() (*bytes.Buffer, error) {
		return Run(inputFileName, &result)
	})
}

func TestProgram(inputFileName string, input string, outputResult string) (bool, error) {
	var reader = *strings.NewReader(input)
	var inputData io.Reader = &reader
	var outputData, err = Run(inputFileName, &inputData)
	if err != nil {
		log.DebugFmt("Cannot run file %s, with input: %s", inputFileName, input)
		return false, err
	}

	var output = utils.RemoveUnnecessarySymbols(outputData.String())
	var result = output == outputResult
	log.DebugFmt("Result of compare test (%s) with real result (%s): %v", outputData, outputResult, result)
	return result, nil
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
		}
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", fileName, out.String())

	for i, test := range tests {
		log.DebugFmt("Test %v:\n%v", i, test)

		var input, output = test.Input, test.Output
		resultOfTest, err := TestProgram(fileName, input, output)
		if err != nil {
			log.DebugFmt("Error when test %s, on test:%s", fileName, test.Id)
			return false, 0, ExecutionError{
				Stage:     Testing,
				TestIndex: i,
				RealError: err,
			}
		}
		log.DebugFmt("Result of test %v: %v", i, resultOfTest)

		if resultOfTest == false {
			return false, i, nil
		}
	}

	return true, 0, nil
}

func compileAndRun(fileName string, runner func() (*bytes.Buffer, error)) (*bytes.Buffer, error) {
	var out, err = CompileFrom(fileName)
	if err != nil {
		log.DebugFmt("Cannot compile %s\nOutput of compiler:\n%s", fileName, out.String())
		return out, err
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", fileName, out.String())

	out, err = runner()
	if err != nil {
		log.DebugFmt("Cannot run file: %s\nOutput of program:\n%s", fileName, out.String())
		return out, err
	}
	log.DebugFmt("Was run file: %s\nOutput of program:\n%s", fileName, out.String())

	return out, nil
}
