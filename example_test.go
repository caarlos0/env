package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

// Basic package usage example.
func Example() {
	type Config struct {
		Foo string `env:"FOO"`
	}

	os.Setenv("FOO", "bar")

	// parse:
	var cfg1 Config
	_ = Parse(&cfg1)

	// parse with generics:
	cfg2, _ := ParseAs[Config]()

	fmt.Print(cfg1.Foo, cfg2.Foo)
	// Output: barbar
}

// Parse the environment into a struct.
func ExampleParse() {
	type Config struct {
		Home string `env:"HOME"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output:  {Home:/tmp/fakehome}
}

// Parse the environment into a struct using generics.
func ExampleParseAs() {
	type Config struct {
		Home string `env:"HOME"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg, err := ParseAs[Config]()
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: env: required environment variable "NOPE" is not set
	// {Nope:}
}

// While `required` demands the environment variable to be set, it doesn't check
// its value. If you want to make sure the environment is set and not empty, you
// need to use the `notEmpty` tag option instead (`env:"SOME_ENV,notEmpty"`).
func ExampleParse_notEmpty() {
	type Config struct {
		Nope string `env:"NOPE,notEmpty"`
	}
	os.Setenv("NOPE", "")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: env: environment variable "NOPE" should not be empty
	// {Nope:}
}

// The `env` tag option `unset` (e.g., `env:"tagKey,unset"`) can be added
// to ensure that some environment variable is unset after reading it.
func ExampleParse_unset() {
	type Config struct {
		Secret string `env:"SECRET,unset"`
	}
	os.Setenv("SECRET", "1234")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v - %s", cfg, os.Getenv("SECRET"))
	// Output: {Secret:1234} -
}

// You can use `envSeparator` to define which character should be used to
// separate array items in a string.
// Similarly, you can use `envKeyValSeparator` to define which character should
// be used to separate a key from a value in a map.
// The defaults are `,` and `:`, respectively.
func ExampleParse_separator() {
	type Config struct {
		Map map[string]string `env:"CUSTOM_MAP" envSeparator:"-" envKeyValSeparator:"|"`
	}
	os.Setenv("CUSTOM_MAP", "k1|v1-k2|v2")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Map:map[k1:v1 k2:v2]}
}

// If you set the `expand` option, environment variables (either in `${var}` or
// `$var` format) in the string will be replaced according with the actual
// value of the variable. For example:
func ExampleParse_expand() {
	type Config struct {
		Expand1 string `env:"EXPAND_1,expand"`
		Expand2 string `env:"EXPAND_2,expand" envDefault:"ABC_${EXPAND_1}"`
	}
	os.Setenv("EXPANDING", "HI")
	os.Setenv("EXPAND_1", "HELLO_${EXPANDING}")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Expand1:HELLO_HI Expand2:ABC_HELLO_HI}
}

// You can automatically initialize `nil` pointers regardless of if a variable
// is set for them or not.
// This behavior can be enabled by using the `init` tag option.
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
		fmt.Println(err)
	}
	fmt.Print(cfg.NilInner, cfg.InitInner)
	// Output: <nil> &{HI}
}

