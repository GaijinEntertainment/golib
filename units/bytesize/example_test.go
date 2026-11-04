package bytesize_test

import (
	"fmt"

	"dev.gaijin.team/go/golib/units/bytesize"
)

func ExampleSize_IEC() {
	fmt.Println(bytesize.Size(1536).IEC(2))
	fmt.Println(bytesize.Size(1048575).IEC(2))
	// Output:
	// 1.5KiB
	// 1MiB
}

func ExampleSize_SI() {
	fmt.Println(bytesize.Size(1536).SI(2))
	fmt.Println(bytesize.Size(256 * bytesize.Kilobyte).SI(0))
	// Output:
	// 1.54kB
	// 256kB
}

func ExampleParse() {
	for _, text := range []string{"1536", "1.5 KiB", "256kB", "1XB"} {
		size, err := bytesize.Parse(text)
		if err != nil {
			fmt.Println(err)

			continue
		}

		fmt.Println(int64(size))
	}

	// Output:
	// 1536
	// 1536
	// 256000
	// invalid byte size (text=1XB)
}
