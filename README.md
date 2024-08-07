<p align="center">
  <img alt="GoReleaser Logo" src="https://becker.software/env.png" height="140" />
  <p align="center">A simple, zero-dependencies library to parse environment variables into structs.</p>
</p>

A simple and zero-dependencies library to parse environment variables into `struct`s.

```go
type config struct {
  Home string `env:"HOME"`
}

// parse
var cfg config
err := env.Parse(&cfg)

// parse with generics
cfg, err := env.ParseAs[config]()
```

You can see the full documentation and list of examples at
[pkg.go.dev](https://pkg.go.dev/github.com/caarlos0/env/v11).

---

## Used and supported by

<p>
  <a href="https://encore.dev">
    <img src="https://user-images.githubusercontent.com/78424526/214602214-52e0483a-b5fc-4d4c-b03e-0b7b23e012df.svg" width="120px" alt="encore icon"></img>
  </a>
  <br/>
  <br/>
  <b>Encore – the platform for building Go-based cloud backends.</b>
  <br/>
</p>

## Caveats

> [!CAUTION]
>
> _Unexported fields_ will be **ignored** by `env`.
> This is by design and will not change.

## Supported types and defaults

Out of the box all built-in types are supported, plus a few others that
are commonly used.

Complete list:

- `string`
- `bool`
- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`
- `float32`
- `float64`
- `time.Duration`
- `encoding.TextUnmarshaler`
- `url.URL`

Pointers, slices and slices of pointers, and maps of those types are also
supported.

## Badges

[![Release](https://img.shields.io/github/release/caarlos0/env.svg?style=for-the-badge)](https://github.com/goreleaser/goreleaser/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE.md)
[![Build status](https://img.shields.io/github/actions/workflow/status/caarlos0/env/build.yml?style=for-the-badge&branch=main)](https://github.com/caarlos0/env/actions?workflow=build)
[![Codecov branch](https://img.shields.io/codecov/c/github/caarlos0/env/main.svg?style=for-the-badge)](https://codecov.io/gh/caarlos0/env)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/caarlos0/env/v11)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## Related projects

- [envdoc](https://github.com/g4s8/envdoc) - generate documentation for environment variables from `env` tags

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/env.svg)](https://starchart.cc/caarlos0/env)
