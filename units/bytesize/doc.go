// Package bytesize represents a number of bytes as a [Size], formats it for people, and parses it from text.
//
// SI units are powers of 1000: kB, MB, GB, TB, PB, EB. IEC units are powers of 1024: KiB, MiB, GiB, TiB, PiB, EiB.
// [Size.SI] and [Size.IEC] format a size in the units of their system, and [Size.String] formats it in IEC units.
// [Parse] reads the units of both systems, and a [Size] field of a configuration or a JSON document decodes from the
// same text. A JSON document can also give the size as a number of bytes.
package bytesize
