package fields

import (
	"strings"
)

const CollectionSep = ", "

// Dict is a dictionary of unique fields.
type Dict map[string]any

// Add adds fields to a [Dict] overwriting existing keys.
func (d Dict) Add(fields ...Field) {
	for _, f := range fields {
		d[f.K] = f.V
	}
}

// ToList converts a [Dict] to a [List].
func (d Dict) ToList() List {
	s := make(List, 0, len(d))

	for k, v := range d {
		s = append(s, Field{k, v})
	}

	return s
}

// WriteTo writes a string representation of a [Dict] to a given builder in the
// `({key}={value}, {key}={value})` format.
func (d Dict) WriteTo(b *strings.Builder) {
	if len(d) == 0 {
		return
	}

	b.WriteString("(")

	sep := false
	for k, v := range d {
		if sep {
			b.WriteString(CollectionSep)
		}

		writeKVTo(b, k, v)

		sep = true
	}

	b.WriteString(")")
}

// String returns a string representation of a [Dict] in the
// `({key}={value}, {key}={value})` format.
func (d Dict) String() string {
	b := strings.Builder{}

	d.WriteTo(&b)

	return b.String()
}
