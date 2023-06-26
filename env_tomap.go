//go:build !windows

package env

import "strings"

func toMap(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		if len(p) == 2 {
			r[p[0]] = p[1]
		}
	}
	return r
}
