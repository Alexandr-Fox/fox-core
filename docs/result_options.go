package docs

type ResultOptions struct {
	Items *[]Result `json:"items,omitempty"`
	Item  *Result   `json:"item,omitempty"`
	Enum  *[]string `json:"enum,omitempty"`
}
