package docs

type FieldOptions struct {
	Items *[]Field  `json:"items,omitempty"`
	Enum  *[]string `json:"enum,omitempty"`
}
