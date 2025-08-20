package main

import (
	"log"

	"github.com/fbundle/go_util/pkg/fs"
	"github.com/fbundle/go_util/pkg/fs/memfile"
	"github.com/fbundle/go_util/pkg/fs/memfs"
)

func main() {
	memFS := memfs.NewFlatMemFS(memfile.NewMemFile)

	if err := fs.Mount(memFS, "tmp/memfs"); err != nil {
		log.Fatal(err)
	}
}
