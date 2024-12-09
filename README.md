<p align="center">
  <img alt="GoReleaser Logo" src="https://becker.software/env.png" height="140" />
  <p align="center">A simple, zero-dependencies library to parse environment variables into structs.</p>
</p>

###### Installation

```bash
go get github.com/caarlos0/env/v11
```

###### Getting started

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

You can see the full documentation and list of examples at [pkg.go.dev](https://pkg.go.dev/github.com/caarlos0/env/v11).

---

## Used and supported by

<p>
  <a href="https://encore.dev">
    <img src="https://user-images.githubusercontent.com/78424526/214602214-52e0483a-b5fc-4d4c-b03e-0b7b23e012df.svg" width="120px" alt="encore icon" />
  </a>
  <br/>
  <br/>
  <b>Encore – the platform for building Go-based cloud backends.</b>
  <br/>
</p>

## Usage

### Caveats

> [!CAUTION]
>
> _Unexported fields_ will be **ignored** by `env`.
> This is by design and will not change.

### Functions

- `Parse`: parse the current environment into a type
- `ParseAs`: parse the current environment into a type using generics
- `ParseWithOptions`: parse the current environment into a type with custom options
- `ParseAsWithOptions`: parse the current environment into a type with custom options and using generics
- `Must`: can be used to wrap `Parse.*` calls to panic on error
- `GetFieldParams`: get the `env` parsed options for a type
- `GetFieldParamsWithOptions`: get the `env` parsed options for a type with custom options

### Supported types

Out of the box all built-in types are supported, plus a few others that are commonly used.

Complete list:

- `bool`
- `float32`
- `float64`
- `int16`
- `int32`
- `int64`
- `int8`
- `int`
- `string`
- `uint16`
- `uint32`
- `uint64`
- `uint8`
- `uint`
- `time.Duration`
- `time.Location`
- `encoding.TextUnmarshaler`
- `url.URL`

Pointers, slices and slices of pointers, and maps of those types are also supported.

You may also add custom parsers for your types.

### Tags

The following tags are provided:

- `env`: sets the environment variable name and optionally takes the tag options described below
- `envDefault`: sets the default value for the field
- `envPrefix`: can be used in a field that is a complex type to set a prefix to all environment variables used in it
- `envSeparator`: sets the character to be used to separate items in slices and maps (default: `,`)
- `envKeyValSeparator`: sets the character to be used to separate keys and their values in maps (default: `:`)

### `env` tag options

Here are all the options available for the `env` tag:

- `,expand`: expands environment variables, e.g. `FOO_${BAR}`
- `,file`: instructs that the content of the variable is a path to a file that should be read
- `,init`: initialize nil pointers
- `,notEmpty`: make the field errors if the environment variable is empty
- `,required`: make the field errors if the environment variable is not set
- `,unset`: unset the environment variable after use

### Parse Options

There are a few options available in the functions that end with `WithOptions`:

- `Environment`: keys and values to be used instead of `os.Environ()`
- `TagName`: specifies another tag name to use rather than the default `env`
- `PrefixTagName`: specifies another prefix tag name to use rather than the default `envPrefix`
- `DefaultValueTagName`: specifies another default tag name to use rather than the default `envDefault`
- `RequiredIfNoDef`: set all `env` fields as required if they do not declare `envDefault`
- `OnSet`: allows to hook into the `env` parsing and do something when a value is set
- `Prefix`: prefix to be used in all environment variables
- `UseFieldNameByDefault`: defines whether or not `env` should use the field name by default if the `env` key is missing
- `FuncMap`: custom parse functions for custom types

### Documentation and examples

Examples are live in [pkg.go.dev](https://pkg.go.dev/github.com/caarlos0/env/v11),
and also in the [example test file](./example_test.go).

## Current state

`env` is considered feature-complete.

I do not intent to add any new features unless they really make sense, and are
requested by many people.

Eventual bug fixes will keep being merged.

## Badges

[![Release](https://img.shields.io/github/release/caarlos0/env.svg?style=for-the-badge)](https://github.com/goreleaser/goreleaser/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE.md)
[![Build status](https://img.shields.io/github/actions/workflow/status/caarlos0/env/build.yml?style=for-the-badge&branch=main)](https://github.com/caarlos0/env/actions?workflow=build)
[![Codecov branch](https://img.shields.io/codecov/c/github/caarlos0/env/main.svg?style=for-the-badge)](https://codecov.io/gh/caarlos0/env)
[![Go docs](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/caarlos0/env/v11)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## Related projects

- [envdoc](https://github.com/g4s8/envdoc) - generate documentation for environment variables from `env` tags

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/env.svg)](https://starchart.cc/caarlos0/env)
