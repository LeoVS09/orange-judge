package executer

import (
	"bytes"
	"io"
	"orange-judge/fileHandling"
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

func RunAndTest(inputFileName string, input string, outputResult string) (bool, error) {
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

func CompileAndTest(fileName string) (bool, int, error) {
	var out, err = CompileFrom(fileName)
	if err != nil {
		log.DebugFmt("Cannot compile %s, compiler return:\n%s", fileName, out.String())
		return false, 0, err
	}
	log.DebugFmt("Compiled file: %s\nOutput of compiler:\n%s", fileName, out.String())

	testNames, err := fileHandling.GetTestsList()
	if err != nil {
		log.DebugFmt("Cannot load list of tests\n%s", err.Error())
		return false, 0, err
	}
	log.DebugFmt("List of tests:\n%s", testNames)

	for i, testName := range testNames {
		testData, err := fileHandling.GetTest(testName)
		if err != nil {
			log.DebugFmt("Error get test: %s", testName)
			return false, 0, err
		}
		log.DebugFmt("Test %v:\n%v", i, testData)

		var input, output = testData[0], testData[1]
		resultOfTest, err := RunAndTest(fileName, input, output)
		if err != nil {
			log.DebugFmt("Error when test %s, on test:%s", fileName, testName)
			return false, 0, err
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
