package schema

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type ExecutionArgs struct {
	Source string
	Args   Node
}

type ExecutionResult struct {
	Data struct {
		Nodes []Node `json:"nodes"`
	} `json:"data"`
	Errors []error `json:"errors,omitempty"`
}

func Execute(schema Schema, args ExecutionArgs) (*ExecutionResult, error) {
	doc, err := parser.ParseQuery(&ast.Source{
		Name:  "query.graphql",
		Input: args.Source,
	})
	if err != nil {
		return nil, err
	}

	for _, op := range doc.Operations[:1] { // Only support one operation for now
		qa, err := parseQueryArgs(op.VariableDefinitions, args.Args)
		if err != nil {
			return nil, err
		}

		fmt.Printf("%+v %+v\n", args.Source, qa)
	}

	return nil, nil
}
