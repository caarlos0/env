//go:build !windows

package env

import "testing"

func TestUnix(t *testing.T) {
	envVars := []string{":=/test/unix", "PATH=:/test_val1:/test_val2", "VAR=REGULARVAR", "FOO=", "BAR"}
	result := ToMap(envVars)
	isEqual(t, map[string]string{
		":":    "/test/unix",
		"PATH": ":/test_val1:/test_val2",
		"VAR":  "REGULARVAR",
		"FOO":  "",
	}, result)
}
