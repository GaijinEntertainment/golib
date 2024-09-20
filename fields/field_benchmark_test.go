package fields_test

import (
	"testing"
)

type anyStruct struct {
	K string
	V any
}

func newAnyStruct(k string, v any) anyStruct {
	return anyStruct{k, v}
}

type anyStructSlice []anyStruct

type anySlice []any

// goos: windows
// goarch: amd64
// pkg: github.com/GaijinEntertainment/golib/fields
// cpu: AMD Ryzen 9 7950X 16-Core Processor
// Benchmark_Container
// Benchmark_Container/any_slice
// Benchmark_Container/any_slice-32                 13572188	87.82 ns/op		320 B/op	1 allocs/op
// Benchmark_Container/struct_slice_direct
// Benchmark_Container/struct_slice_direct-32		16315009	83.82 ns/op		320 B/op	1 allocs/op
// Benchmark_Container/struct_slice_constructor
// Benchmark_Container/struct_slice_constructor-32	13774736	89.33 ns/op		320 B/op	1 allocs/op
// PASS
func Benchmark_Container(b *testing.B) {
	b.Run("any slice", func(b *testing.B) {
		b.ReportAllocs()

		var arr anySlice

		for i := 0; i < b.N; i++ {
			arr = make(anySlice, 0, 20)
			for range 10 {
				arr = append(arr, "key", "value")
			}
		}

		b.Logf("%d", len(arr))
	})

	b.Run("struct slice direct", func(b *testing.B) {
		b.ReportAllocs()

		var arr anyStructSlice

		for i := 0; i < b.N; i++ {
			arr = make(anyStructSlice, 0, 10)
			for range 10 {
				arr = append(arr, anyStruct{"key", "value"})
			}
		}

		b.Logf("%d", len(arr))
	})

	b.Run("struct slice constructor", func(b *testing.B) {
		b.ReportAllocs()

		var arr anyStructSlice

		for i := 0; i < b.N; i++ {
			arr = make(anyStructSlice, 0, 10)
			for range 10 {
				arr = append(arr, newAnyStruct("key", "value"))
			}
		}

		b.Logf("%d", len(arr))
	})
}
