//nolint
package build_test

import (
	"fmt"
	"log"

	"dev.gaijin.team/go/golib/build"
)

func Example() {
	// Get build information
	info := build.GetInfo()

	// Print human-readable build info
	fmt.Println("Build information:", info)

	// Access specific fields
	fmt.Printf("Built with Go %s for %s/%s\n", info.GoVersion, info.GoOS, info.GoArch)

	// Check if built from modified source
	if info.VCSModified == "true" {
		fmt.Println("Warning: Built from modified source!")
	}
}

func ExampleInfo_ToFields() {
	// Get build information
	info := build.GetInfo()

	// Convert to fields.List for use in structured logging
	fields := info.ToFields()

	// Use with custom logger (just printing here for the example)
	log.Printf("Application starting with build info: %s", fields)

	// In a real application, you might do something like:
	// logger.Info("Application starting", fields...)
}
