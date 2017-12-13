parsers
=======
This directory contains pre-built, custom parsers that can be used with `env.ParseWithFuncs`
to facilitate the parsing of envs that are not basic types.

Example Usage:

```golang
package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/caarlos0/env"
	"github.com/caarlos0/env/parsers"
)

type config struct {
	ExampleURL url.URL `env:"EXAMPLE_URL" envDefault:"https://google.com"`
}

func main() {
	cfg := config{}

	if err := env.ParseWithFuncs(&cfg, env.CustomParsers{
		parsers.URLType: parsers.URLFunc,
	}); err != nil {
		log.Fatal("Unable to parse envs: ", err)
	}

	fmt.Printf("Scheme: %v Host: %v\n", cfg.ExampleURL.Scheme, cfg.ExampleURL.Host)
}
```