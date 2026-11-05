---
paths:
  - "**/*.go"
---

<go>

- Outside tests and adapter packages, import only the standard library and golib packages. A third-party dependency gets
  its own adapter package, such as `logger/zapadapter` for zap, so that an importer of one package does not build the
  dependencies of another.
- Create and wrap errors with `dev.gaijin.team/go/golib/e`, and put each dynamic value in a field:
  `e.NewFrom("read body", err, fields.F("path", path))`. `fields`, which `e` imports, uses the standard library.

</go>
