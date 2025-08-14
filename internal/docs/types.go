package docs

type Type string

const (
	String Type = "string"
	Int    Type = "int"
	Float  Type = "float"
	Bool   Type = "bool"
	Array  Type = "array"
	Object Type = "object"
	Enum   Type = "enum"
)

func (t *Type) ToPointer() *Type {
	r := *t
	return &r
}

func ToPointer(t Type) *Type {
	return &t
}
