package bytesize_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/units/bytesize"
)

func Test_Size(t *testing.T) {
	t.Parallel()

	t.Run("SI and IEC", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name          string
			give          bytesize.Size
			givePrecision int
			wantSI        string
			wantIEC       string
		}{
			{name: "zero", give: 0, givePrecision: 2, wantSI: "0B", wantIEC: "0B"},
			{name: "bytes", give: 1023, givePrecision: 2, wantSI: "1.02kB", wantIEC: "1023B"},
			{name: "bytes below the base stay in bytes", give: 999, givePrecision: 0, wantSI: "999B", wantIEC: "999B"},
			{name: "whole unit", give: 1024, givePrecision: 2, wantSI: "1.02kB", wantIEC: "1KiB"},
			{name: "fraction with a leading zero", give: 1060, givePrecision: 2, wantSI: "1.06kB", wantIEC: "1.04KiB"},
			{name: "fraction rounded to zero", give: 1026, givePrecision: 2, wantSI: "1.03kB", wantIEC: "1KiB"},
			{name: "trailing zero dropped", give: 1536, givePrecision: 2, wantSI: "1.54kB", wantIEC: "1.5KiB"},
			{name: "IEC rounding carries", give: 1048575, givePrecision: 2, wantSI: "1.05MB", wantIEC: "1MiB"},
			{name: "SI rounding carries", give: 999_999, givePrecision: 2, wantSI: "1MB", wantIEC: "976.56KiB"},
			{name: "zero precision", give: 1536, givePrecision: 0, wantSI: "2kB", wantIEC: "2KiB"},
			{name: "negative precision", give: 1536, givePrecision: -1, wantSI: "2kB", wantIEC: "2KiB"},
			{name: "half rounds away from zero", give: 2560, givePrecision: 0, wantSI: "3kB", wantIEC: "3KiB"},
			{name: "precision beyond float64", give: 1024, givePrecision: 17, wantSI: "1.024kB", wantIEC: "1KiB"},
			{name: "negative size", give: -1536, givePrecision: 2, wantSI: "-1.54kB", wantIEC: "-1.5KiB"},
			{name: "largest size", give: math.MaxInt64, givePrecision: 2, wantSI: "9.22EB", wantIEC: "8EiB"},
			{name: "smallest size", give: math.MinInt64, givePrecision: 2, wantSI: "-9.22EB", wantIEC: "-8EiB"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, tt.wantSI, tt.give.SI(tt.givePrecision))
				assert.Equal(t, tt.wantIEC, tt.give.IEC(tt.givePrecision))
			})
		}
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "1.5KiB", bytesize.Size(1536).String())
		assert.Equal(t, "3MiB", bytesize.Size(3*bytesize.Mebibyte).String())
	})
}
