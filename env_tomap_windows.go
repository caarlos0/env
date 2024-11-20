//go:build windows

package env

import "strings"

func toMap(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		// On Windows, environment variables can start with '='. If so, Split at next character.
		// See env_windows.go in the Go source: https://github.com/golang/go/blob/go1.18/src/syscall/env_windows.go#L58
		e, prefixEqualSign := strings.CutPrefix(e, "=")
		p := strings.SplitN(e, "=", 2)
		if prefixEqualSign {
			p[0] = "=" + p[0]
		}

		if len(p) == 2 {
			r[p[0]] = p[1]
		}
	}
	return r
}
