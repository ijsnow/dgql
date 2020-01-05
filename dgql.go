package dgql

import (
	"context"
	"fmt"

	"gitlab.com/jago-eng/dgql/client"
	"gitlab.com/jago-eng/dgql/schema"
)

type Schema schema.Schema

type query struct {
	c *client.Client
}

type result struct {
	Nodes []map[string]interface{}
}

func (q query) Nodes(ctx context.Context, args schema.QueryArgs) ([]schema.Node, error) {
	var r result
	err := q.c.Query(ctx, "{ nodes (func: has(name)) { uid name friend { uid name } } }", &r)
	if err != nil {
		return nil, err
	}

	for _, node := range r.Nodes {
		for field := range node {
			fmt.Printf("%+v %+v %+v\n", field, q.c.Predicates[field], node[field])

		}

		if v, ok := node["friend"]; ok {
			b_temp := v.([]interface{})
			fmt.Printf("%+v\n", b_temp)
		}
	}

	// https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json

	panic("not impl")
}
