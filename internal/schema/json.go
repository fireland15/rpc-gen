package schema

// JsonType represents a JSON primitive type.
//
// It is used to identify the underlying primitive kind of a JSON value
// as defined by the JSON specification.
type JsonType string

const (
	// JsonTypeBoolean represents a JSON boolean value (true or false).
	JsonTypeBoolean JsonType = "boolean"

	// JsonTypeNumber represents a JSON numeric value.
	JsonTypeNumber JsonType = "number"

	// JsonTypeString represents a JSON string value.
	JsonTypeString JsonType = "string"
)

// isValidJsonType ensure `t` is a valid value.
func isValidJsonType(t JsonType) bool {
	switch t {
	case JsonTypeBoolean, JsonTypeNumber, JsonTypeString:
		return true
	default:
		return false
	}
}
