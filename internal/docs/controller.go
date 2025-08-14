package docs

type ControllerDoc struct {
	Fields  *Field  `json:"fields,omitempty"`
	Results *Result `json:"results,omitempty"`
}
