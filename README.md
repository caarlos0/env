# env [![Build Status](https://img.shields.io/circleci/project/caarlos0/env/master.svg) [![](https://godoc.org/github.com/caarlos0/env?status.svg)](http://godoc.org/github.com/caarlos0/env) [![Stories in Ready](https://badge.waffle.io/caarlos0/env.svg?label=ready&title=Ready)](http://waffle.io/caarlos0/env)

A KISS way to deal with environment variables in Go.

## Why

At first, it was boring for me to write down an entire function just to
get some `var` from the environment and default to another in case it's missing.

For that manner, I wrote a `GetOr` function in the
[go-idioms](https://github.com/caarlos0/go-idioms) project.

Then, I got pissed about writing `os.Getenv`, `os.Setenv`, `os.Unsetenv`...
it kind of make more sense to me write it as `env.Get`, `env.Set`, `env.Unset`.
So I did.

Then I got a better idea: to use `struct` tags to do all that work for me.

## Example

A very basic example (check the `examples` folder):

```go
package main

import (
	"fmt"
	"os"

	"gopkg.in/caarlos0/env.v1"
)

type config struct {
	Home         string `env:"HOME"`
	Port         int    `env:"PORT" envDefault:"3000"`
	IsProduction bool   `env:"PRODUCTION"`
}

func main() {
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	env.Parse(&cfg)
	fmt.Println(cfg)
}
```

You can run it like this:

```sh
$ PRODUCTION=true go run examples/first.go
{/tmp/fakehome 3000 true}
```

## Supported types and defaults

Currently we only support `string`, `bool` and `int`.

If you set the `envDefault` tag for something, this value will be used in the
case of absence of it in the environment. If you don't do that AND the
environment variable is also not set, the zero-value
of the type will be used: empty for `string`s, `false` for `bool`s
and `0` for `int`s.

