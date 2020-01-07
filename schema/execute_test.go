package schema

type testSchema struct{}

/**
TODO: impl the interface below for testSchema

type Query interface {
	Nodes(context.Context, string, QueryArgs) ([]Node, error)
}

type MutationResult struct {
	UIDs []UID
}

type Mutation interface {
	Add(context.Context, []Node) (*MutationResult, error)
	Update(context.Context, Filter, Node) (*MutationResult, error)
	Delete(context.Context, Filter) (*MutationResult, error)
}

type Schema interface {
	Query() Query
	Mutation() Mutation
}
*/

/*
func TestExecute(t *testing.T) {
	tests := []struct {
		q    string
		args Node
	}{
		{
			q: `query QueryNodes($filter: Filter, $order: Order, $first: Int, $offset: Int) {
	nodes(filter: $filter) {
		uid
		name
	}
}`,
		},
	}

	for _, tc := range tests {
		doc, errs := parser.ParseQuery(&ast.Source{
			Name:  "query.graphql",
			Input: tc.q,
		})
		if errs != nil {
			t.Error(errs)
			continue
		}

		got, err := parseQueryArgs(doc.Operations[0].VariableDefinitions, tc.args)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := pretty.Diff(*got, tc.want); len(diff) > 0 {
			t.Errorf("unexpected result\nq: %s\ndiff: %v", tc.q, diff)
		}
	}
}
*/
