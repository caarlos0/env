package env_test

import (
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
)

const key = "WHATEVER"

func TestGetenv(t *testing.T) {
	value := "bla"
	env.Set(key, value)
	defer env.Unset(key)
	assert.Equal(t, value, env.GetOr(key, "default"))
}

func TestGetenvUnseted(t *testing.T) {
	value := "default"
	assert.Equal(t, value, env.GetOr(key, value))
}
