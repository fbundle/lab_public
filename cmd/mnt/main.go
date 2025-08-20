package main

import (
	"log"

	"github.com/fbundle/go_util/pkg/pathfs"
)

func main() {
	const mountPoint = "/tmp/mnt"
	if err := pathfs.MountMemPathFS(mountPoint); err != nil {
		log.Fatalf("mount failed: %v", err)
	}
}
