package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	Some  string `env:"somevar"`
	Other bool   `env:"othervar"`
	Port  int    `env:"PORT"`
}

func TestParsesEnv(t *testing.T) {
	Set("somevar", "somevalue")
	Set("othervar", "true")
	Set("PORT", "8080")
	defer Unset("somevar")
	defer Unset("othervar")
	defer Unset("PORT")

	cfg := Config{}
	err := Parse(&cfg)

	assert.Nil(t, err)
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	err := Parse(&cfg)

	assert.Nil(t, err)
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
