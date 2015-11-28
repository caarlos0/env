package main

import (
	"fmt"

	"github.com/dselans/env"
)

type config struct {
	Home         string   `env:"HOME"`
	Port         int      `env:"PORT" default:"3000"`
	IsProduction bool     `env:"PRODUCTION"`
	Hosts        []string `env:"HOSTS" envSeparator:":"`
}

func main() {
	cfg := config{}
	env.Parse(&cfg)
	fmt.Println(cfg)
}
