package main

import (
	"errors"
	"flag"
	"math/rand"
	"orange-judge/configuration"
	"orange-judge/executer"
	"orange-judge/fileHandling"
	"orange-judge/log"
	"orange-judge/router"
	"orange-judge/utils"
	"os"
	"time"
)

const defaultConfigName = "config.json"
const configReadDelay = time.Minute

func main() {
	isDebug := flag.Bool("d", false, "Print more debug output.")
	configName := flag.String("c", defaultConfigName, "Configuration file name")
	isNeedTestCompiler := flag.Bool("tc", true, "Test cannot be use compiler in environment")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if *isNeedTestCompiler == true {
		configureLogging(configuration.Testing)
	}

	configuration.Read(*configName, configReadDelay)
	config, err := configuration.GetConfigData()
	log.Check("Configuration error:", err)

	err = createOrClearWorkFolders(config)
	log.Check("Cannot clear work environment", err)

	if *isNeedTestCompiler == true {
		isTestSuccessful, err := testCompiler()
		log.Check("Error when compile and run test program", err)
		if isTestSuccessful {
			log.Log("compiler test: ok")
		} else {
			log.Log("compiler test: fatal")
			panic(errors.New("Test program return not expected result"))
		}
	}

	if *isDebug == true {
		configureLogging(configuration.Development)
	} else {
		configureLogging(configuration.Production)
	}

	err = router.ListenAndServe(config.Port)
	log.Check("Error serving", err)
}

func configureLogging(env configuration.Environment) {

	if env == configuration.Testing {
		log.VerboseAll()
		log.Log("Test mode enabled")
		return
	}

	if env == configuration.Development {
		log.VerboseAll()
		log.Log("Debug mode enabled")
		return
	}

	if env == configuration.Production {
		log.VerboseProduction()
		return
	}

	log.VerboseOnlyErrors()
}

func createOrClearWorkFolders(config *configuration.ConfigFile) error {
	var names = [...]string{config.Directories.Compiled, config.Directories.Uploaded, config.Directories.Test}
	for _, name := range names {
		var err = fileHandling.CreateIfNotExistFolder(name)
		if err != nil {
			log.DebugFmt("Cannot create or clear folder: %s", name)
			return err
		}
	}

	if fileHandling.IsExist(config.TestListFileName) {
		//var err = fileHandling.ClearTestList()
		//if err != nil {
		//	log.DebugFmt("Cannot clear tests list file %s", config.TestListFileName)
		//	return err
		//}
		return nil
	}

	var f, err = os.Create(config.TestListFileName)
	if err != nil {
		log.DebugFmt("Cannot create tests list file %s", config.TestListFileName)
		return err
	}

	return f.Close()
}

func testCompiler() (bool, error) {
	log.Debug("Start test compiler environment...")
	const input = "2 3"
	const output = "6"
	const fileName = "test.cpp"

	var testSourceFile, err = fileHandling.LoadFile(fileName)
	if err != nil {
		log.DebugFmt("Cannot load test source file %s", fileName)
		return false, err
	}

	newFileName, err := fileHandling.SaveSourceFile(testSourceFile)
	if err != nil {
		log.DebugFmt("Cannot save test source file %s", fileName)
		return false, err
	}

	out, status, err := executer.RunFromSourceWithOAR(newFileName, input)
	if err != nil {
		return false, err
	}

	if status != 0 {
		log.LogFmt("Not successful exit code: %v", status)
	}

	var result = utils.RemoveUnnecessarySymbols(out.String())
	if result != output {
		return false, nil
	}
	return true, nil
}
