// Package build provides tools for accessing and working with build information
// of the binary.
//
// The package automatically extracts build information from Go's runtime debug data,
// including Go version, target OS and architecture, and version control system details.
// It also provides methods to format this information as human-readable text or structured data.
//
// # Usage
//
// To get the current build information:
//
//	info := build.GetInfo()
//	fmt.Println(info)  // Prints human-readable build info
//
// To access specific build details:
//
//	info := build.GetInfo()
//	if info.VCSModified == "true" {
//		fmt.Println("Built from modified source!")
//	}
//
// # Field Conversion
//
// Convert build info to fields.List for integration with logging systems:
//
//	logger.Info("Starting application", build.GetInfo().ToFields()...)
//
// # Custom Build Information
//
// Info can also include custom information that is expected to be set by developer at build time
// via ldflags as follows:
//
//	`go build -ldflags "-X 'dev.gaijin.team/go/golib/build.time=<current-time>' -X 'dev.gaijin.team/go/golib/build.vcsTag=v1.0.0'"``
//
//nolint
package build
