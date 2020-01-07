package schema

type ExecutionResult struct {
	Data struct {
		Nodes []Node `json:"nodes"`
	} `json:"data"`
	Errors []error `json:"errors,omitempty"`
}

//func (Schema) Execute(source string, vars Node) (*ExecutionResult, error) {
//_, err := parser.ParseQuery(&ast.Source{
//Name:  "query.graphql",
//Input: source,
//})
//if err != nil {
//return nil, err
//}

//fmt.Printf("%+v %+v\n", source, vars)

//return nil, nil
//}
