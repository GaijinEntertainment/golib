// Package build provides tools for accessing and working with build information
// of the binary.
//
// The package automatically extracts build information from Go's runtime debug data,
// including Go version, target OS and architecture, and version control system details.
// It also provides methods to format this information as human-readable text or structured
// data.
//
// The Info struct has the following fields:
//   - Version (user-defined, settable via ldflags or taken from debug info)
//   - BuildTime (user-defined, settable via ldflags)
//   - GoVersion
//   - GoOS
//   - GoArch
//   - VCSRevision
//   - VCSModified
//
// The ToFields() method returns fields in the same order as described.
//
// Any of following fields can be set via ldflags, and, if so, they won't be overridden
// by parsed info, except for GoVersion.
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
// via ldflags as follows (note the variables names):
//
//	go build -ldflags "-X 'dev.gaijin.team/go/golib/build.version=v1.0.0' -X 'dev.gaijin.team/go/golib/build.buildTime=<current-time>'"
//
// If not set, these fields will be populated from Go build info or default to "unknown".
//nolint
package build
