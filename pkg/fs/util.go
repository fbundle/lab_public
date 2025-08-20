package fs

import (
	"strings"
)

func ensurePath(path []string) []string {
	var pathOut []string = nil
	for _, name := range path {
		if len(name) == 0 {
			continue
		}
		if strings.Contains(name, "/") {
			panic("invalid path")
		}
		pathOut = append(pathOut, name)
	}
	return pathOut
}

func pathToKey(path []string) string {
	path = ensurePath(path)
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
