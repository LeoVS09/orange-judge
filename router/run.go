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

func runHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			var err, ok = r.(error)
			if ok == false {
				http.Error(w, "Unexpected error", http.StatusInternalServerError)
				return
			}

			log.Error("Error when run", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	log.DebugFmt("run handler request: %s", r.URL.Path)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var requestBody runProgramRequestBody
	var err = requestBodyParse(r.Body, &requestBody)
	log.Panic("Cannot parse request data", err)

	log.DebugFmt("Code of program:\n%s", requestBody.Code)

	fileName, err := fileHandling.SaveSourceFile([]byte(requestBody.Code))
	log.Panic("Error save uploaded file", err)

	err, _, _, tests := database.GetProblemData(databaseClient, requestBody.ProblemId)
	log.Panic("Error when get problem data", err)

	defer func() {
		r := recover()
		if r == nil {
			return
		}

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
			ProblemId:               requestBody.ProblemId,
			IsAllTestsSuccessful:    false,
			FailedTest:              failedTest,
			IsCompilationSuccessful: compilationSuccessful,
		}
		log.DebugFmt("Response to client: %s", result)
		responseBody, err := json.Marshal(result)
		log.Panic("Cannot marshal result for response", err)

		fmt.Fprintf(w, "%s", responseBody)
	}()

	isTestsSuccessful, testNumber, err := executer.CompileAndTest(fileName, tests)
	log.Panic("Error when compile, run and test program", err)

	var result = runProgramResponseBody{
		ProblemId:               requestBody.ProblemId,
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
