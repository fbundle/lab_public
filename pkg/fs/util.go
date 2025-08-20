package fs

import (
	"strings"
)

func ensurePath(path []string) {
	for _, name := range path {
		if len(name) == 0 || strings.Contains(name, "/") {
			panic("invalid path")
		}
	}
}

func pathToKey(path []string) string {
	ensurePath(path)
	return strings.Join(path, "/")
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