// You can define the default value for a field by either using the
// `envDefault` tag, or when initializing the `struct`.
//
// Default values defined as `struct` tags will overwrite existing values
// during `Parse`.
func ExampleParse_setDefaults() {
	type Config struct {
		Foo string `env:"DEF_FOO"`
		Bar string `env:"DEF_BAR" envDefault:"bar"`
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

// You might want to listen to value sets and, for example, log something or do
// some other kind of logic.
func ExampleParseWithOptions_onSet() {
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

// By default, env supports anything that implements the `TextUnmarshaler`
// interface, which includes `time.Time`.
//
// The upside is that depending on the format you need, you don't need to change
// anything.
//
// The downside is that if you do need time in another format, you'll need to
// create your own type and implement `TextUnmarshaler`.
func ExampleParse_customTimeFormat() {
	// type MyTime time.Time
	//
	// func (t *MyTime) UnmarshalText(text []byte) error {
	// 	tt, err := time.Parse("2006-01-02", string(text))
	// 	*t = MyTime(tt)
	// 	return err
	// }

	type Config struct {
		SomeTime MyTime `env:"SOME_TIME"`
	}
	os.Setenv("SOME_TIME", "2021-05-06")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Print(cfg.SomeTime)
	// Output: {0 63755856000 <nil>}
}

// Parse using extra options.
func ExampleParseWithOptions_customTypes() {
	type Thing struct {
		desc string
	}

	type Config struct {
		Thing Thing `env:"THING"`
	}

	os.Setenv("THING", "my thing")

	c := Config{}
	err := ParseWithOptions(&c, Options{
		FuncMap: map[reflect.Type]ParserFunc{
			reflect.TypeOf(Thing{}): func(v string) (interface{}, error) {
				return Thing{desc: v}, nil
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(c.Thing.desc)
	// Output: my thing
}

// Make all fields required by default.
func ExampleParseWithOptions_allFieldsRequired() {
	type Config struct {
		Username string `env:"EX_USERNAME" envDefault:"admin"`
		Password string `env:"EX_PASSWORD"`
	}

	var cfg Config
	if err := ParseWithOptions(&cfg, Options{
		RequiredIfNoDef: true,
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cfg)
	// Output: env: required environment variable "EX_PASSWORD" is not set
	// {Username:admin Password:}
}

// Set a custom environment.
// By default, `os.Environ()` is used.
func ExampleParseWithOptions_setEnv() {
	type Config struct {
		Username string `env:"EX_USERNAME" envDefault:"admin"`
		Password string `env:"EX_PASSWORD"`
	}

	var cfg Config
	if err := ParseWithOptions(&cfg, Options{
		Environment: map[string]string{
			"EX_USERNAME": "john",
			"EX_PASSWORD": "cena",
		},
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cfg)
	// Output: {Username:john Password:cena}
}

// Handling slices of complex types.
func ExampleParse_complexSlices() {
	type Test struct {
		Str string `env:"STR"`
		Num int    `env:"NUM"`
	}
	type Config struct {
		Foo []Test `envPrefix:"FOO"`
	}

	os.Setenv("FOO_0_STR", "a")
	os.Setenv("FOO_0_NUM", "1")
	os.Setenv("FOO_1_STR", "b")
	os.Setenv("FOO_1_NUM", "2")

	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cfg)
	// Output: {Foo:[{Str:a Num:1} {Str:b Num:2}]}
}

// Setting prefixes for inner types.
func ExampleParse_prefix() {
	type Inner struct {
		Foo string `env:"FOO,required"`
	}
	type Config struct {
		A Inner `envPrefix:"A_"`
		B Inner `envPrefix:"B_"`
	}
	os.Setenv("A_FOO", "a")
	os.Setenv("B_FOO", "b")
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {A:{Foo:a} B:{Foo:b}}
}

// Setting prefixes for the entire config.
func ExampleParseWithOptions_prefix() {
	type Config struct {
		Foo string `env:"FOO"`
	}
	os.Setenv("MY_APP_FOO", "a")
	var cfg Config
	if err := ParseWithOptions(&cfg, Options{
		Prefix: "MY_APP_",
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Foo:a}
}

// Use a different tag name than `env` and `envDefault`.
func ExampleParseWithOptions_tagName() {
	type Config struct {
		Home string `json:"HOME"`
		Page string `json:"PAGE" def:"world"`
	}
	os.Setenv("HOME", "hello")
	var cfg Config
	if err := ParseWithOptions(&cfg, Options{
		TagName:             "json",
		DefaultValueTagName: "def",
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Home:hello Page:world}
}

// If you don't want to set the `env` tag on every field, you can use the
// `UseFieldNameByDefault` option.
//
// It will use the field name to define the environment variable name.
// So, `Foo` becomes `FOO`, `FooBar` becomes `FOO_BAR`, and so on.
func ExampleParseWithOptions_useFieldName() {
	type Config struct {
		Foo string
	}
	os.Setenv("FOO", "bar")
	var cfg Config
	if err := ParseWithOptions(&cfg, Options{
		UseFieldNameByDefault: true,
	}); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Foo:bar}
}

// The `env` tag option `file` (e.g., `env:"tagKey,file"`) can be added in
// order to indicate that the value of the variable shall be loaded from a
// file.
//
// The path of that file is given by the environment variable associated with
// it.
func ExampleParse_fromFile() {
	f, _ := os.CreateTemp("", "")
	_, _ = f.WriteString("super secret")
	_ = f.Close()

	type Config struct {
		Secret string `env:"SECRET,file"`
	}
	os.Setenv("SECRET", f.Name())
	var cfg Config
	if err := Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", cfg)
	// Output: {Secret:super secret}
}

// TODO: envSeperator
//

func ExampleParse_errorHandling() {
	type Config struct {
		Username string `env:"EX_ERR_USERNAME" envDefault:"admin"`
		Password string `env:"EX_ERR_PASSWORD,notEmpty"`
	}

	var cfg Config
	if err := Parse(&cfg); err != nil {
		if errors.Is(err, EmptyVarError{}) {
			fmt.Println("oopsie")
		}
		aggErr := AggregateError{}
		if ok := errors.As(err, &aggErr); ok {
			for _, er := range aggErr.Errors {
				switch v := er.(type) {
				// Handle the error types you need:
				// ParseError
				// NotStructPtrError
				// NoParserError
				// NoSupportedTagOptionError
				// EnvVarIsNotSetError
				// EmptyEnvVarError
				// LoadFileContentError
				// ParseValueError
				case EmptyVarError:
					fmt.Println("daisy")
				default:
					fmt.Printf("Unknown error type %v", v)
				}
			}
		}
	}

	fmt.Printf("%+v", cfg)
	// Output: oopsie
	// daisy
	// {Username:admin Password:}
}

// You can avoid setting defaults for non zero values
// This could be useful for loading data from config file first
// and then filling the rest from env
func Example_setDefaultsForZeroValuesOnly() {
	type Config struct {
		Username string `env:"USERNAME" envDefault:"admin"`
		Password string `env:"PASSWORD" envDefault:"qwerty"`
	}

	cfg := Config{
		Username: "root",
	}

	if err := ParseWithOptions(&cfg, Options{
		Environment:                  map[string]string{},
		SetDefaultsForZeroValuesOnly: true,
	}); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", cfg)
	// Without SetDefaultsForZeroValuesOnly, the username would have been 'admin'.
	// Output: {Username:root Password:qwerty}
}
