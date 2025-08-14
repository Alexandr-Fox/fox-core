package docs

type Result struct {
	Name     string         `json:"name"`
	Type     Type           `json:"type"`
	Optional bool           `json:"optional"`
	Options  *ResultOptions `json:"options,omitempty"`
}
