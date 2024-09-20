package fields_test

import (
	"strings"
	"testing"

	"github.com/GaijinEntertainment/golib/fields"

	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	t.Parallel()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()

		d := fields.Dict{"foo": "bar"}

		require.Equal(t, fields.Dict{"foo": "bar"}, d)

		d.Add(fields.F("baz", "qux"), fields.F("foo", "baz"))
		require.Equal(t, fields.Dict{"foo": "baz", "baz": "qux"}, d)
	})

	t.Run("ToList", func(t *testing.T) {
		t.Parallel()

		d := fields.Dict{"foo": "bar", "baz": "qux"}

		require.ElementsMatch(t, fields.List{{"foo", "bar"}, {"baz", "qux"}}, d.ToList())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		d := fields.Dict{}
		require.Equal(t, "", d.String())

		d = fields.Dict{"foo": "bar"}
		require.ElementsMatch(t, []string{"foo=bar"}, collectionStringToKVElements(d.String()))

		d = fields.Dict{"foo": "bar", "baz": "qux"}
		require.ElementsMatch(t, []string{"foo=bar", "baz=qux"}, collectionStringToKVElements(d.String()))

	})
}

func collectionStringToKVElements(str string) []string {
	if len(str) < 2 {
		return nil
	}

	return strings.Split(str[1:len(str)-1], fields.CollectionSep)
}
