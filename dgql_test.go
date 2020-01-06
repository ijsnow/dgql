package dgql

import (
	"context"
	"testing"

	"gitlab.com/jago-eng/dgql/client"
	"gitlab.com/jago-eng/dgql/schema"
)

func TestAddAndQuery(t *testing.T) {
	ctx := context.Background()
	client, err := client.NewClient(ctx, client.ClientOptions{"localhost:9080"})
	if err != nil {
		t.Fatal(err)
	}

	s := Schema{
		c: client,
	}

	nodes := []schema.Node{
		schema.Node{
			"name": "Isaac Snow",
			"age":  25,
			"friend": []schema.Node{
				schema.Node{
					"name": "Luci",
					"age":  26,
				},
			},
		},
	}

	mres, err := s.Add(ctx, nodes)
	if err != nil {
		t.Fatal(err)
	}

	if len(mres.UIDs) != 2 {
		t.Fatalf("expected %d uids, got %v", 2, len(mres.UIDs))
	}

	qs := `query QueryNodes($filter: Filter) {
		nodes(filter: $filter) {
			uid
			name
			age
			friend {
				uid
				name
				age
			}
		}
	}`

	qres, err := s.Nodes(ctx, qs, schema.QueryArgs{
		Filter: &schema.Filter{
			UIDs: &mres.UIDs,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(qres) != 2 {
		t.Fatalf("expected %d nodes, got %v", 2, len(mres.UIDs))
	}

	ures, err := s.Update(ctx, schema.Filter{UIDs: &mres.UIDs}, schema.Node{
		"name": "updates",
	})
	if err != nil {
		t.Fatal(err)
	}

	qs = `query QueryNodes($filter: Filter) {
		nodes(filter: $filter) {
			uid
			name
		}
	}`

	qres, err = s.Nodes(ctx, qs, schema.QueryArgs{
		Filter: &schema.Filter{
			UIDs: &ures.UIDs,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range qres {
		if node["name"] != "test" {
			t.Errorf("updated did not work; expected node.name to be `test`, got `%v`", node["name"])
		}
	}

	dres, err := s.Delete(ctx, schema.Filter{UIDs: &mres.UIDs})
	if err != nil {
		t.Fatal(err)
	}

	qs = `query QueryNodes($filter: Filter) {
		nodes(filter: $filter) {
			uid
			name
		}
	}`

	qres, err = s.Nodes(ctx, qs, schema.QueryArgs{
		Filter: &schema.Filter{
			UIDs: &dres.UIDs,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range qres {
		if len(node) > 1 {
			t.Errorf("expected to get no fields, got %+v", node)
		}
	}
}