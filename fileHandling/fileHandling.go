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

	var testListString = utils.BytesToString(testListData)

	return strings.Split(testListString, "\n"), nil
}

func SaveTestList(list []string) error {
	var data = strings.Join(list, "\n")
	return SaveFile(path.Join(testFilesDir, testListFileName), []byte(data))
}

func AddTestToList(name string) error {
	var list, err = GetTestsList()
	if err != nil {
		log.DebugFmt("Cannot load file with list of tests\n%s", err.Error())
		return err
	}

	var resultList = append(list, name)

	err = SaveTestList(resultList)
	if err != nil {
		log.DebugFmt("Cannot save file with list of tests\n%s", err.Error())
	}
	return err
}

func ClearTestList() error {
	return SaveFile(path.Join(testFilesDir, testListFileName), []byte(""))
}
