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

	time   = unknown
	vcsTag = unknown
)

// Info represents the build information for a Go binary.
// It contains details about the Go environment used for building,
// version control system state, and custom build-time information.
type Info struct {
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

	//
	// User-defined information, set at build time.
	//

	// BuildTime is the build time.
	BuildTime string `json:"build-time"`
	// BuildTag is the vcs tag of the commit, binary built from.
	BuildTag string `json:"build-tag"`
}

// String returns a human-readable representation of build information.
// The format includes Go version, OS, architecture, VCS revision, and build details.
// If the repository was modified (dirty), this is indicated in the output.
// If a build tag is present, it's included in square brackets. after the revision.
func (i Info) String() string {
	b := &strings.Builder{}

	b.WriteString("built with ")
	b.WriteString(i.GoVersion)
	b.WriteString(" for ")
	b.WriteString(i.GoOS)
	b.WriteString(" ")
	b.WriteString(i.GoArch)
	b.WriteString(" from ")
	b.WriteString(i.VCSRevision)

	if i.BuildTag != unknown {
		b.WriteString(" [")
		b.WriteString(i.BuildTag)
		b.WriteRune(']')
	}

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
		fields.F("go-version", i.GoVersion),
		fields.F("go-os", i.GoOS),
		fields.F("go-arch", i.GoArch),
		fields.F("vcs-revision", i.VCSRevision),
		fields.F("vcs-modified", i.VCSModified),
		fields.F("build-time", i.BuildTime),
		fields.F("build-tag", i.BuildTag),
	}
}

// GetInfo returns the build information.
// It automatically parses build information from runtime debug data
// the first time it's called, and caches the result for subsequent calls.
func GetInfo() Info {
	if !buildInfoParsed {
		buildInfoParsed = true

		parseBuildInfo()
	}

	return Info{
		GoVersion:   goVersion,
		GoOS:        goOS,
		GoArch:      goArch,
		VCSRevision: vcsRevision,
		VCSModified: vcsModified,
		BuildTime:   time,
		BuildTag:    vcsTag,
	}
}

//nolint:gochecknoglobals
var settingsProcessors = map[string]func(s debug.BuildSetting){
	"GOOS":         func(s debug.BuildSetting) { goOS = s.Value },
	"GOARCH":       func(s debug.BuildSetting) { goArch = s.Value },
	"vcs.revision": func(s debug.BuildSetting) { vcsRevision = s.Value },
	"vcs.modified": func(s debug.BuildSetting) { vcsModified = s.Value },
}

func parseBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	goVersion = info.GoVersion

	for _, s := range info.Settings {
		if f, ok := settingsProcessors[s.Key]; ok {
			f(s)
		}
	}
}
