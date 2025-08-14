package docs

type Field struct {
	Name     string        `json:"name"`
	Desc     string        `json:"desc"`
	Type     Type          `json:"type"`
	Required bool          `json:"required"`
	Options  *FieldOptions `json:"options,omitempty"`
}
