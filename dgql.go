package dgql

import (
	"context"

	"github.com/vektah/gqlparser/ast"
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

func (s *Schema) Nodes(ctx context.Context, doc ast.QueryDocument, args schema.Args) ([]schema.Node, error) {
	return s.c.Query(ctx, doc, args)
}

func (s *Schema) Mutation() schema.Mutation {
	return s
}

func (s *Schema) Add(ctx context.Context, args schema.Args) (*schema.MutationResult, error) {
	return s.c.Add(ctx, args)
}

func (s *Schema) Update(ctx context.Context, args schema.Args) (*schema.MutationResult, error) {
	return s.c.Update(ctx, args)
}

func (s *Schema) Delete(ctx context.Context, args schema.Args) (*schema.MutationResult, error) {
	return s.c.Delete(ctx, args)
}
