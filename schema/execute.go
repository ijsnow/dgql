package schema

import (
	"context"
	"errors"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type ExecutionArgs struct {
	Source string
	Args   Node
}

type ResultData struct {
	Nodes []Node `json:"nodes"`
}

type ExecutionResult struct {
	Data   interface{} `json:"data,omitempty`
	Errors []error     `json:"errors,omitempty"`
}

func Execute(ctx context.Context, schema Schema, args ExecutionArgs) (*ExecutionResult, error) {
	doc, err := parser.ParseQuery(&ast.Source{
		Name:  "query.graphql",
		Input: args.Source,
	})
	if err != nil {
		return nil, err
	}

	for _, op := range doc.Operations[:1] { // Only support one operation for now
		qa, err := parseArgs(op.VariableDefinitions, args.Args)
		if err != nil {
			return nil, err
		}

		res := ExecutionResult{}

		switch op.Operation {
		case "query":
			query := schema.Query()

			nodes, err := query.Nodes(ctx, *doc, *qa)
			if err != nil {
				res.Errors = []error{err}
			} else {
				res.Data = ResultData{nodes}
			}

		case "mutation":
			sel, ok := op.SelectionSet[0].(*ast.Field)
			if !ok {
				return nil, errors.New("could not cast mutation selection")
			}
			mutation := schema.Mutation()

			switch sel.Name {
			case "add":
				mres, err := mutation.Add(ctx, *qa)
				if err != nil {
					res.Errors = []error{err}
				} else {
					res.Data = mres
				}
			case "update":
				mres, err := mutation.Update(ctx, *qa)
				if err != nil {
					res.Errors = []error{err}
				} else {
					res.Data = mres
				}
			case "delete":
				mres, err := mutation.Delete(ctx, *qa)
				if err != nil {
					res.Errors = []error{err}
				} else {
					res.Data = mres
				}
			default:
				panic("should not reach 1")
			}
		default:
			panic("should not reach 2")
		}

		return &res, nil
	}

	return nil, nil
}
