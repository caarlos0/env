package main

import (
	"fmt"

	"github.com/caarlos0/env"
)

func ExampleParse() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" default:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	cfg := config{}
	env.Parse(&cfg)
	fmt.Println(cfg)
}
