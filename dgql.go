package dgql

import (
	"context"

	"github.com/vektah/gqlparser/ast"
	"gitlab.com/jago-eng/dgql/client"
	"gitlab.com/jago-eng/dgql/schema"
)

type Schema schema.Schema

type query struct {
	c *client.Client
}

type result struct {
	Nodes []schema.Node
}

func (q query) Nodes(ctx context.Context, query *ast.QueryDocument, args schema.QueryArgs) ([]schema.Node, error) {
	var r result
	err := q.c.Query(ctx, "{ nodes (func: has(name)) { uid name friend { uid name } } }", &r)
	if err != nil {
		return nil, err
	}

	// https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json

	return r.Nodes, nil
}
