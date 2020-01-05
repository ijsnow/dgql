package dgql

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/jago-eng/dgql/client"
	"gitlab.com/jago-eng/dgql/schema"
)

func TestQuery(t *testing.T) {
	ctx := context.Background()
	client, err := client.NewClient(ctx, client.ClientOptions{"localhost:9080"})
	if err != nil {
		t.Fatal(err)
	}

	q := query{
		c: client,
	}

	if _, err := q.Nodes(ctx, schema.QueryArgs{}); err != nil {
		fmt.Errorf("unexpected error: %v", err)
	}
}

/*
func TestToNode(t *testing.T) {
	predicates := []client.Predicate{
		client.Predicate{
			Name: "UID",
			Type: "uid",
		},
		client.Predicate{
			Name: "String",
			Type: "string",
		},
		client.Predicate{
			Name: "Int",
			Type: "int",
		},
		client.Predicate{
			Name: "Float",
			Type: "float",
		},
		client.Predicate{
			Name: "DateTime",
			Type: "dateTime",
		},
		client.Predicate{
			Name: "Password",
			Type: "password",
		},
		client.Predicate{
			Name:   "ListOfString",
			Type:   "string",
			IsList: true,
		},
		client.Predicate{
			Name:   "Children",
			Type:   "uid",
			IsList: true,
		},
	}
}
*/
