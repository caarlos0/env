package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Some        string `env:"somevar"`
	Other       bool   `env:"othervar"`
	Port        int    `env:"PORT"`
	NotAnEnv    string
	DatabaseURL string   `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`
	Strings     []string `env:"STRINGS"`
	SepStrings  []string `env:"SEPSTRINGS" envSeparator:":"`
	Numbers     []int    `env:"NUMBERS"`
	Bools       []bool   `env:"BOOLS"`
}

func TestParsesEnv(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("STRINGS", "string1,string2,string3")
	os.Setenv("SEPSTRINGS", "string1:string2:string3")
	os.Setenv("NUMBERS", "1,2,3,4")
	os.Setenv("BOOLS", "t,TRUE,0,1")

	defer os.Setenv("somevar", "")
	defer os.Setenv("othervar", "")
	defer os.Setenv("PORT", "")
	defer os.Setenv("STRINGS", "")
	defer os.Setenv("SEPSTRINGS", "")
	defer os.Setenv("NUMBERS", "")
	defer os.Setenv("BOOLS", "")

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.Strings)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.SepStrings)
	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Numbers)
	assert.Equal(t, []bool{true, true, false, true}, cfg.Bools)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "", cfg.Some)
	assert.Equal(t, false, cfg.Other)
	assert.Equal(t, 0, cfg.Port)
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
	defer os.Setenv("othervar", "")

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("PORT", "should-be-an-int")
	defer os.Setenv("PORT", "")

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
		WontWork int64 `env:"BLAH"`
	}
	os.Setenv("BLAH", "10")
	cfg := config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Setenv("WONTWORK", "")

	cfg := &config{}
	assert.Error(t, env.Parse(cfg))
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	cfg := &config{}
	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Setenv("WONTWORK", "")

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
