package build

import (
	"runtime/debug"
	"strings"

	"dev.gaijin.team/go/golib/fields"
)

const unknown = "unknown"

//nolint:gochecknoglobals
var (
	buildInfoParsed = false

	goVersion = unknown
	goOS      = unknown
	goArch    = unknown

	vcsRevision = unknown
	vcsModified = unknown

	buildTime = unknown
	version   = unknown
)

// Info represents the build information for a Go binary.
// It contains details about the Go environment used for building,
// version control system state, and custom build-time information.
type Info struct {
	//
	// User-defined information, set at build time.
	//

	// Version is the version of the binary. In case parameter is not set,
	// via linker flags it will be taken from debug.ReadBuildInfo().
	Version string `json:"version"`
	// BuildTime is the build time.
	BuildTime string `json:"build-time"`

	//
	// Basic go information, detected automatically.
	//

	// GoVersion is the version of Go used to build the binary.
	GoVersion string `json:"go-version"`
	// GoOS is the operating system used to build the binary.
	GoOS string `json:"go-os"`
	// GoArch is the architecture used to build the binary.
	GoArch string `json:"go-arch"`

	//
	// VCS information, detected automatically, if available.
	//

	// VCSRevision is the version control system revision.
	VCSRevision string `json:"vcs-revision"`
	// VCSModified is the version control system modified status. "true" in case
	// repo is "dirty".
	VCSModified string `json:"vcs-modified"`
}

// String returns a human-readable representation of build information.
// The format includes version, Go version, OS, architecture, VCS revision, and build details.
// If the repository was modified (dirty), this is indicated in the output.
func (i Info) String() string {
	b := &strings.Builder{}

	b.WriteString(i.Version)
	b.WriteString(" built with ")
	b.WriteString(i.GoVersion)
	b.WriteString(" for ")
	b.WriteString(i.GoOS)
	b.WriteString(" ")
	b.WriteString(i.GoArch)
	b.WriteString(" from ")
	b.WriteString(i.VCSRevision)

	if i.VCSModified == "true" {
		b.WriteString(" (dirty)")
	}

	b.WriteString(" build time ")
	b.WriteString(i.BuildTime)

	return b.String()
}

// ToFields converts Info to a fields.List containing all build information.
// This is useful for structured logging or when you need to process
// build information as key-value pairs.
func (i Info) ToFields() fields.List {
	return fields.List{
		fields.F("version", i.Version),
		fields.F("go-version", i.GoVersion),
		fields.F("go-os", i.GoOS),
		fields.F("go-arch", i.GoArch),
		fields.F("vcs-revision", i.VCSRevision),
		fields.F("vcs-modified", i.VCSModified),
		fields.F("build-time", i.BuildTime),
	}
}

// GetInfo returns the build information.
// It automatically parses build information from runtime debug data
// the first time it's called, and caches the result for subsequent calls.
// Values set via linker flags (-X) are preserved and not overwritten.
func GetInfo() Info {
	if !buildInfoParsed {
		buildInfoParsed = true

		parseBuildInfo()
	}

	return Info{
		Version:     version,
		GoVersion:   goVersion,
		GoOS:        goOS,
		GoArch:      goArch,
		VCSRevision: vcsRevision,
		VCSModified: vcsModified,
		BuildTime:   buildTime,
	}
}

//nolint:gochecknoglobals
var settingsValues = map[string]*string{
	"GOOS":         &goOS,
	"GOARCH":       &goArch,
	"vcs.revision": &vcsRevision,
	"vcs.modified": &vcsModified,
	"build-time":   &buildTime,
	"version":      &version,
}

func parseBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	goVersion = info.GoVersion

	if version == unknown {
		version = info.Main.Version
	}

	for _, s := range info.Settings {
		if ptr, ok := settingsValues[s.Key]; ok && *ptr == unknown {
			*ptr = s.Value
		}
	}
}
