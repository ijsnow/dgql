package schema

import (
	"context"
)

type StringTermFilter struct {
	Equals     *string
	AllOfTerms *string
	AnyOfTerms *string
}

type UID string

type Filter struct {
	UIDs *[]UID
	Name *StringTermFilter
	And  *Filter
	Or   *Filter
	Not  *Filter
}

type Orderable struct {
	By string
}

type Order struct {
	Asc  *Orderable
	Desc *Orderable
	Then *Order
}

// Field is a tuple of [key, value]
type Field [2]interface{}

type Node []Field

type QueryArgs struct {
	Filter *Filter
	Order  *Order
	First  *int
	Offset *int
}

type Query interface {
	Nodes(context.Context, QueryArgs) ([]Node, error)
}

type MutationResult struct {
	UIDs []UID
}

type Mutation interface {
	Add(context.Context, []Node) (*MutationResult, error)
	Update(context.Context, Filter, Node) (*MutationResult, error)
	Delete(context.Context, Filter) (*MutationResult, error)
}

type Schema struct {
	Query    Query
	Mutation Mutation
}
