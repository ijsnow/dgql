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

func (q query) Nodes(ctx context.Context, args schema.QueryArgs) ([]schema.Node, error) {
	fmt.Printf("%+v", fields)

	var result map[string]interface{}
	err := q.c.Query(ctx, "", &result)
	if err != nil {
		return nil, err
	}

	fmt.Printf("result => %+v", result)

	// https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json

	panic("not impl")
}
