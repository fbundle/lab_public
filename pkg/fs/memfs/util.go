package memfs

import (
	"strings"
)

func ensurePath(path []string) bool {
	for _, name := range path {
		if len(name) == 0 || strings.Contains(name, "/") {
			return false
		}
	}
	return true
}

func pathToKey(path []string) (string, bool) {
	if !ensurePath(path) {
		return "", false
	}
	return strings.Join(path, "/"), true
}

func keyToPath(key string) []string {
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}
	if len(key) > 0 && key[len(key)-1] == '/' {
		key = key[:len(key)-1]
	}
	return strings.Split(key, "/")
}
