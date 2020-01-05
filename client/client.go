package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"gitlab.com/jago-eng/dgql/schema"
	"google.golang.org/grpc"
)

type Client struct {
	conn       *grpc.ClientConn
	d          *dgo.Dgraph
	s          *schema.Schema
	predicates []predicate
}

func (c *Client) Close() {
	c.conn.Close()
}

type ClientOptions struct {
	Host string
}

func NewClient(ctx context.Context, opt ClientOptions) (*Client, error) {
	conn, err := grpc.Dial(opt.Host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn: conn,
		d:    dgo.NewDgraphClient(api.NewDgraphClient(conn)),
	}

	return c, c.initialize(ctx)
}

type predicate struct {
	Name       string   `json:"predicate"`
	Type       string   `json:"type"`
	IsList     bool     `json:"list"`
	Tokenizers []string `json:"tokenizers"`
	Index      bool     `json:"index"`
}

var omitPredicateNames = map[string]struct{}{
	"dgraph.type": struct{}{},
}

type sch struct {
	Predicates []predicate `json:"schema"`
}

func (c *Client) initialize(ctx context.Context) error {
	res, err := c.d.NewTxn().Query(ctx, "schema {}")
	if err != nil {
		return err
	}

	var s sch
	err = json.Unmarshal(res.Json, &s)
	if err != nil {
		return err
	}

	c.predicates = s.Predicates
	fmt.Printf("%+v\n", s.Predicates)

	// preds := toPredicates(s.Predicates)

	// c.s = schema.BuildSchema(preds)

	return nil
}

func (c *Client) Query(ctx context.Context, query string, target interface{}) error {
	txn := c.d.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.Query(ctx, query)
	if err != nil {
		return err
	}

	return json.Unmarshal(res.Json, target)
}
