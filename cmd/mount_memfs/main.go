package main

import (
	"log"

	"github.com/fbundle/go_util/pkg/fs"
)

func main() {
	if err := fs.MountMemPathFS("tmp/memfs"); err != nil {
		log.Fatal(err)
	}
}
