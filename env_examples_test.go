package env

import (
	"fmt"
	"log"
	"os"
)

func ExampleParse() {
	type Config struct {
		Home string `env:"HOME"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", cfg)
	// Output:  {Home:/tmp/fakehome}
}

func ExampleParseAs() {
	type Config struct {
		Home string `env:"HOME"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg, err := ParseAs[Config]()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", cfg)
	// Output:  {Home:/tmp/fakehome}
}

func ExampleParse_required() {
	type Config struct {
		Nope string `env:"NOPE,required"`
	}
	var cfg Config
	if err := Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: env: required environment variable "NOPE" is not set
	// {Nope:}
}

func ExampleParse_notEmpty() {
	type Config struct {
		Nope string `env:"NOPE,notEmpty"`
	}
	os.Setenv("NOPE", "")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: env: environment variable "NOPE" should not be empty
	// {Nope:}
}

func ExampleParse_unset() {
	type Config struct {
		Secret string `env:"SECRET,unset"`
	}
	os.Setenv("SECRET", "1234")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v - %s", cfg, os.Getenv("SECRET"))
	// Output: {Secret:1234} -
}

func ExampleParse_expand() {
	type Config struct {
		Expand1 string `env:"EXPAND_1,expand"`
		Expand2 string `env:"EXPAND_2,expand" envDefault:"ABC_${EXPAND_1}"`
	}
	os.Setenv("EXPANDING", "HI")
	os.Setenv("EXPAND_1", "HELLO_${EXPANDING}")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Expand1:HELLO_HI Expand2:ABC_HELLO_HI}
}

func ExampleParse_init() {
	type Inner struct {
		A string `env:"OLA" envDefault:"HI"`
	}
	type Config struct {
		NilInner  *Inner
		InitInner *Inner `env:",init"`
	}
	var cfg Config
	if err := Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Print(cfg.NilInner, cfg.InitInner)
	// Output: <nil> &{HI}
}

func ExampleParse_setDefaults() {
	type Config struct {
		Foo string `env:"FOO"`
		Bar string `env:"BAR" envDefault:"bar"`
	}
	cfg := Config{
		Foo: "foo",
	}
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Foo:foo Bar:bar}
}

func ExampleParse_onSet() {
	type config struct {
		Home         string `env:"HOME,required"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		NoEnvTag     bool
		Inner        struct{} `envPrefix:"INNER_"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg config
	if err := ParseWithOptions(&cfg, Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			fmt.Printf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: Set HOME to /tmp/fakehome (default? false)
	// Set PORT to 3000 (default? true)
	// Set PRODUCTION to  (default? false)
	// {Home:/tmp/fakehome Port:3000 IsProduction:false NoEnvTag:false Inner:{}}
}
