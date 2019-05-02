package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orange-judge/database"
	"orange-judge/fileHandling"
	"orange-judge/log"
)

var databaseClient = database.InitClient()

func SetHandlers() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)

	http.HandleFunc("/run", runHandler)

	http.HandleFunc("/test/upload", testUploadHandler)
}

func ListenAndServe(port int) error {
	SetHandlers()

	log.LogFmt("Serving at localhost:%v...", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var body, err = fileHandling.LoadPage("index")
	log.Panic("Cannot load index page", err)
	w.Write(body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	var body = r.FormValue("body")
	var err = fileHandling.SaveFile("test-uploaded.cpp", []byte(body))
	log.Panic("Error save uploaded file", err)
	http.Redirect(w, r, "/", http.StatusFound)
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

type runProgramRequestBody struct {
	ProblemId string `json:"problemId"`
	Code      string `json:"code"`
}

type runProgramResponseBody struct {
	ProblemId               string `json:"problemId"`
	IsAllTestsSuccessful    bool   `json:"isAllTestsSuccessful"`
	FailedTest              int    `json:"failedTest"`
	IsCompilationSuccessful bool   `json:"isCompilationSuccessful"`
	StatusCode              int    `json:"status"`
}

func requestBodyParse(requestBody io.ReadCloser, v interface{}) error {
	var buffer = new(bytes.Buffer)
	buffer.ReadFrom(requestBody)
	var body = buffer.String()
	log.DebugFmt("Request data: %s", body)

	return json.Unmarshal([]byte(body), v)
}

type testUploadRequestBody struct {
	Text string `json:"text"`
}

type testUploadResponseBody struct {
	IsSuccessfulAdded bool `json:"isSuccessfulAdded"`
}

func testUploadHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			var err, ok = r.(error)
			if ok == false {
				http.Error(w, "Unexpected error", http.StatusInternalServerError)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	log.Debug("New test upload start")

	var result = testUploadResponseBody{
		IsSuccessfulAdded: false,
	}
	defer func() {
		log.DebugFmt("Response to client: %s", result)
		responseBody, err := json.Marshal(result)
		log.Panic("Cannot marshal result for response", err)

		fmt.Fprintf(w, "%s", responseBody)
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestBody testUploadRequestBody
	var err = requestBodyParse(r.Body, &requestBody)
	log.Panic("Cannot parse request data", err)

	name, err := fileHandling.SaveTestFile([]byte(requestBody.Text))
	log.Panic("Error save test file", err)

	err = fileHandling.AddTestToList(name)
	log.Panic("Cannot add test to list", err)

	result.IsSuccessfulAdded = true
}
