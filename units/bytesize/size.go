package bytesize

import (
	"math/big"
	"strings"
)

// Size is a number of bytes.
type Size int64

// Units of a [Size]. The constants are untyped, so they also size values of other integer types, such as
// make([]byte, 64*bytesize.Kibibyte).
const (
	Byte = 1

	Kilobyte = 1000 * Byte
	Megabyte = 1000 * Kilobyte
	Gigabyte = 1000 * Megabyte
	Terabyte = 1000 * Gigabyte
	Petabyte = 1000 * Terabyte
	Exabyte  = 1000 * Petabyte

	Kibibyte = 1024 * Byte
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
	Tebibyte = 1024 * Gibibyte
	Pebibyte = 1024 * Tebibyte
	Exbibyte = 1024 * Pebibyte
)

const stringPrecision = 2

type unit struct {
	symbol string
	size   int64
}

//nolint:gochecknoglobals // Read-only unit tables, in ascending order.
var (
	siUnits = []unit{
		{"B", Byte}, {"kB", Kilobyte}, {"MB", Megabyte}, {"GB", Gigabyte}, {"TB", Terabyte}, {"PB", Petabyte},
		{"EB", Exabyte},
	}
	iecUnits = []unit{
		{"B", Byte}, {"KiB", Kibibyte}, {"MiB", Mebibyte}, {"GiB", Gibibyte}, {"TiB", Tebibyte}, {"PiB", Pebibyte},
		{"EiB", Exbibyte},
	}
)

// String formats s as [Size.IEC] does with precision 2, such as "1.5KiB".
func (s Size) String() string {
	return s.IEC(stringPrecision)
}

// SI formats s in an SI unit from B to EB: the smallest one in which the absolute value, rounded half away from zero
// to precision fractional digits, is below 1000. It drops the trailing zeros of the fraction, so 1500 bytes with
// precision 2 is "1.5kB". A negative precision counts as zero.
func (s Size) SI(precision int) string {
	return s.format(siUnits, precision)
}

// IEC formats s in an IEC unit from B to EiB: the smallest one in which the absolute value, rounded half away from
// zero to precision fractional digits, is below 1024. It drops the trailing zeros of the fraction, so 1536 bytes with
// precision 2 is "1.5KiB". A negative precision counts as zero.
func (s Size) IEC(precision int) string {
	return s.format(iecUnits, precision)
}

func (s Size) format(units []unit, precision int) string {
	precision = max(precision, 0)

	base := big.NewRat(units[1].size, 1)

	u := units[0]
	num := s.digitsIn(u, precision)

	for _, next := range units[1:] {
		if rounded, _ := new(big.Rat).SetString(num); rounded.Abs(rounded).Cmp(base) < 0 {
			break
		}

		u, num = next, s.digitsIn(next, precision)
	}

	return trimFraction(num) + u.symbol
}

func (s Size) digitsIn(u unit, precision int) string {
	return new(big.Rat).SetFrac64(int64(s), u.size).FloatString(precision)
}

func trimFraction(num string) string {
	if !strings.Contains(num, ".") {
		return num
	}

	return strings.TrimSuffix(strings.TrimRight(num, "0"), ".")
}
