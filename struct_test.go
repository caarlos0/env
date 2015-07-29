package env

import (
	"testing"

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
	Set("somevar", "somevalue")
	Set("othervar", "true")
	Set("PORT", "8080")
	defer Unset("somevar")
	defer Unset("othervar")
	defer Unset("PORT")

	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "", cfg.Some)
	assert.Equal(t, false, cfg.Other)
	assert.Equal(t, 0, cfg.Port)
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.Error(t, Parse(&thisShouldBreak))
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidBool(t *testing.T) {
	Set("othervar", "should-be-a-bool")
	defer Unset("othervar")

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	Set("PORT", "should-be-an-int")
	defer Unset("PORT")

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestParsesDefaultConfig(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "postgres://localhost:5432/db", cfg.DatabaseURL)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Empty(t, cfg.NotAnEnv)
}
