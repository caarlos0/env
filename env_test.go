package env_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Some        string `env:"somevar"`
	Other       bool   `env:"othervar"`
	Port        int    `env:"PORT"`
	UintVal     uint   `env:"UINTVAL"`
	NotAnEnv    string
	DatabaseURL string        `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`
	Strings     []string      `env:"STRINGS"`
	SepStrings  []string      `env:"SEPSTRINGS" envSeparator:":"`
	Numbers     []int         `env:"NUMBERS"`
	Numbers64   []int64       `env:"NUMBERS64"`
	Bools       []bool        `env:"BOOLS"`
	Duration    time.Duration `env:"DURATION"`
	Float32     float32       `env:"FLOAT32"`
	Float64     float64       `env:"FLOAT64"`
	Float32s    []float32     `env:"FLOAT32S"`
	Float64s    []float64     `env:"FLOAT64S"`
}

type ParentStruct struct {
	InnerStruct *InnerStruct
	unexported  *InnerStruct
	Ignored     *http.Client
}

type InnerStruct struct {
	Inner string `env:"innervar"`
}

func TestParsesEnv(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("STRINGS", "string1,string2,string3")
	os.Setenv("SEPSTRINGS", "string1:string2:string3")
	os.Setenv("NUMBERS", "1,2,3,4")
	os.Setenv("NUMBERS64", "1,2,2147483640,-2147483640")
	os.Setenv("BOOLS", "t,TRUE,0,1")
	os.Setenv("DURATION", "1s")
	os.Setenv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
	os.Setenv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
	os.Setenv("FLOAT32S", "1.0,2.0,3.0")
	os.Setenv("FLOAT64S", "1.0,2.0,3.0")
	os.Setenv("UINTVAL", "44")

	defer os.Clearenv()

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, uint(44), cfg.UintVal)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.Strings)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.SepStrings)
	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Numbers)
	assert.Equal(t, []int64{1, 2, 2147483640, -2147483640}, cfg.Numbers64)
	assert.Equal(t, []bool{true, true, false, true}, cfg.Bools)
	d, _ := time.ParseDuration("1s")
	assert.Equal(t, d, cfg.Duration)
	f32 := float32(3.40282346638528859811704183484516925440e+38)
	assert.Equal(t, f32, cfg.Float32)
	f64 := float64(1.797693134862315708145274237317043567981e+308)
	assert.Equal(t, f64, cfg.Float64)
	assert.Equal(t, []float32{float32(1.0), float32(2.0), float32(3.0)}, cfg.Float32s)
	assert.Equal(t, []float64{float64(1.0), float64(2.0), float64(3.0)}, cfg.Float64s)
}

func TestParsesEnvInner(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "someinnervalue", cfg.InnerStruct.Inner)
}

func TestParsesEnvInnerNil(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{}
	assert.NoError(t, env.Parse(&cfg))
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "", cfg.Some)
	assert.Equal(t, false, cfg.Other)
	assert.Equal(t, 0, cfg.Port)
	assert.Equal(t, uint(0), cfg.UintVal)
	assert.Equal(t, 0, len(cfg.Strings))
	assert.Equal(t, 0, len(cfg.SepStrings))
	assert.Equal(t, 0, len(cfg.Numbers))
	assert.Equal(t, 0, len(cfg.Bools))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.Error(t, env.Parse(&thisShouldBreak))
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("othervar", "should-be-a-bool")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("PORT", "should-be-an-int")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidUint(t *testing.T) {
	os.Setenv("UINTVAL", "-44")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidBoolsSlice(t *testing.T) {
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}

	os.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestInvalidDuration(t *testing.T) {
	os.Setenv("DURATION", "should-be-a-valid-duration")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestParsesDefaultConfig(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "postgres://localhost:5432/db", cfg.DatabaseURL)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Empty(t, cfg.NotAnEnv)
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	os.Setenv("BLAH", "a")
	cfg := config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Clearenv()

	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	cfg := &config{}
	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Clearenv()

	assert.Error(t, env.Parse(cfg))
}

func TestNoErrorRequiredSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}

	os.Setenv("IS_REQUIRED", "val")
	defer os.Clearenv()
	assert.NoError(t, env.Parse(cfg))
	assert.Equal(t, "val", cfg.IsRequired)
}

func TestErrorRequiredNotSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED,required"`
	}

	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestCustomParser(t *testing.T) {
	type foo struct {
		name string
	}

	type config struct {
		Var foo `env:"VAR"`
	}

	os.Setenv("VAR", "test")

	customParserFunc := func(v string) (interface{}, error) {
		return foo{name: v}, nil
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(foo{}): customParserFunc,
	})

	assert.NoError(t, err)
	assert.Equal(t, cfg.Var.name, "test")
}

func TestParseWithFuncsNoPtr(t *testing.T) {
	type foo struct{}
	err := env.ParseWithFuncs(foo{}, nil)
	assert.Error(t, err)
	assert.Equal(t, err, env.ErrNotAStructPtr)
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	err := env.ParseWithFuncs(&c, nil)
	assert.Error(t, err)
	assert.Equal(t, err, env.ErrNotAStructPtr)
}

func TestCustomParserError(t *testing.T) {
	type foo struct {
		name string
	}

	type config struct {
		Var foo `env:"VAR"`
	}

	os.Setenv("VAR", "test")

	customParserFunc := func(v string) (interface{}, error) {
		return nil, errors.New("something broke")
	}

	cfg := &config{}
	err := env.ParseWithFuncs(cfg, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(foo{}): customParserFunc,
	})

	assert.Empty(t, cfg.Var.name, "Var.name should not be filled out when parse errors")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Custom parser error: something broke")
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}

	os.Setenv("FOO", "foo")

	cfg := &config{}
	err := env.Parse(cfg)

	assert.Error(t, err)
	assert.Equal(t, env.ErrUnsupportedType, err)
}
func TestEmptyOption(t *testing.T) {
	type config struct {
		Var string `env:"VAR,"`
	}

	cfg := &config{}

	os.Setenv("VAR", "val")
	defer os.Clearenv()
	assert.NoError(t, env.Parse(cfg))
	assert.Equal(t, "val", cfg.Var)
}

func TestErrorOptionNotRecognized(t *testing.T) {
	type config struct {
		Var string `env:"VAR,not_supported!"`
	}

	cfg := &config{}
	assert.Error(t, env.Parse(cfg))

}

func ExampleParse() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	env.Parse(&cfg)
	fmt.Println(cfg)
	// Output: {/tmp/fakehome 3000 false}
}

func ExampleParseRequiredField() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		SecretKey    string `env:"SECRET_KEY,required"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	err := env.Parse(&cfg)
	fmt.Println(err)
	// Output: Required environment variable SECRET_KEY is not set
}

func ExampleParseMultipleOptions() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
		SecretKey    string `env:"SECRET_KEY,required,option1"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	err := env.Parse(&cfg)
	fmt.Println(err)
	// Output: Env tag option option1 not supported.
}
