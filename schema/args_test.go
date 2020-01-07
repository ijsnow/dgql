package schema

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

func toUIDSlcPtr(v []string) *[]UID {
	uids := make([]UID, len(v))
	for idx, uid := range v {
		uids[idx] = UID(uid)
	}

	return &uids
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		q    string
		args Node
		want QueryArgs
	}{
		{
			q: `query QueryNodes($filter: Filter, $order: Order, $first: Int, $offset: Int) {
	nodes(filter: $filter) {
		uid
		name
	}
}`,
			args: Node{
				"filter": Node{
					"uids": []string{"0x1", "0x2"},
					"term": Node{
						"name": "a",
						"all":  "test",
						"any":  "test",
					},
				},
				"order": Node{
					"asc":  "a",
					"desc": "b",
				},
				"first":  20,
				"offset": 20,
			},
			want: QueryArgs{
				Filter: &Filter{
					UIDs: toUIDSlcPtr([]string{"0x1", "0x2"}),
					Term: &TermFilter{
						Name: "a",
						// Setting both of these should be an error, but we can just have dgraph raise this.
						All: strPtr("test"),
						Any: strPtr("test"),
					},
				},
				Order: &Order{
					// Setting both of these should be an error, but we can just have dgraph raise this.
					Asc:  strPtr("a"),
					Desc: strPtr("b"),
				},
				First:  intPtr(20),
				Offset: intPtr(20),
			},
		},
	}

	for _, tc := range tests {
		doc, errs := parser.ParseQuery(&ast.Source{
			Name:  "query.graphql",
			Input: tc.q,
		})
		if errs != nil {
			t.Error(errs)
			continue
		}

		got, err := parseQueryArgs(doc.Operations[0].VariableDefinitions, tc.args)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := pretty.Diff(*got, tc.want); len(diff) > 0 {
			t.Errorf("unexpected result\nq: %s\ndiff: %v", tc.q, diff)
		}
	}
}
