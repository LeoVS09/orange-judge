package main

import (
	"errors"
	"flag"
	"math/rand"
	"orange-judge/configuration"
	"orange-judge/executer"
	"orange-judge/log"
	"orange-judge/router"
	"strings"
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

	configStore := configuration.Read(*configName, configReadDelay)
	config, err := configuration.ToConfigData(configStore)
	log.Check("Configuration error:", err)

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

func testCompiler() (bool, error) {
	log.Debug("Start test compiler environment...")
	const input = "1 2 3"
	const output = "3 2 1"
	var out, err = executer.TestRunFromSource(input, "test")
	if err != nil {
		return false, err
	}

	var result = removeUnnecessarySymbols(out.String())
	if result != output {
		return false, nil
	}
	return true, nil
}

func removeUnnecessarySymbols(data string) string {
	var removeSymbols = []string{"\n", "\t", "\r"}
	for _, symbol := range removeSymbols {
		data = strings.Replace(data, symbol, "", -1)
	}
	return data
}
