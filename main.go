package main

import (
	"flag"
	"math/rand"
	"orange-judge/configuration"
	"orange-judge/executer"
	"orange-judge/log"
	"orange-judge/router"
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
		err = testCompiler()
		log.Check("Error when compile and run test program", err)
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

func testCompiler() error {
	log.Debug("Start test compiler environment...")
	var _, err = executer.TestRunFromSource("test")
	return err
}
