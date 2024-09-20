package fields

import (
	"fmt"
	"strings"
)

// Field is a generic key-value pair container.
type Field struct {
	K string
	V any
}

// F is a shorthand for creating a new Field.
func F(key string, value any) Field {
	return Field{K: key, V: value}
}

func writeKVTo(b *strings.Builder, key string, value any) {
	b.WriteString(key)
	b.WriteRune('=')

	switch val := value.(type) {
	case string:
		b.WriteString(val)

	case fmt.Stringer:
		b.WriteString(val.String())

	case error:
		b.WriteString(val.Error())

	default:
		b.WriteString(fmt.Sprintf("%v", val))
	}
}

// WriteTo writes a string representation of a [Field] to a given builder in the
// `{key}={value}`.
func (f Field) WriteTo(b *strings.Builder) {
	writeKVTo(b, f.K, f.V)
}

// String returns a string representation of a [Field] in the `{key}={value}`
// format.
func (f Field) String() string {
	b := &strings.Builder{}
	f.WriteTo(b)

	return b.String()
}
