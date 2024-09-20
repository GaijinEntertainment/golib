package fields

import (
	"strings"
)

// List is list of fields in order of insertion.
type List []Field

// Add adds fields to a [List], modifying it.
func (l *List) Add(fields ...Field) {
	*l = append(*l, fields...)
}

// ToDict converts a [List] to a [Dict] overwriting existing keys.
func (l List) ToDict() Dict {
	d := make(Dict, len(l))

	for i := range len(l) {
		d[l[i].K] = l[i].V
	}

	return d
}

// WriteTo writes a string representation of a [List] to a given builder in the
// `({key}={value}, {key}={value})` format.
func (l List) WriteTo(b *strings.Builder) {
	if len(l) == 0 {
		return
	}

	b.WriteString("(")

	sep := false
	for i := 0; i < len(l); i++ {
		if sep {
			b.WriteString(", ")
		}

		writeKVTo(b, l[i].K, l[i].V)

		sep = true
	}

	b.WriteString(")")
}

// String returns a string representation of a [List] in the
// `({key}={value}, {key}={value})` format.
func (l List) String() string {
	b := strings.Builder{}

	l.WriteTo(&b)

	return b.String()
}
