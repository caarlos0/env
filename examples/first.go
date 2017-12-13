package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	ExampleFoo   Foo           `env:"EXAMPLE_FOO"`
}

type Foo struct {
	Name string
}

func main() {
	cfg := config{}

	// Parse for built-in types
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Unable to parse envs: ", err)
	}

	// OR w/ a custom parser for `Foo`
	//
	// if err := env.ParseWithFuncs(&cfg, env.CustomParsers{
	// 	reflect.TypeOf(Foo{}): fooParser,
	// }); err != nil {
	// 	log.Fatal("Unable to parse envs: ", err)
	// }

	fmt.Printf("%+v\n", cfg)
}

func fooParser(value string) (interface{}, error) {
	return Foo{
		Name: strings.ToUpper(value),
	}, nil
}
