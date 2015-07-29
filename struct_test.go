package env_test

import (
	"fmt"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Some        string `env:"somevar"`
	Other       bool   `env:"othervar"`
	Port        int    `env:"PORT"`
	NotAnEnv    string
	DatabaseURL string `env:"DATABASE_URL" default:"postgres://localhost:5432/db"`
}

func TestParsesEnv(t *testing.T) {
	env.Set("somevar", "somevalue")
	env.Set("othervar", "true")
	env.Set("PORT", "8080")
	defer env.Unset("somevar")
	defer env.Unset("othervar")
	defer env.Unset("PORT")

	cfg := Config{}
	assert.NoError(t, env.Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
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
	env.Set("othervar", "should-be-a-bool")
	defer env.Unset("othervar")

	cfg := Config{}
	assert.Error(t, env.Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	env.Set("PORT", "should-be-an-int")
	defer env.Unset("PORT")

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

func ExampleParse() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" default:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	cfg := config{}
	env.Set("HOME", "/tmp/fakehome")
	env.Parse(&cfg)
	fmt.Println(cfg)
	// Output: {/tmp/fakehome 3000 false}
}
