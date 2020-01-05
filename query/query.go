package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"gitlab.com/jago-eng/dgql/schema"
)

type Query string
type QueryVars map[string]interface{}

func FromSource(in string, args schema.QueryArgs) (Query, error) {
	source := ast.Source{
		Name:  "query.graphql",
		Input: in,
	}

	doc, err := parser.ParseQuery(&source)
	if err != nil {
		return "", err
	}

	q, errr := write(doc, args)
	return Query(q), errr
}

func write(doc *ast.QueryDocument, args schema.QueryArgs) (string, error) {
	var sb strings.Builder
	qb := queryBuilder{&sb, 0, args}

	qb.push()

	op := doc.Operations[0]

	err := qb.writeOuterMostSelectionSet(op.SelectionSet)
	if err != nil {
		return "", err
	}

	qb.pop()

	return sb.String(), nil
}

var specialFields = map[string]struct{}{
	"uid": struct{}{},
}

type queryBuilder struct {
	sb     *strings.Builder
	levels int
	args   schema.QueryArgs
}

func (qb *queryBuilder) write(s string) { qb.sb.WriteString(s) }
func (qb *queryBuilder) writef(s string, vals ...interface{}) {
	qb.sb.WriteString(fmt.Sprintf(s, vals...))
}
func (qb *queryBuilder) writeR(s rune)  { qb.sb.WriteRune(s) }
func (qb *queryBuilder) writeInt(i int) { qb.sb.WriteString(strconv.Itoa(i)) }
func (qb *queryBuilder) newLine()       { qb.writeR('\n') }
func (qb *queryBuilder) leftPad() {
	for i := 0; i < qb.levels; i++ {
		qb.writeR('\t')
	}
}

func (qb *queryBuilder) push() {
	qb.levels += 1
	qb.writeR('{')
	qb.newLine()
}

func (qb *queryBuilder) pop() {
	qb.levels -= 1
	qb.newLine()
	qb.leftPad()
	qb.writeR('}')
}

func (qb *queryBuilder) writeln(s string) {
	qb.leftPad()
	qb.sb.WriteString(s)
	qb.newLine()
}

func (qb *queryBuilder) writeOuterMostSelectionSet(sset ast.SelectionSet) error {
	for _, selection := range sset {
		if field, ok := selection.(*ast.Field); ok {
			qb.writeField(field)
			err := qb.writeQueryArgs(field.SelectionSet)
			if err != nil {
				return err
			}

			qb.writeSelectionSet(field.SelectionSet)
			if err != nil {
				return err
			}

			continue
		}

		return errors.New("unable to convert selection to field")
	}

	return nil
}

func (qb *queryBuilder) writeSelectionSet(sset ast.SelectionSet) error {
	qb.push()
	for idx, selection := range sset {
		if idx > 0 {
			qb.newLine()
		}
		field, ok := selection.(*ast.Field)
		if !ok {
			return errors.New("unable to convert selection to field")
		}

		qb.writeField(field)

		if field.SelectionSet != nil {
			qb.writeR(' ')
			qb.writeSelectionSet(field.SelectionSet)
		}
	}
	qb.pop()
	return nil
}

/*
	type Filter struct {
		UIDs *[]UID
		Name *StringTermFilter
		And  *Filter
		Or   *Filter
		Not  *Filter
	}

	type QueryArgs struct {
		Filter *Filter
		Order  *Order
		First  *int
		Offset *int
	}
*/

func (qb *queryBuilder) writeQueryArgs(sset ast.SelectionSet) error {
	hasWritten := false

	qb.write("(func: ")

	if qb.args.Filter != nil {
		hasWritten = true
		if qb.args.Filter.UIDs != nil {
			qb.write("uid(")
			qb.write(strings.Join(schema.UIDsToStrings(*qb.args.Filter.UIDs), ", "))
			qb.writeR(')')
		}

		if qb.args.Filter.String != nil {
			filter := qb.args.Filter.String

			if filter.AnyOfTerms != nil {
				qb.write("anyofterms(")
				qb.write(filter.Name)
				qb.write(`, "`)
				qb.write(*filter.AnyOfTerms)
				qb.write(`")`)
			}

			if filter.AllOfTerms != nil {
				qb.write("allofterms(")
				qb.write(filter.Name)
				qb.write(`, "`)
				qb.write(*filter.AllOfTerms)
				qb.write(`")`)
			}
		}
	}

	if !hasWritten && sset != nil {
		qb.write("has(")
		names := []string{}
		for _, selection := range sset {
			field, ok := selection.(*ast.Field)
			if !ok {
				return errors.New("unable to convert selection to field")
			}

			if _, ok := specialFields[field.Name]; ok {
				continue
			}

			names = append(names, field.Name)
		}
		qb.write(strings.Join(names, ", "))
		qb.writeR(')')
	}

	if qb.args.First != nil {
		qb.write(", first: ")
		qb.writeInt(*qb.args.First)
	}

	if qb.args.Offset != nil {
		qb.write(", offset: ")
		qb.writeInt(*qb.args.Offset)
	}

	if qb.args.Order != nil {
		if qb.args.Order.Asc != nil {
			qb.write(", orderasc: ")
			qb.write(*qb.args.Order.Asc)
		} else if qb.args.Order.Desc != nil {
			qb.write(", orderdesc: ")
			qb.write(*qb.args.Order.Desc)
		}
	}

	qb.write(") ")

	return nil
}

func (qb *queryBuilder) writeField(field *ast.Field) {
	qb.leftPad()
	qb.write(field.Name)
}
