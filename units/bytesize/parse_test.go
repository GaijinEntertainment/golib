package bytesize_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/units/bytesize"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("accepts", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give string
			want bytesize.Size
		}{
			{name: "bytes", give: "1536", want: 1536},
			{name: "byte symbol", give: "1536B", want: 1536},
			{name: "SI unit", give: "256kB", want: 256 * bytesize.Kilobyte},
			{name: "IEC unit", give: "1MiB", want: bytesize.Mebibyte},
			{name: "SI unit in other letter case", give: "1KB", want: bytesize.Kilobyte},
			{name: "IEC unit in other letter case", give: "1kib", want: bytesize.Kibibyte},
			{name: "space before the unit", give: "1 KiB", want: bytesize.Kibibyte},
			{name: "spaces around the text", give: " 1KiB\n", want: bytesize.Kibibyte},
			{name: "exact fraction", give: "4.1MB", want: 4_100_000},
			{name: "fraction rounds to the nearest byte", give: "1.1MiB", want: 1_153_434},
			{name: "half rounds away from zero", give: "0.5B", want: 1},
			{name: "negative half rounds away from zero", give: "-0.5B", want: -1},
			{name: "negative size", give: "-1.5KiB", want: -1536},
			{name: "positive sign", give: "+1KiB", want: bytesize.Kibibyte},
			{name: "largest unit", give: "8EB", want: 8 * bytesize.Exabyte},
			{name: "largest size", give: "9223372036854775807", want: math.MaxInt64},
			{name: "largest size with a fraction", give: "9223372036854775807.4", want: math.MaxInt64},
			{name: "smallest size", give: "-9223372036854775808", want: math.MinInt64},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := bytesize.Parse(tt.give)
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("rejects", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			give    string
			wantErr error
		}{
			{name: "empty text", give: "", wantErr: bytesize.ErrSyntax},
			{name: "unit without a number", give: "KiB", wantErr: bytesize.ErrSyntax},
			{name: "point without a fraction", give: "1.KiB", wantErr: bytesize.ErrSyntax},
			{name: "fraction without an integer", give: ".5KiB", wantErr: bytesize.ErrSyntax},
			{name: "unknown unit", give: "1XB", wantErr: bytesize.ErrSyntax},
			{name: "exponent", give: "1e3", wantErr: bytesize.ErrSyntax},
			{name: "decimal comma", give: "1,5MB", wantErr: bytesize.ErrSyntax},
			{name: "space inside the unit", give: "1 Ki B", wantErr: bytesize.ErrSyntax},
			{name: "beyond the largest size", give: "9223372036854775808", wantErr: bytesize.ErrRange},
			{name: "rounding beyond the largest size", give: "9223372036854775807.5", wantErr: bytesize.ErrRange},
			{name: "beyond the largest size in units", give: "8EiB", wantErr: bytesize.ErrRange},
			{name: "beyond the smallest size", give: "-9223372036854775809", wantErr: bytesize.ErrRange},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				_, err := bytesize.Parse(tt.give)
				require.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
