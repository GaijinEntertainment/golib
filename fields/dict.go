package fields

import (
	"strings"
)

const CollectionSep = ", "

// Dict is a map-based collection of unique fields, keyed by string.
// It provides efficient lookup and overwrites duplicate keys.
type Dict map[string]any

// Add inserts or updates fields in the Dict, overwriting existing keys if present.
//
// Example:
//
//	d := fields.Dict{"foo": "bar"}
//	d.Add(fields.F("baz", 42), fields.F("foo", "qux")) // d["foo"] == "qux"
func (d Dict) Add(fields ...Field) {
	for _, f := range fields {
		d[f.K] = f.V
	}
}

// ToList converts the Dict to a List, with order unspecified.
// Each key-value pair becomes a Field in the resulting List.
func (d Dict) ToList() List {
	s := make(List, 0, len(d))

	for k, v := range d {
		s = append(s, Field{k, v})
	}

	return s
}

// WriteTo writes the Dict as a string in the format "(key1=val1, key2=val2)" to the provided builder.
// If the Dict is empty, nothing is written. The order of fields is unspecified.
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

// String returns the Dict as a string in the format "(key1=val1, key2=val2)".
// Returns an empty string if the Dict is empty. The order of fields is unspecified.
func (d Dict) String() string {
	b := strings.Builder{}

	d.WriteTo(&b)

	return b.String()
}
