package bytesize

import (
	"encoding/json"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"dev.gaijin.team/go/golib/e"
)

var (
	// ErrSyntax reports text that is not a decimal number with an optional unit.
	ErrSyntax = e.New("invalid byte size")
	// ErrRange reports a size outside the range of [Size].
	ErrRange = e.New("byte size out of range")
)

var sizeRe = regexp.MustCompile(`^([+-]?[0-9]+(?:\.[0-9]+)?)\s*([A-Za-z]*)$`)

// Parse reads a size from a decimal number and an optional unit, such as "1536", "1.5KiB", or "256 kB". The unit is a
// symbol of [Size.SI] or [Size.IEC] in any letter case, so "kb" is 1000 bytes and "KIB" is 1024 bytes. A number
// without a unit counts bytes. The number can have a sign and a fraction with digits on both sides of the point, but
// no exponent. The size rounds to the nearest byte, half away from zero. Parse ignores spaces around the text and
// between the number and the unit.
//
// The returned error matches [ErrSyntax] for text of another form, and [ErrRange] for a size that a [Size] cannot hold.
func Parse(text string) (Size, error) {
	match := sizeRe.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return 0, ErrSyntax.WithField("text", text)
	}

	unitSize, isUnit := lookupUnit(match[2])
	if !isUnit {
		return 0, ErrSyntax.WithField("text", text)
	}

	value, _ := new(big.Rat).SetString(match[1])

	// FloatString rounds half away from zero.
	digits := value.Mul(value, new(big.Rat).SetInt64(unitSize)).FloatString(0)

	size, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return 0, ErrRange.WithField("text", text)
	}

	return Size(size), nil
}

// UnmarshalText sets s from text that [Parse] accepts.
func (s *Size) UnmarshalText(text []byte) error {
	size, err := Parse(string(text))
	if err != nil {
		return err
	}

	*s = size

	return nil
}

// UnmarshalJSON sets s from a JSON string or a JSON number, each read as [Parse] text. Its error matches [ErrSyntax]
// for the number 1e3 and for a JSON value of another type. A JSON null leaves s unchanged.
func (s *Size) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	text := string(data)

	if strings.HasPrefix(text, `"`) {
		if err := json.Unmarshal(data, &text); err != nil {
			return e.NewFrom("decode byte size", err)
		}
	}

	return s.UnmarshalText([]byte(text))
}

func lookupUnit(symbol string) (int64, bool) {
	if symbol == "" {
		return Byte, true
	}

	for _, units := range [...][]unit{siUnits, iecUnits} {
		for _, u := range units {
			if strings.EqualFold(u.symbol, symbol) {
				return u.size, true
			}
		}
	}

	return 0, false
}
