package fields_test

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"dev.gaijin.team/go/golib/fields"
)

// goos: windows
// goarch: amd64
// pkg: dev.gaijin.team/go/golib/fields
// cpu: AMD Ryzen 9 7950X 16-Core Processor
// Benchmark_Dict
//
//	dict_benchmark_test.go:31: dataset uniquiness: 0.7
//
// Benchmark_Dict/10
// Benchmark_Dict/10-32		5656870		219.7 ns/op		576 B/op	1 allocs/op
//
//	dict_benchmark_test.go:31: dataset uniquiness: 0.64
//
// Benchmark_Dict/100
// Benchmark_Dict/100-32	292149		4511 ns/op		9663B/op	7 allocs/op
//
//	dict_benchmark_test.go:31: dataset uniquiness: 0.603
//
// Benchmark_Dict/1000
// Benchmark_Dict/1000-32	28310		41133 ns/op		83231B/op         24 allocs/op
// PASS.
func Benchmark_Dict(b *testing.B) {
	const prealloc = 10

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		dataset := generateDataset(size)
		b.Logf("dataset uniquiness: %g", float64(len(uniqSlice(dataset)))/float64(size))

		b.Run(strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				dict := make(fields.Dict, prealloc)

				for _, key := range dataset {
					dict.Add(fields.F(key, "val"))
				}
			}
		})
	}
}

// In order to properly test the performance of the dictionary, we need to
// generate dataset with a certain hitrate of overwrites.
func generateDataset(size int) []string {
	const hitRate = 0.1

	result := make([]string, size)

	for i := range size {
		result[i] = strconv.Itoa(rand.Intn(size - int(math.Floor(float64(size)*hitRate)) + 1)) //nolint:gosec
	}

	return result
}

// uniqSlice returns a slice with unique elements from the input slice.
func uniqSlice[T comparable](s []T) []T {
	m := make(map[T]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}
