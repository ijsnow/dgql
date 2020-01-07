package dgql

import (
	"context"

	"gitlab.com/jago-eng/dgql/client"
	"gitlab.com/jago-eng/dgql/schema"
)

type Schema struct {
	c *client.Client
}

func NewSchema(ctx context.Context, opt client.ClientOptions) (*Schema, error) {
	c, err := client.NewClient(ctx, opt)
	if err != nil {
		return nil, err
	}

	return &Schema{c}, err
}

// https://stackoverflow.com/questions/42152750/golang-is-there-an-easy-way-to-unmarshal-arbitrary-complex-json

func (s *Schema) Query() schema.Query {
	return s
}

func (s *Schema) Nodes(ctx context.Context, qs string, args schema.QueryArgs) ([]schema.Node, error) {
	return s.c.Query(ctx, qs, args)
}

func (s *Schema) Mutation() schema.Mutation {
	return s
}

func (s *Schema) Add(ctx context.Context, nodes []schema.Node) (*schema.MutationResult, error) {
	return s.c.Add(ctx, nodes)
}

func (s *Schema) Update(ctx context.Context, filter schema.Filter, patch schema.Node) (*schema.MutationResult, error) {
	return s.c.Update(ctx, filter, patch)
}

func (s *Schema) Delete(ctx context.Context, filter schema.Filter) (*schema.MutationResult, error) {
	return s.c.Delete(ctx, filter)
}
