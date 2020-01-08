package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/vektah/gqlparser/ast"
	"github.com/ijsnow/dgql/dgql/query"
	"github.com/ijsnow/dgql/dgql/schema"
	"google.golang.org/grpc"
)

type Client struct {
	conn       *grpc.ClientConn
	d          *dgo.Dgraph
	s          *schema.Schema
	Predicates map[string]Predicate
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

type Predicate struct {
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
	Predicates []Predicate `json:"schema"`
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

	c.Predicates = map[string]Predicate{}

	for _, pred := range s.Predicates {
		c.Predicates[pred.Name] = pred
	}

	return nil
}

type result struct {
	Nodes []schema.Node `json:"nodes"`
}

func (c *Client) Query(ctx context.Context, doc ast.QueryDocument, args schema.Args) ([]schema.Node, error) {
	qb, err := query.FromSource(&doc, args)
	if err != nil {
		return nil, err
	}

	txn := c.d.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.Query(ctx, string(qb))
	if err != nil {
		return nil, err
	}

	var r result
	err = json.Unmarshal(res.Json, &r)
	if err != nil {
		return nil, err
	}

	return r.Nodes, nil
}

func (c *Client) Add(ctx context.Context, args schema.Args) (*schema.MutationResult, error) {
	txn := c.d.NewTxn()
	defer txn.Discard(ctx)

	pb, err := json.Marshal(args.Nodes)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{
		SetJson: pb,
	}
	req := &api.Request{CommitNow: true, Mutations: []*api.Mutation{mu}}
	res, err := txn.Do(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	uids := []schema.UID{}
	for _, uid := range res.GetUids() {
		uids = append(uids, schema.UID(uid))
	}

	return &schema.MutationResult{
		UIDs: uids,
	}, nil
}

func (c *Client) Update(ctx context.Context, args schema.Args) (*schema.MutationResult, error) {
	vq := query.BuildVarQuery(*args.Filter)

	patch := *args.Patch
	patch["uid"] = "uid(node)"

	pb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{
		SetJson: pb,
	}

	req := &api.Request{
		Query:     vq.String(),
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
	}

	txn := c.d.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	uids := []schema.UID{}
	for _, uid := range res.GetUids() {
		uids = append(uids, schema.UID(uid))
	}

	return &schema.MutationResult{
		UIDs: uids,
	}, nil
}

func (c *Client) Delete(ctx context.Context, args schema.Args) (*schema.MutationResult, error) {
	txn := c.d.NewTxn()
	defer txn.Discard(ctx)

	vq := query.BuildVarQuery(*args.Filter)

	mu := &api.Mutation{}
	dgo.DeleteEdges(mu, "uid(node)", "*")

	req := &api.Request{
		Query:     vq.String(),
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
	}

	res, err := txn.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	uids := []schema.UID{}
	for _, uid := range res.GetUids() {
		uids = append(uids, schema.UID(uid))
	}

	return &schema.MutationResult{
		UIDs: uids,
	}, nil
}
