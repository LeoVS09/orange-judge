package fileHandling

import (
	"io/ioutil"
	"orange-judge/log"
	"orange-judge/utils"
	"path"
	"strings"
)

// TODO: this constans must be in config
const uploadedFilesDir = "uploaded"
const testFilesDir = "tests"
const testListFileName = "tests.txt"

func SaveUploaded(dir, extension string, body []byte) (string, error) {
	var name = utils.GenerateHash(5)
	var err = SaveFile(path.Join(dir, name+"."+extension), body)
	return name, err
}

func SaveSourceFile(body []byte) (string, error) {
	return SaveUploaded(uploadedFilesDir, "cpp", body)
}

func SaveTestFile(body []byte) (string, error) {
	return SaveUploaded(testFilesDir, "txt", body)
}

func GetTest(fileName string) ([]string, error) {
	var data, err = LoadFile(path.Join(testFilesDir, fileName+".txt"))
	if err != nil {
		log.DebugFmt("Cannot read test file")
		return nil, err
	}

	var test = utils.BytesToString(data)
	return strings.Split(test, "\n"), nil
}

func LoadPage(name string) ([]byte, error) {
	return LoadFile(name + ".html")
}

func LoadFile(name string) ([]byte, error) {
	body, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func SaveFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0600)
}

func GetTestsList() ([]string, error) {
	testListData, err := LoadFile(path.Join(testFilesDir, testListFileName))
	if err != nil {
		return nil, err
	}
	log.DebugFmt("List of tests names was read")

	var testListString = utils.BytesToString(testListData)

	result, err := strings.Split(testListString, "\n"), nil
	if err != nil {
		return nil, err
	}

	log.DebugFmt("List of tests %v", result)
	log.DebugFmt("List of tests length: %v", len(result))
	if len(result) == 1 && (result[0] == "" || result[0] == "\n" || result[0] == "\t" || result[0] == "\r" || result[0] == "\t\r") {
		log.Debug("Return default value")
		return make([]string, 0), nil
	}
	if len(result) > 0 {
		log.DebugFmt("First test length: %v", len(result[0]))
	}

	return result, nil
}

func SaveTestList(list []string) error {
	log.DebugFmt("List of test length: %v", len(list))
	var data = strings.Join(list, "\n")
	log.DebugFmt("List of test value: %v", data)
	return SaveFile(path.Join(testFilesDir, testListFileName), []byte(data))
}

func AddTestToList(name string) error {
	var list, err = GetTestsList()
	if err != nil {
		log.DebugFmt("Cannot load file with list of tests\n%s", err.Error())
		return err
	}
	log.DebugFmt("List of test length: %v", len(list))

	var resultList = append(list, name)

	log.DebugFmt("List of test length: %v", len(resultList))
	err = SaveTestList(resultList)
	if err != nil {
		log.DebugFmt("Cannot save file with list of tests\n%s", err.Error())
	}
	return err
}

func ClearTestList() error {
	return SaveFile(path.Join(testFilesDir, testListFileName), []byte(""))
}
