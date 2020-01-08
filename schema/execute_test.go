package schema

import (
	"context"
	"errors"
	"testing"

	"github.com/kr/pretty"
	"github.com/vektah/gqlparser/ast"
)

var (
	errNodes  = errors.New("nodes")
	errAdd    = errors.New("add")
	errUpdate = errors.New("update")
	errDelete = errors.New("delete")
)

type testSchema struct {
	called string
	args   *Args
}

func (s *testSchema) Query() Query {
	return s
}

func (s *testSchema) Mutation() Mutation {
	return s
}

func (s *testSchema) Nodes(ctx context.Context, doc ast.QueryDocument, args Args) ([]Node, error) {
	s.called = "nodes"
	s.args = &args
	return nil, nil
}

func (s *testSchema) Add(ctx context.Context, args Args) (*MutationResult, error) {
	s.called = "add"
	s.args = &args
	return nil, nil
}

func (s *testSchema) Update(ctx context.Context, args Args) (*MutationResult, error) {
	s.called = "update"
	s.args = &args
	return nil, nil
}

func (s *testSchema) Delete(ctx context.Context, args Args) (*MutationResult, error) {
	s.called = "delete"
	s.args = &args
	return nil, nil
}

func TestExecute(t *testing.T) {
	tests := []struct {
		q          string
		args       Node
		wantCalled string
		wantArgs   *Args
	}{
		{
			q: `query QueryNodes($filter: Filter, $order: Order, $first: Int, $offset: Int) {
	nodes(filter: $filter)
}`,
			args: Node{
				"filter": Node{
					"uids": []string{"0x1", "0x2"},
				},
			},
			wantCalled: "nodes",
			wantArgs: &Args{
				Filter: &Filter{
					UIDs: toUIDSlcPtr([]string{"0x1", "0x2"}),
				},
			},
		},
		{
			q: `mutation AddNodes($nodes: [Node!]!) {
	add(nodes: $nodes)
}`,
			args: Node{
				"nodes": []Node{
					Node{
						"uid": "0x1",
					},
					Node{
						"uid": "0x2",
					},
				},
			},
			wantCalled: "add",
			wantArgs: &Args{
				Nodes: []Node{
					Node{
						"uid": "0x1",
					},
					Node{
						"uid": "0x2",
					},
				},
			},
		},
		{
			q: `mutation UpdateNodes($filter: Filter, $patch: Node) {
	update(nodes: $nodes)
}`,
			args: Node{
				"filter": Node{
					"uids": []string{"0x1", "0x2"},
				},
				"patch": Node{
					"name": "a",
				},
			},
			wantCalled: "update",
			wantArgs: &Args{
				Filter: &Filter{
					UIDs: toUIDSlcPtr([]string{"0x1", "0x2"}),
				},
				Patch: &Node{
					"name": "a",
				},
			},
		},
		{
			q: `mutation DeleteNodes($filter: Filter) {
	delete(nodes: $nodes)
}`,
			args: Node{
				"filter": Node{
					"uids": []string{"0x1", "0x2"},
				},
			},
			wantCalled: "delete",
			wantArgs: &Args{
				Filter: &Filter{
					UIDs: toUIDSlcPtr([]string{"0x1", "0x2"}),
				},
			},
		},
	}

	ctx := context.Background()

	for _, tc := range tests {
		ts := &testSchema{}

		_, err := Execute(ctx, ts, ExecutionArgs{
			Source: tc.q,
			Args:   tc.args,
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if ts.called == "" {
			t.Fatalf("no method was called: \n%s\n", tc.q)
		}

		if ts.called != tc.wantCalled {
			t.Errorf("correct method did not get called\ngot: %s\nwant: %s", ts.called, tc.wantCalled)
		}

		if diff := pretty.Diff(ts.args, tc.wantArgs); len(diff) > 0 {
			t.Errorf("incorrect query args\ndiff: %v", diff)
		}
	}
}
