package build_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/build"
	"dev.gaijin.team/go/golib/fields"
)

func TestInfo_String(t *testing.T) { //nolint:tparallel
	t.Parallel()

	tests := []struct {
		name     string
		info     build.Info
		expected string
	}{
		{
			name: "all fields populated",
			info: build.Info{
				GoVersion:   "go1.17.5",
				GoOS:        "linux",
				GoArch:      "amd64",
				VCSRevision: "abcdef123456",
				VCSModified: "false",
				BuildTime:   "2023-07-15T12:34:56Z",
				BuildTag:    "v1.0.0",
			},
			expected: "built with go1.17.5 for linux amd64 from abcdef123456 [v1.0.0] build time 2023-07-15T12:34:56Z",
		},
		{
			name: "with unknown build tag",
			info: build.Info{
				GoVersion:   "go1.18.0",
				GoOS:        "darwin",
				GoArch:      "arm64",
				VCSRevision: "987654abcdef",
				VCSModified: "unknown",
				BuildTime:   "2023-08-20T10:11:12Z",
				BuildTag:    "unknown",
			},
			expected: "built with go1.18.0 for darwin arm64 from 987654abcdef build time 2023-08-20T10:11:12Z",
		},
		{
			name: "with dirty repo",
			info: build.Info{
				GoVersion:   "go1.19.0",
				GoOS:        "windows",
				GoArch:      "amd64",
				VCSRevision: "1a2b3c4d5e6f",
				VCSModified: "true",
				BuildTime:   "2023-09-25T15:16:17Z",
				BuildTag:    "v2.0.0",
			},
			//revive:disable-next-line:line-length-limit
			expected: "built with go1.19.0 for windows amd64 from 1a2b3c4d5e6f [v2.0.0] (dirty) build time 2023-09-25T15:16:17Z",
		},
		{
			name: "unknown build tag and dirty repo",
			info: build.Info{
				GoVersion:   "go1.20.0",
				GoOS:        "freebsd",
				GoArch:      "amd64",
				VCSRevision: "a1b2c3d4e5f6",
				VCSModified: "true",
				BuildTime:   "2023-10-30T20:21:22Z",
				BuildTag:    "unknown",
			},
			expected: "built with go1.20.0 for freebsd amd64 from a1b2c3d4e5f6 (dirty) build time 2023-10-30T20:21:22Z",
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.info.String())
		})
	}
}

func TestInfo_ToFields(t *testing.T) {
	t.Parallel()

	info := build.Info{
		GoVersion:   "go1.17.5",
		GoOS:        "linux",
		GoArch:      "amd64",
		VCSRevision: "abcdef123456",
		VCSModified: "false",
		BuildTime:   "2023-07-15T12:34:56Z",
		BuildTag:    "v1.0.0",
	}

	expectedFields := fields.List{
		fields.F("go-version", "go1.17.5"),
		fields.F("go-os", "linux"),
		fields.F("go-arch", "amd64"),
		fields.F("vcs-revision", "abcdef123456"),
		fields.F("vcs-modified", "false"),
		fields.F("build-time", "2023-07-15T12:34:56Z"),
		fields.F("build-tag", "v1.0.0"),
	}

	resultFields := info.ToFields()

	// Test length
	assert.Equal(t, len(expectedFields), len(resultFields), "fields list length should match")

	// Test each field
	for i := range len(expectedFields) {
		assert.Equal(t, expectedFields[i].K, resultFields[i].K, "Field key should match")
		assert.Equal(t, expectedFields[i].V, resultFields[i].V, "Field value should match")
	}
}
