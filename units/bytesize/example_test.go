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
