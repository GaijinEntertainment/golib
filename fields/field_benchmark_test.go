package fields_test

import (
	"strconv"
	"testing"
)

// goos: windows
// goarch: amd64
// pkg: dev.gaijin.team/go/golib/fields
// cpu: AMD Ryzen 9 7950X 16-Core Processor
// Benchmark_Container
// Benchmark_Container/any_slice/10
// Benchmark_Container/any_slice/10-32				   12352566		88.78 ns/op		  320 B/op	1 allocs/op
// Benchmark_Container/struct_slice_direct/10
// Benchmark_Container/struct_slice_direct/10-32	   13581465		84.25 ns/op		  320 B/op	1 allocs/op
// Benchmark_Container/struct_slice_constructor/10
// Benchmark_Container/struct_slice_constructor/10-32  15991406		82.81 ns/op		  320 B/op	1 allocs/op
// Benchmark_Container/any_slice/100
// Benchmark_Container/any_slice/100-32					2184120		492.0 ns/op		 3456 B/op	1 allocs/op
// Benchmark_Container/struct_slice_direct/100
// Benchmark_Container/struct_slice_direct/100-32		2369796		506.0 ns/op		 3456 B/op	1 allocs/op
// Benchmark_Container/struct_slice_constructor/100
// Benchmark_Container/struct_slice_constructor/100-32	2096652		497.8 ns/op		 3456 B/op	1 allocs/op
// Benchmark_Container/any_slice/1000
// Benchmark_Container/any_slice/1000-32				 230362		 4502 ns/op		32768 B/op	1 allocs/op
// Benchmark_Container/struct_slice_direct/1000
// Benchmark_Container/struct_slice_direct/1000-32		 231802		 4666 ns/op		32768 B/op	1 allocs/op
// Benchmark_Container/struct_slice_constructor/1000
// Benchmark_Container/struct_slice_constructor/1000-32	 280346		 4505 ns/op		32768 B/op	1 allocs/op
// PASS.
func Benchmark_Container(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		benchAnySlice(b, size)
		benchStructDirect(b, size)
		benchStructConstructor(b, size)
	}
}

type anySlice []any

func benchAnySlice(b *testing.B, size int) {
	b.Helper()

	b.Run("any slice/"+strconv.Itoa(size), func(b *testing.B) {
		b.ReportAllocs()

		var arr anySlice

		for i := 0; i < b.N; i++ {
			arr = make(anySlice, 0, size*2)
			for range size {
				arr = append(arr, "key", "value")
			}
		}

		b.Logf("%d", len(arr))
	})
}

type anyStruct struct {
	K string
	V any
}

type anyStructSlice []anyStruct

func benchStructDirect(b *testing.B, size int) {
	b.Helper()

	b.Run("struct slice direct/"+strconv.Itoa(size), func(b *testing.B) {
		b.ReportAllocs()

		var arr anyStructSlice

		for i := 0; i < b.N; i++ {
			arr = make(anyStructSlice, 0, size)
			for range size {
				arr = append(arr, anyStruct{"key", "value"})
			}
		}

		b.Logf("%d", len(arr))
	})
}

func newAnyStruct(k string, v any) anyStruct {
	return anyStruct{k, v}
}

func benchStructConstructor(b *testing.B, size int) {
	b.Helper()

	b.Run("struct slice constructor/"+strconv.Itoa(size), func(b *testing.B) {
		b.ReportAllocs()

		var arr anyStructSlice

		for i := 0; i < b.N; i++ {
			arr = make(anyStructSlice, 0, size)
			for range size {
				arr = append(arr, newAnyStruct("key", "value"))
			}
		}

		b.Logf("%d", len(arr))
	})
}
