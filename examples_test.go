package env

import (
	"fmt"
	"os"
	"reflect"
)

func ExampleParseAs() {
	type config struct {
		Home         string         `env:"HOME,required"`
		Port         int            `env:"PORT" envDefault:"3000"`
		IsProduction bool           `env:"PRODUCTION"`
		TempFolder   string         `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/.tmp"`
		StringInts   map[string]int `env:"MAP_STRING_INT" envDefault:"k1:1,k2:2"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg, err := ParseAs[config]()
	if err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output:  {Home:/tmp/fakehome Port:3000 IsProduction:false TempFolder:/tmp/fakehome/.tmp StringInts:map[k1:1 k2:2]}
}

func ExampleParse() {
	type inner struct {
		Foo string `env:"FOO" envDefault:"foobar"`
	}
	type config struct {
		Home         string         `env:"HOME,required"`
		Port         int            `env:"PORT" envDefault:"3000"`
		IsProduction bool           `env:"PRODUCTION"`
		TempFolder   string         `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/.tmp"`
		StringInts   map[string]int `env:"MAP_STRING_INT" envDefault:"k1:1,k2:2"`
		Inner        inner
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg config
	if err := Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output:  {Home:/tmp/fakehome Port:3000 IsProduction:false TempFolder:/tmp/fakehome/.tmp StringInts:map[k1:1 k2:2] Inner:{Foo:foobar}}
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

func ExampleParse_defaults() {
	type config struct {
		A string `env:"FOO" envDefault:"foo"`
		B string `env:"FOO"`
	}

	// env FOO is not set

	cfg := config{
		A: "A",
		B: "B",
	}
	if err := Parse(&cfg); err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {A:foo B:B}
}

func ExampleParseWithOptions() {
	type thing struct {
		desc string
	}

	type conf struct {
		Thing thing `env:"THING"`
	}

	os.Setenv("THING", "my thing")

	c := conf{}

	err := ParseWithOptions(&c, Options{FuncMap: map[reflect.Type]ParserFunc{
		reflect.TypeOf(thing{}): func(v string) (interface{}, error) {
			return thing{desc: v}, nil
		},
	}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Thing.desc)
	// Output:
	// my thing
}

func ExampleParse_expandVars() {
	type config struct {
		Host    string `env:"HOST" envDefault:"localhost"`
		Port    int    `env:"PORT" envDefault:"3000"`
		Address string `env:"ADDRESS,expand" envDefault:"$HOST:${PORT}"`
	}

	cfg := config{}
	if err := Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)
	// Output: {Host:localhost Port:3000 Address:localhost:3000}
}

func ExampleParse_unset() {
	type config struct {
		SecretKey string `env:"SECRET_KEY,unset"`
	}
	os.Setenv("SECRET_KEY", "asd")
	cfg := config{}
	if err := Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)
	fmt.Println(os.Getenv("SECRET_KEY"))
	// Output: {SecretKey:asd}
}
