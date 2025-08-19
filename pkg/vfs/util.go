package vfs

import "strings"

func pathClean(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}
func pathConsume(path string) (string, string) {
	path = pathClean(path)
	parts := strings.SplitN(path, "/", 2)
	parts = append(parts, "")
	return parts[0], parts[1]
}
