package domain

// ParameterType represents the type of a parameter
type ParameterType string

const (
	ParameterTypeString  ParameterType = "string"
	ParameterTypeInteger ParameterType = "integer"
	ParameterTypeNumber  ParameterType = "number"
	ParameterTypeBoolean ParameterType = "boolean"
	ParameterTypeObject  ParameterType = "object"
	ParameterTypeArray   ParameterType = "array"
)

// Parameter represents a parameter for an operation
type Parameter struct {
	Name        string        `json:"name"`
	Type        ParameterType `json:"type"`
	Description string        `json:"description"`
	Required    bool          `json:"required"`
	Enum        []string      `json:"enum,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
}

func NewParameter(name string, parameterType ParameterType, description string, required bool, enum []string, defaultVal interface{}) *Parameter {
	return &Parameter{
		Name:        name,
		Type:        parameterType,
		Description: description,
		Required:    required,
		Enum:        enum,
		Default:     defaultVal,
	}
}
