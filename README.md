# env [![Build Status](https://drone.io/github.com/caarlos0/env/status.png)](https://drone.io/github.com/caarlos0/env/latest) [![Coverage Status](https://coveralls.io/repos/caarlos0/env/badge.svg?branch=master&service=github)](https://coveralls.io/github/caarlos0/env?branch=master)

A KISS way to deal with environment variables in Go.

## Why

At first, it was boring for me to write down an entire function just to
get some var from the environment and default to another in case it's missing.

For that manner, I wrote the `GetOr` function.

Then, I got pissed about writing `os.Getenv`, `os.Setenv`, `os.Unsetenv`...
it kind of make more sense to me write it as `env.Get`, `env.Set`, `env.Unset`.
So I did.

Then I got a better idea: to use struct tags to do that work for me.

## Example

The most basic example (check the `examples` folder):

```go
package main

import (
	"fmt"

	"github.com/caarlos0/env"
)

type config struct {
	Home         string `env:"HOME"`
	Port         int    `env:"PORT"`
	IsProduction bool   `env:"PRODUCTION"`
}

func main() {
	cfg := config{}
	env.Parse(&cfg)
	fmt.Println(cfg)
}
```

You can run it like this:

```sh
$ PORT=8080 PRODUCTION=true go run examples/first.go
```
