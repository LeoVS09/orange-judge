package database

import (
	"context"
	"github.com/machinebox/graphql"
	"orange-judge/log"
)

const databaseUrl = "http://localhost:8765/graphql"

// create a graphql client share across requests
func InitClient() *graphql.Client {
	var client = graphql.NewClient(databaseUrl)

	client.Log = func(s string) {
		log.Log(s)
	}

	return client
}

type Problem struct {
	NodeId string
	Id     string

	InputType struct {
		NodeId string
		Id     string
		name   string
	}

	OutputType struct {
		NodeId string
		Id     string
		name   string
	}

	LimitTime   int
	LimitMemory int

	Tests struct {
		Nodes []struct {
			NodeId string
			Id     string
			Index  int
			Input  string
			Output string
			Public bool
		}
	}
}

type ProblemData struct {
	Problem Problem
}

type Test struct {
	Id     string
	Input  string
	Output string
}

const problemDataQuery = `
		query ($id: UUID!) {
		   problem(id: $id) {
			  nodeId
			  id
			  
			  inputType {
				 nodeId
				 id
				 name
			  }
			  
			  outputType {
				 nodeId
				 id
				 name
			  }
			  
			  limitTime
			  limitMemory
			  
			  tests {
				 nodes {
					nodeId
					id
					index
					input 
					output
					public
				 }
			  }
		   }
		}
	`

// TODO: use another problem structure for return, without tests inside
func GetProblemData(client *graphql.Client, id string) (*Problem, []Test, error) {

	var req = graphql.NewRequest(problemDataQuery)
	req.Var("id", id)

	req.Header.Set("Cache-Control", "no-cache")

	var ctx = context.Background()

	var respData ProblemData
	err := client.Run(ctx, req, &respData)
	if err != nil {
		log.WarningFmt("Error when run graphql request %v", err)
		return nil, []Test{}, err
	}

	log.DebugFmt("Response graphql: %v", respData)

	var problem = respData.Problem
	var input = problem.InputType.name
	var output = problem.OutputType.name

	log.DebugFmt("Input type: %s", input)
	log.DebugFmt("Output type: %s", output)

	var tests = problem.Tests.Nodes
	var testsLength = len(tests)
	log.DebugFmt("Tests count: %v", testsLength)

	var resultTests = make([]Test, testsLength)

	for i, test := range tests {
		resultTests[test.Index] = Test{
			Id:     test.Id,
			Input:  test.Input,
			Output: test.Output,
		}
		log.DebugFmt("Test [%v]: %v", i, test)
	}

	return &problem, resultTests, nil
}
