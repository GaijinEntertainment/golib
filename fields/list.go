package fields

import (
	"strings"
)

// List is an ordered collection of Field values, preserving insertion order.
// Such collection do not check for duplicate keys.
type List []Field

// Add appends one or more fields to the List, modifying it.
//
// Example:
//
//	var l fields.List
//	l.Add(fields.F("foo", "bar"), fields.F("baz", 42))
func (l *List) Add(fields ...Field) {
	*l = append(*l, fields...)
}

// ToDict converts the List to a Dict, overwriting duplicate keys with the last occurrence.
//
// Example:
//
//	l := fields.List{fields.F("foo", 1), fields.F("foo", 2)}
//	d := l.ToDict() // d["foo"] == 2
func (l List) ToDict() Dict {
	d := make(Dict, len(l))

	for i := range l {
		d[l[i].K] = l[i].V
	}

	return d
}

// WriteTo writes the List as a string in the format "(key1=val1, key2=val2)" to the provided builder.
// If the List is empty, nothing is written.
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

// String returns the List as a string in the format "(key1=val1, key2=val2)".
// Returns an empty string if the List is empty.
func (l List) String() string {
	b := strings.Builder{}

	l.WriteTo(&b)

	return b.String()
}
