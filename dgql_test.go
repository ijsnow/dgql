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
