package main

import (
	"flag"
	"fmt"
	"net/http"
	"orange-judge/config"
	"orange-judge/executer"
	"orange-judge/log"
	"time"
)

const defaultConfigName = "config.json"
const configReadDelay = time.Minute

func main() {
	debug := flag.Bool("d", false, "Print more debug output.")
	configName := flag.String("c", defaultConfigName, "Configuration file name")
	flag.Parse()

	configureLogging(*debug)

	configStore := config.Read(*configName, configReadDelay)
	configData, err := config.ToConfigData(configStore)
	log.Check("Configuration error:", err)

	http.HandleFunc("/", handler)

	log.LogFmt("Serving at localhost:%v...", configData.Port)
	//log.Check("Error serving",
	//	http.ListenAndServe(fmt.Sprintf(":%v", configData.Port), nil),
	//)

	executer.Run()
}

func configureLogging(debug bool) {
	if debug == true {
		log.VerboseAll()
	} else {
		log.VerboseProduction()
	}

	log.UseSync()       // bug with async logging
	log.SetNoColoring() // bug with coloring on windows
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi %s sucks!", r.URL.Path[1:])
}
