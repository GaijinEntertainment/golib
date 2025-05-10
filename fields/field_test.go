package fields_test

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/fields"
)

func TestField(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			name     string
			in       fields.Field
			expected string
		}{
			{
				name:     "string value",
				in:       fields.F("key", "value"),
				expected: "key=value",
			},
			{
				name:     "int value",
				in:       fields.F("key", 42),
				expected: "key=42",
			},
			{
				name:     "error value",
				in:       fields.F("key", errors.New("error")), //nolint:err113
				expected: "key=error",
			},
			{
				name:     "stringer value",
				in:       fields.F("key", net.IP{192, 168, 1, 1}),
				expected: "key=192.168.1.1",
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) { //nolint:paralleltest
				if got := tc.in.String(); got != tc.expected {
					assert.Equal(t, tc.expected, got)
				}
			})
		}
	})
}
