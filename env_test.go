package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const key = "WHATEVER"

func TestGetenv(t *testing.T) {
	value := "bla"
	Set(key, value)
	defer Unset(key)
	assert.Equal(t, value, GetOr(key, "default"))
}

func TestGetenvUnseted(t *testing.T) {
	value := "default"
	assert.Equal(t, value, GetOr(key, value))
}

type env1 struct {
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

	env := env1{}
	err := ParseEnv(env, &env)

	assert.Nil(t, err)
	assert.Equal(t, "somevalue", env.Some)
	assert.Equal(t, true, env.Other)
	assert.Equal(t, 8080, env.Port)
}
