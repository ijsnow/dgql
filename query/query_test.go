package query

import (
	"testing"

	"gitlab.com/jago-eng/dgql/schema"
)

func toUIDSlcPtr(v []string) *[]schema.UID {
	uids := make([]schema.UID, len(v))
	for idx, uid := range v {
		uids[idx] = schema.UID(uid)
	}

	return &uids
}

func strPtr(s string) *string {
	return &s
}
func intPtr(i int) *int {
	return &i
}

func TestFromSource(t *testing.T) {
	tests := []struct {
		q    string
		want string
		args schema.QueryArgs
	}{
		{
			q: `{
	nodes {
		uid
		name
		age
		friend {
			name
			age
		}
	}
}`,
			want: `{
	nodes(func: has(name, age, friend)) {
		uid
		name
		age
		friend {
			name
			age
		}
	}
}`,
			args: schema.QueryArgs{},
		},
		{
			q: `query QueryNodes($filter: Filter) {
	nodes(filter: $filter) {
		uid
		name
	}
}`,
			want: `{
	nodes(func: uid(0x1, 0x2)) {
		uid
		name
	}
}`,
			args: schema.QueryArgs{
				Filter: &schema.Filter{
					UIDs: toUIDSlcPtr([]string{"0x1", "0x2"}),
				},
			},
		},
		{
			q: `query QueryNodes($filter: Filter) {
	nodes(filter: $filter) {
		uid
		name
	}
}`,
			want: `{
	nodes(func: anyofterms(name, "test")) {
		uid
		name
	}
}`,
			args: schema.QueryArgs{
				Filter: &schema.Filter{
					String: &schema.StringTermFilter{
						Name:       "name",
						AnyOfTerms: strPtr("test"),
					},
				},
			},
		},
		{
			q: `query QueryNodes($filter: Filter) {
	nodes(filter: $filter) {
		uid
		name
	}
}`,
			want: `{
	nodes(func: allofterms(name, "test")) {
		uid
		name
	}
}`,
			args: schema.QueryArgs{
				Filter: &schema.Filter{
					String: &schema.StringTermFilter{
						Name:       "name",
						AllOfTerms: strPtr("test"),
					},
				},
			},
		},
		{
			q: `query QueryNodes($filter: Filter) {
	nodes(filter: $filter) {
		uid
		name
	}
}`,
			want: `{
	nodes(func: has(name), first: 20, offset: 20) {
		uid
		name
	}
}`,
			args: schema.QueryArgs{
				First:  intPtr(20),
				Offset: intPtr(20),
			},
		},
		{
			q: `query QueryNodes($order: Order) {
	nodes(order: $order) {
		uid
		name
	}
}`,
			want: `{
	nodes(func: has(name), orderasc: name) {
		uid
		name
	}
}`,
			args: schema.QueryArgs{
				Order: &schema.Order{
					Asc: strPtr("name"),
				},
			},
		},
		{
			q: `query QueryNodes($order: Order) {
	nodes(order: $order) {
		uid
		name
	}
}`,
			want: `{
	nodes(func: has(name), orderdesc: name) {
		uid
		name
	}
}`,
			args: schema.QueryArgs{
				Order: &schema.Order{
					Desc: strPtr("name"),
				},
			},
		},
	}

	for _, tc := range tests {
		got, err := FromSource(tc.q, tc.args)
		if err != nil {
			t.Error(err)
			continue
		}

		if string(got) != tc.want {
			t.Errorf("unexpected result\nq: `%s`\ngot: `%s`\nwant: `%s`", tc.q, got, tc.want)
		}
	}
}

func TestBuildVarQuery(t *testing.T) {
	tests := []struct {
		want   string
		filter schema.Filter
	}{
		{
			want: `query {
	node as var(func: uid(0x1, 0x2))
}`,
			filter: schema.Filter{
				UIDs: toUIDSlcPtr([]string{"0x1", "0x2"}),
			},
		},
	}

	for _, tc := range tests {
		got := BuildVarQuery(tc.filter)

		if string(got) != tc.want {
			t.Errorf("unexpected result\ngot: `%s`\nwant: `%s`", got, tc.want)
		}
	}
}
