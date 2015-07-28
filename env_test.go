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
	Some string `env:"somevar"`
}

func TestParsesEnv(t *testing.T) {
	key := "somevar"
	value := "somevalue"
	Set(key, value)
	defer Unset(key)
	env := env1{}
	err := ParseEnv(env, &env)
	assert.Nil(t, err)
	assert.Equal(t, value, env.Some)
}
