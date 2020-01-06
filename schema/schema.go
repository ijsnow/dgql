package schema

import (
	"context"
)

type StringTermFilter struct {
	Name       string
	AllOfTerms *string
	AnyOfTerms *string
}

type UID string

func (u UID) String() string {
	return string(u)
}

func UIDsToStrings(uids []UID) []string {
	out := make([]string, len(uids))
	for idx, uid := range uids {
		out[idx] = uid.String()
	}
	return out
}

type Filter struct {
	UIDs   *[]UID
	String *StringTermFilter
	And    *Filter
	// Or     *Filter
	// Not    *Filter
}

type Order struct {
	Asc  *string
	Desc *string
	// Then *Order
}

type Node map[string]interface{}

func (n Node) SetField(field string, value interface{}) {
	n[field] = value
}

type QueryArgs struct {
	Filter *Filter
	Order  *Order
	First  *int
	Offset *int
}

type Query interface {
	Nodes(context.Context, string, QueryArgs) ([]Node, error)
}

type MutationResult struct {
	UIDs []UID
}

type Mutation interface {
	Add(context.Context, []Node) (*MutationResult, error)
	Update(context.Context, Filter, Node) (*MutationResult, error)
	Delete(context.Context, Filter) (*MutationResult, error)
}

type Schema interface {
	Query() Query
	Mutation() Mutation
}
