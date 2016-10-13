# env [![Build Status](https://travis-ci.org/caarlos0/env.svg?branch=master)](https://travis-ci.org/caarlos0/env) [![Coverage Status](https://coveralls.io/repos/caarlos0/env/badge.svg?branch=master&service=github)](https://coveralls.io/github/caarlos0/env?branch=master) [![](https://godoc.org/github.com/caarlos0/env?status.svg)](http://godoc.org/github.com/caarlos0/env) [![](http://goreportcard.com/badge/caarlos0/env)](http://goreportcard.com/report/caarlos0/env)

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
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
}

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)
}
```

You can run it like this:

```sh
$ PRODUCTION=true HOSTS="host1:host2:host3" DURATION=1s go run examples/first.go
{Home:/your/home Port:3000 IsProduction:true Hosts:[host1 host2 host3] Duration:1s}
```

## Supported types and defaults

The library has support for the following types:

* `string`
* `int`
* `bool`
* `float32`
* `float64`
* `[]string`
* `[]int`
* `[]bool`
* `[]float32`
* `[]float64`

If you set the `envDefault` tag for something, this value will be used in the
case of absence of it in the environment. If you don't do that AND the
environment variable is also not set, the zero-value
of the type will be used: empty for `string`s, `false` for `bool`s
and `0` for `int`s.

By default, slice types will split the environment value on `,`; you can change this behavior by setting the `envSeparator` tag.

## Required fields

The `env` tag option `required` (e.g., `env:"tagKey,required"`) can be added
to ensure that some environment variable is set.  In the example above,
an error is returned if the `config` struct is changed to:


```go
type config struct {
    Home         string   `env:"HOME"`
    Port         int      `env:"PORT" envDefault:"3000"`
    IsProduction bool     `env:"PRODUCTION"`
    Hosts        []string `env:"HOSTS" envSeparator:":"`
    SecretKey    string   `env:"SECRET_KEY,required"`
}
```
