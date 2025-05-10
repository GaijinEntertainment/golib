package fields_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/fields"
)

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()

		l := fields.List{{"foo", "baz"}}
		l.Add(fields.F("foo", "bar"), fields.F("baz", "qux"))
		require.Equal(t, fields.List{{"foo", "baz"}, {"foo", "bar"}, {"baz", "qux"}}, l)
	})

	t.Run("ToDict", func(t *testing.T) {
		t.Parallel()

		l := fields.List{{"foo", "bar"}, {"foo", "baz"}, {"baz", "qux"}}
		require.Equal(t, fields.Dict{"foo": "baz", "baz": "qux"}, l.ToDict())

		l = fields.List{{"foo", "baz"}, {"foo", "bar"}, {"baz", "qux"}}
		require.Equal(t, fields.Dict{"foo": "bar", "baz": "qux"}, l.ToDict())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		l := fields.List{}
		require.Empty(t, l.String())

		l = fields.List{{"foo", "bar"}}
		require.Equal(t, "(foo=bar)", l.String())

		l = fields.List{{"foo", "bar"}, {"baz", "qux"}}
		require.Equal(t, "(foo=bar, baz=qux)", l.String())
	})

	t.Run("All early exit", func(t *testing.T) {
		t.Parallel()

		l := fields.List{
			{"foo", "bar"},
			{"baz", "qux"},
		}

		var seen []string
		for k := range l.All() {
			seen = append(seen, k)
			break // stop after first
		}

		require.Len(t, seen, 1)
	})
}
