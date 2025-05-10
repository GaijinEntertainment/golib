package fields_test

import (
	"strconv"
	"testing"

	"dev.gaijin.team/go/golib/fields"
)

// goos: windows
// goarch: amd64
// pkg: dev.gaijin.team/go/golib/fields
// cpu: AMD Ryzen 9 7950X 16-Core Processor
// Benchmark_List
//
//	list_benchmark_test.go:29: dataset uniquiness: 0.7
//
// Benchmark_List/10
// Benchmark_List/10-32		10873563	120.2 ns/op		320 B/op	1 allocs/op
//
//	list_benchmark_test.go:29: dataset uniquiness: 0.55
//
// Benchmark_List/100
// Benchmark_List/100-32	465217		2913 ns/op		11648B/op	5 allocs/op
//
//	list_benchmark_test.go:29: dataset uniquiness: 0.604
//
// Benchmark_List/1000
// Benchmark_List/1000-32	51418		23741 ns/op		86656B/op	8 allocs/op
// PASS.
func Benchmark_List(b *testing.B) {
	const prealloc = 10

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		dataset := generateDataset(size)
		b.Logf("dataset uniquiness: %g", float64(len(uniqSlice(dataset)))/float64(size))

		b.Run(strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				list := make(fields.List, 0, prealloc)

				for _, key := range dataset {
					list.Add(fields.F(key, "val"))
				}
			}
		})
	}
}
