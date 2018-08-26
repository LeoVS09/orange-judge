package main

import (
	"flag"
	"fmt"
	"net/http"
	"orange-judge/src/config"
	"orange-judge/src/log"
	"time"
)

const defaultConfigName = "config.json"
const configReadDelay = time.Minute

func main() {
	verbosePtr := flag.Bool("verbose", false, "Print more debug output.")
	configName := flag.String("config", defaultConfigName, "Configuration file name")
	flag.Parse()

	if *verbosePtr == true {
		log.VerboseAll()
	} else {
		log.VerboseProduction()
	}

	log.UseSync()       // bug with async logging
	log.SetNoColoring() // bug with coloring on windows

	configStore := config.Read(*configName, configReadDelay)
	configData, err := config.ToConfigData(configStore)
	log.Check("Configuration error:", err)

	http.HandleFunc("/", handler)

	log.LogFmt("Serving at localhost:%v...", configData.Port)
	log.Check("Error serving",
		http.ListenAndServe(fmt.Sprintf(":%v", configData.Port), nil),
	)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi %s sucks!", r.URL.Path[1:])
}
