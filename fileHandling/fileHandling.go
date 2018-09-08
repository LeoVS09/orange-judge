package fileHandling

import (
	"io/ioutil"
	"orange-judge/configuration"
	"orange-judge/log"
	"orange-judge/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func IsExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		// folder does not exist
		return false
	}
	return true
}

func CreateFolder(folder string) error {
	log.DebugFmt("Create folder %s", folder)
	return os.MkdirAll(folder, os.ModePerm)
}

func ClearFolder(folder string) error {
	d, err := os.Open(folder)
	if err != nil {
		log.DebugFmt("Cannot open for clear folder %s", folder)
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		var filename = filepath.Join(folder, name)
		log.DebugFmt("Remove %s", filename)
		err = os.Remove(filename)
		if err != nil {
			log.DebugFmt("Cannot remove %s", filename)
			return err
		}
	}
	return nil
}

func CreateOrClearFolder(folder string) error {
	if IsExist(folder) {
		return ClearFolder(folder)
	}
	return CreateFolder(folder)
}

func SaveUploaded(dir, extension string, body []byte) (string, error) {
	var name = utils.GenerateHash(5)
	var err = SaveFile(path.Join(dir, name+"."+extension), body)
	return name, err
}

func SaveSourceFile(body []byte) (string, error) {
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	return SaveUploaded(config.Directories.Uploaded, "cpp", body)
}

func SaveTestFile(body []byte) (string, error) {
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	return SaveUploaded(config.Directories.Test, "txt", body)
}

func GetTest(fileName string) ([]string, error) {
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	data, err := LoadFile(path.Join(config.Directories.Test, fileName+".txt"))
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
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	testListData, err := LoadFile(config.TestListFileName)
	if err != nil {
		return nil, err
	}
	log.DebugFmt("List of tests names was read")

	var testListString = utils.BytesToString(testListData)

	result, err := strings.Split(testListString, "\n"), nil
	if err != nil {
		return nil, err
	}

	if len(result) == 1 && (result[0] == "" || result[0] == "\n" || result[0] == "\t" || result[0] == "\r" || result[0] == "\t\r") {
		log.Debug("Return default value")
		return make([]string, 0), nil
	}

	return result, nil
}

func SaveTestList(list []string) error {
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	var data = strings.Join(list, "\n")

	return SaveFile(config.TestListFileName, []byte(data))
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
	var config, err = configuration.GetConfigData()
	log.Check("Configuration error:", err)

	return SaveFile(config.TestListFileName, []byte(""))
}
