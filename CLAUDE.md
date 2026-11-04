# golib

Go packages shared by Gaijin Entertainment projects, published as the public module `dev.gaijin.team/go/golib`. Other
modules import each package, so a change is judged by its effect on importers. The repository is public: a committed
file carries no internal hostname, internal service name, or local path.

<placement>

- A logging backend binds to `logger` in its own `logger/<backend>adapter` package.
- A new general-purpose package gets its own top-level directory, named for what it provides.

</placement>

<toolchain>

- The `go` directive in `go.mod` is the older of the two newest Go releases, written as `go 1.N`. The review CI
  (`_tests/tests.py`) tests in a Go image tagged by that directive, and the image exists only for minor versions, so
  `go 1.N.0` fails the test job. `go mod tidy` raises the directive to `go 1.N.0` when a dependency declares that
  version: hold such a dependency at its newest release that declares an older one.
- GitHub Actions test on the two newest Go releases, on Linux, macOS, and Windows (`.github/workflows/ci.yml`).
- `vendor/` is committed. After a change of `go.mod`, run `go mod tidy && go mod vendor`.
- Lint with the golangci-lint version in the `###golangci-lint:version:` comment atop `.golangci.yaml`. The review CI
  runs that version.

</toolchain>

<verification>

- `go test ./...`
- `golangci-lint run ./...`

</verification>

<git-commit-config>
<extra-instructions>

- Write the subject as `[<scope>] <summary>`. The scope is the path of the changed package, or the area of a change that
  no package owns.

</extra-instructions>
</git-commit-config>
