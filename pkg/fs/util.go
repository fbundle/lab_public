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

type twoWayMap[T1 comparable, T2 comparable] struct {
	m12 map[T1]T2
	m21 map[T2]T1
}

func (m *twoWayMap[T1, T2]) Set(k1 T1, k2 T2) {
	m.m12[k1] = k2
	m.m21[k2] = k1
}

func (m *twoWayMap[T1, T2]) Get1(k2 T2) (T1, bool) {
	v, ok := m.m21[k2]
	return v, ok
}
func (m *twoWayMap[T1, T2]) Get2(k1 T1) (T2, bool) {
	v, ok := m.m12[k1]
	return v, ok
}
