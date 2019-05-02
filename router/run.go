package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orange-judge/database"
	"orange-judge/executer"
	"orange-judge/fileHandling"
	"orange-judge/log"
)

func pushErrorToResponse(w http.ResponseWriter, r interface{}) {
	var err, ok = r.(error)
	if ok == false {
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}

	log.Error("Error when run", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func pushExecutionErrorToResponse(w http.ResponseWriter, problemId string, r interface{}) {
	executionError, ok := r.(executer.ExecutionError)
	if ok == false {
		panic(r)
	}

	var compilationSuccessful = true
	var failedTest = 0

	switch executionError.Stage {
	case executer.Compilation:
		compilationSuccessful = false
		break
	case executer.Testing:
		failedTest = executionError.TestIndex
		break
	}

	var result = runProgramResponseBody{
		ProblemId:               problemId,
		IsAllTestsSuccessful:    false,
		FailedTest:              failedTest,
		IsCompilationSuccessful: compilationSuccessful,
		StatusCode:              executionError.Status,
	}
	log.DebugFmt("Response to client: %s", result)
	responseBody, err := json.Marshal(result)
	log.Panic("Cannot marshal result for response", err)

	fmt.Fprintf(w, "%s", responseBody)
}

func setResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func parseRunRequestBody(r *http.Request) *runProgramRequestBody {
	var requestBody runProgramRequestBody
	var err = requestBodyParse(r.Body, &requestBody)
	log.Panic("Cannot parse request data", err)

	log.DebugFmt("Code of program:\n%s", requestBody.Code)
	return &requestBody
}

func pushRunResultToResponse(w http.ResponseWriter, problemId string, isTestsSuccessful bool, testNumber int) {
	var result = runProgramResponseBody{
		ProblemId:               problemId,
		IsAllTestsSuccessful:    false,
		FailedTest:              0,
		IsCompilationSuccessful: true,
	}
	defer func() {
		log.DebugFmt("Response to client: %s", result)
		responseBody, err := json.Marshal(result)
		log.Panic("Cannot marshal result for response", err)

		fmt.Fprintf(w, "%s", responseBody)
	}()

	if isTestsSuccessful {
		result.IsAllTestsSuccessful = true
		return
	}

	result.FailedTest = testNumber
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			pushErrorToResponse(w, r)
		}
	}()

	log.DebugFmt("run handler request: %s", r.URL.Path)
	setResponseHeaders(w)

	var requestBody = parseRunRequestBody(r)

	fileName, err := fileHandling.SaveSourceFile([]byte(requestBody.Code))
	log.Panic("Error save uploaded file", err)

	_, tests, err := database.GetProblemData(databaseClient, requestBody.ProblemId)
	log.Panic("Error when get problem data", err)

	defer func() {
		if r := recover(); r != nil {
			pushExecutionErrorToResponse(w, requestBody.ProblemId, r)
		}
	}()

	isTestsSuccessful, testNumber, err := executer.CompileAndTest(fileName, tests)
	log.Panic("Error when compile, run and test program", err)

	pushRunResultToResponse(w, requestBody.ProblemId, isTestsSuccessful, testNumber)
}
