package main

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/fbundle/go_util/pkg/fuse_util"
	"github.com/fbundle/go_util/pkg/fuse_util/mem"
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseutil"
)

func mustRunCmd(command string) {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	_, _ = output, err
}

func mount(fs fuse_util.FileStore, mountpoint string) error {
	if mountpoint == "" {
		return errors.New("mountpoint is required")
	}

	server := fuseutil.NewFileSystemServer(fuse_util.NewFuseFileSystem(fs))
	mfs, err := fuse.Mount(mountpoint, server, &fuse.MountConfig{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}
	return mfs.Join(context.Background())
}
func main() {
	func() {
		mustRunCmd("fusermount -u mnt")
		mustRunCmd("mkdir mnt")
	}()
	defer func() {
		mustRunCmd("fusermount -u mnt")
		mustRunCmd("rm -rf mnt")
	}()

	files := fuse_util_mem.NewMemFileStore()
	if err := mount(files, "mnt"); err != nil {
		log.Fatal(err)
	}
}
