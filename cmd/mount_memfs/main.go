package main

import (
	"log"

	"github.com/fbundle/go_util/pkg/fs"
)

func main() {
	memfs := fs.NewFlatMemFS(fs.NewMemFile)

	if err := fs.Mount(memfs, "tmp/memfs"); err != nil {
		log.Fatal(err)
	}
}
