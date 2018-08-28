package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"orange-judge/config"
	"orange-judge/executer"
	"orange-judge/log"
	"time"
)

const defaultConfigName = "config.json"
const configReadDelay = time.Minute

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	isDebug := flag.Bool("d", false, "Print more debug output.")
	configName := flag.String("c", defaultConfigName, "Configuration file name")
	isNeedTestCompiler := flag.Bool("tc", true, "Test cannot be use compiler in environment")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	var resultConfigurator = configureLogging(*isDebug, *isNeedTestCompiler)

	configStore := config.Read(*configName, configReadDelay)
	configData, err := config.ToConfigData(configStore)
	log.Check("Configuration error:", err)

	testCompiler(*isNeedTestCompiler)
	resultConfigurator()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	var result = resultRunProgram{}
	http.HandleFunc("/run", runHandler(&result))
	http.HandleFunc("/run/result/", resultHandler(&result))

	log.LogFmt("Serving at localhost:%v...", configData.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", configData.Port), nil)
	log.Check("Error serving", err)
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func configureLogging(debug bool, test bool) func() {
	log.UseSync()       // bug with async logging
	log.SetNoColoring() // bug with coloring on windows

	if test == true {
		log.VerboseAll()
	}

	return func() {
		if debug == true {
			log.VerboseAll()
		} else {
			log.VerboseProduction()
		}
	}
}

func testCompiler(test bool) {
	if !test {
		return
	}
	log.Debug("Start test compiler environment...")
	var _, err = executer.RunFromSource("test")
	log.Check("Error when compile and run test program", err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var body, err = loadPage("index")
	log.Panic("Cannot load index page", err)
	w.Write(body)
}

func loadPage(name string) ([]byte, error) {
	return loadFile(name + ".html")
}

func loadFile(name string) ([]byte, error) {
	body, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	var body = r.FormValue("body")
	var err = saveFile("test-uploaded.cpp", []byte(body))
	log.Panic("Error save uploaded file", err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func saveFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0600)
}

type resultRunProgram struct {
	id  string
	out bytes.Buffer
}

func resultHandler(result *resultRunProgram) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.DebugFmt("result handler response: %s", result.out.String())
		fmt.Fprintf(w, "%s", result.out.String())
	}
}

func runHandler(result *resultRunProgram) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.DebugFmt("run handler request: %s", r.URL.Path)

		const fileName = "test-uploaded"
		var body = r.FormValue("body")

		var err = saveFile(fileName+".cpp", []byte(body))
		log.Panic("Error save uploaded file", err)

		resultOut, err := executer.RunFromSource(fileName)
		log.Panic("Error when compile and run test program", err)

		var resultId = RandStringRunes(50)
		// TODO: use array with managed content for hold results
		result.id = resultId
		result.out = *resultOut

		http.Redirect(w, r, "/run/result/"+resultId, http.StatusFound)
	}
}
