package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"github.com/fbundle/go_util/pkg/fuse_util/node"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// Wrappers
type nodeWrapper struct {
	fs.Inode
	n node.Node
}

func (n *nodeWrapper) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child := n.n.Get(name)
	if child == nil {
		return nil, syscall.ENOENT
	}

	stable := fs.StableAttr{Mode: fuse.S_IFREG}
	if child.File() == nil {
		stable.Mode = fuse.S_IFDIR
	}

	wrapper := &nodeWrapper{n: child}
	return n.NewInode(ctx, wrapper, stable), 0
}

func (n *nodeWrapper) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	entries := []fuse.DirEntry{}

	n.n.Iter(func(name string, child node.Node) bool {
		mode := uint32(fuse.S_IFREG)
		if child.File() == nil {
			mode = fuse.S_IFDIR
		}
		entries = append(entries, fuse.DirEntry{Name: name, Mode: mode})
		return true
	})

	return fs.NewListDirStream(entries), 0
}

func (n *nodeWrapper) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if file := n.n.File(); file != nil {
		out.Size = uint64(file.Size())
		out.Mode = fuse.S_IFREG | 0644
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return 0
}

func (n *nodeWrapper) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	if file := n.n.File(); file != nil {
		return &fileWrapper{f: file}, fuse.FOPEN_KEEP_CACHE, 0
	}
	return nil, 0, syscall.EISDIR
}
func (n *nodeWrapper) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	// create a new directory Node in your model
	child := node.NewNode(nil)
	n.n.Set(name, child)

	stable := fs.StableAttr{Mode: fuse.S_IFDIR}
	wrapper := &nodeWrapper{n: child}
	return n.NewInode(ctx, wrapper, stable), 0
}

func (n *nodeWrapper) Create(ctx context.Context, name string, flags uint32, mode uint32) (fs.FileHandle, *fs.Inode, syscall.Errno) {
	// create a new file Node in your model
	child := node.NewNode(node.NewMemFile())
	n.n.Set(name, child)

	stable := fs.StableAttr{Mode: fuse.S_IFREG}
	wrapper := &nodeWrapper{n: child}
	ino := n.NewInode(ctx, wrapper, stable)

	fh := &fileWrapper{f: child.File()}
	return fh, ino, 0
}

type fileWrapper struct {
	fs.FileHandle
	f node.File
}

func (f *fileWrapper) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	buf := make([]byte, len(dest))
	n := f.f.Read(int(off), buf)
	return fuse.ReadResultData(buf[:n]), 0
}

func (f *fileWrapper) Write(ctx context.Context, data []byte, off int64) (uint32, syscall.Errno) {
	n := f.f.Write(int(off), data)
	return uint32(n), 0
}

func (f *fileWrapper) Flush(ctx context.Context) syscall.Errno {
	return 0
}

func (f *fileWrapper) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	return 0
}

func (f *fileWrapper) Release(ctx context.Context) syscall.Errno {
	return 0
}

// Mount
func Mount(root node.Node, mountpoint string) error {
	wrapped := &nodeWrapper{n: root}
	opts := &fs.Options{
		MountOptions: fuse.MountOptions{
			AllowOther: false,
			Debug:      true,
		},
	}

	server, err := fs.Mount(mountpoint, wrapped, opts)
	if err != nil {
		return err
	}

	log.Printf("Mounted at %s\n", mountpoint)
	server.Wait()
	return nil
}

func main() {
	var root node.Node = node.NewNode(nil)

	mnt := "tmp"
	_ = os.MkdirAll(mnt, 0755)
	if err := Mount(root, mnt); err != nil {
		log.Fatal(err)
	}
}
