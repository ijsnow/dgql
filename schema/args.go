package schema

import (
	"errors"
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

var nodeArgs = map[string]struct{}{
	"filter": struct{}{},
	"order":  struct{}{},
}

func parseArgs(defs ast.VariableDefinitionList, args Node) (*QueryArgs, error) {
	var err error
	qa := QueryArgs{}

	for _, def := range defs {
		arg, ok := args[def.Variable]
		if !ok {
			continue
		}

		switch def.Variable {
		case "filter":
			node, ok := arg.(Node)
			if !ok {
				return nil, fmt.Errorf("could not cast argument to node: %s", def.Variable)
			}
			qa.Filter = &Filter{}
			err = parseFilterNode(node, qa.Filter)
		case "order":
			node, ok := arg.(Node)
			if !ok {
				return nil, fmt.Errorf("could not cast argument to node: %s", def.Variable)
			}
			qa.Order = &Order{}
			err = parseOrderNode(node, qa.Order)
		case "first":
			first, ok := arg.(int)
			if !ok {
				return nil, errors.New("could not cast first")
			}
			qa.First = &first
		case "offset":
			offset, ok := arg.(int)
			if !ok {
				return nil, errors.New("could not cast offset")
			}
			qa.Offset = &offset
		default:
			return nil, fmt.Errorf("illegal argument name: %s", def.Variable)
		}
		if err != nil {
			return nil, err
		}
	}

	return &qa, nil
}

func parseFilterNode(node Node, target *Filter) error {
	for f, v := range node {
		switch f {
		case "uids":
			strs, ok := v.([]string)
			if !ok {
				return errors.New("could not cast uids")
			}
			uids := StringsToUIDs(strs)
			target.UIDs = &uids
		case "term":
			str, ok := v.(Node)
			if !ok {
				return errors.New("could not cast string term filter")
			}
			target.Term = &TermFilter{}
			parseTermFilterNode(str, target.Term)
		default:
			return fmt.Errorf("illegal filter field: %s", f)
		}
	}

	return nil
}

func parseTermFilterNode(node Node, target *TermFilter) error {
	var ok bool

	for f, v := range node {
		switch f {
		case "name":
			var name string
			name, ok = v.(string)
			if !ok {
				return errors.New("could not cast name")
			}
			target.Name = name
		case "all":
			t, ok := v.(string)
			if !ok {
				return errors.New("could not cast name")
			}
			target.All = &t
		case "any":
			t, ok := v.(string)
			if !ok {
				return errors.New("could not cast name")
			}
			target.Any = &t
		default:
			return fmt.Errorf("illegal term filter field: %s", f)
		}
	}

	return nil
}

func parseOrderNode(node Node, target *Order) error {
	for f, v := range node {
		switch f {
		case "asc":
			asc, ok := v.(string)
			if !ok {
				return errors.New("could not cast asc")
			}
			target.Asc = &asc
		case "desc":
			desc, ok := v.(string)
			if !ok {
				return errors.New("could not cast desc")
			}
			target.Desc = &desc
		default:
			return fmt.Errorf("illegal order field: %s", f)
		}
	}

	return nil
}
