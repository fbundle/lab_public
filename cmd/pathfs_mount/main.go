package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/fbundle/go_util/pkg/vfs"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// split and sanitize a FUSE path into PathFS segments.
func splitClean(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	parts := strings.Split(p, "/")
	res := make([]string, 0, len(parts))
	for _, s := range parts {
		if s == "." || s == "" || s == ".." {
			continue
		}
		res = append(res, s)
	}
	return res
}

// NodeDir represents a directory backed by PathFS prefixes.
type NodeDir struct {
	fs.Inode
	store vfs.PathFS
	path  []string // path prefix represented by this directory
}

// NodeFile represents a file backed by vfs.File.
type NodeFile struct {
	fs.Inode
	file vfs.File
	name string
}

// FileHandle implements fs.FileHandle for NodeFile I/O.
type FileHandle struct{ file vfs.File }

// Directory attribute
func (d *NodeDir) Getattr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR | 0o755
	out.Uid = uint32(os.Getuid())
	out.Gid = uint32(os.Getgid())
	out.Nlink = 1
	return 0
}

// File attribute
func (n *NodeFile) Getattr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFREG | 0o644
	out.Uid = uint32(os.Getuid())
	out.Gid = uint32(os.Getgid())
	out.Nlink = 1
	out.Size = n.file.Length()
	return 0
}

// Dir: Lookup child by name; decides between file and directory based on backend keys.
func (d *NodeDir) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	childPath := append(append([]string{}, d.path...), name)

	// Determine if child exists as file or directory by scanning.
	isFile := false
	isDir := false
	var theFile vfs.File

	d.store.Walk(func(p []string, f vfs.File) bool {
		if len(p) < len(childPath) {
			return true
		}
		match := true
		for i := range childPath {
			if p[i] != childPath[i] {
				match = false
				break
			}
		}
		if !match {
			return true
		}
		if len(p) == len(childPath) {
			isFile = true
			theFile = f
			return false
		}
		isDir = true
		return true
	})

	if !isFile && !isDir {
		return nil, syscall.ENOENT
	}

	if isDir && !isFile {
		ch := &NodeDir{store: d.store, path: childPath}
		return d.NewInode(ctx, ch, fs.StableAttr{Mode: uint32(fuse.S_IFDIR)}), 0
	}

	// Prefer file when both match (file node at exact path also has deeper keys).
	ch := &NodeFile{file: theFile, name: name}
	return d.NewInode(ctx, ch, fs.StableAttr{Mode: uint32(fuse.S_IFREG)}), 0
}

// Dir: Readdir lists immediate children under this prefix.
func (d *NodeDir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	children := make(map[string]uint32)

	d.store.Walk(func(p []string, f vfs.File) bool {
		if len(p) < len(d.path) {
			return true
		}
		for i := range d.path {
			if p[i] != d.path[i] {
				return true
			}
		}
		if len(p) == len(d.path) {
			return true
		}
		name := p[len(d.path)]
		mode := uint32(fuse.S_IFREG)
		if len(p) > len(d.path)+1 {
			mode = uint32(fuse.S_IFDIR)
		}
		// directory wins over file if both seen
		if prev, ok := children[name]; !ok || prev != uint32(fuse.S_IFDIR) {
			children[name] = mode
		}
		return true
	})

	list := make([]fuse.DirEntry, 0, len(children))
	for name, mode := range children {
		list = append(list, fuse.DirEntry{Name: name, Mode: mode})
	}
	return fs.NewListDirStream(list), 0
}

// Dir: Mkdir is accepted (directories are implicit), no specific storage needed.
func (d *NodeDir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child := &NodeDir{store: d.store, path: append(append([]string{}, d.path...), name)}
	return d.NewInode(ctx, child, fs.StableAttr{Mode: uint32(fuse.S_IFDIR)}), 0
}

// Dir: Create a file; ensures an empty file exists in backend.
func (d *NodeDir) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {
	parts := append(append([]string{}, d.path...), name)
	f, err := d.store.OpenOrCreate(parts)
	if err != nil {
		return nil, nil, 0, syscall.EIO
	}
	n := &NodeFile{file: f, name: name}
	inode := d.NewInode(ctx, n, fs.StableAttr{Mode: uint32(fuse.S_IFREG)})
	return inode, &FileHandle{file: f}, 0, 0
}

// Dir: Unlink deletes a file if present at exact path.
func (d *NodeDir) Unlink(ctx context.Context, name string) syscall.Errno {
	parts := append(append([]string{}, d.path...), name)
	if err := d.store.Unlink(parts); err != nil {
		return syscall.ENOENT
	}
	return 0
}

// File: Open returns a handle.
func (n *NodeFile) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	return &FileHandle{file: n.file}, 0, 0
}

// File: Setattr supports truncate.
func (n *NodeFile) Setattr(ctx context.Context, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if sz, ok := in.GetSize(); ok {
		if err := n.file.Truncate(uint64(sz)); err != nil {
			return syscall.EIO
		}
	}
	return n.Getattr(ctx, out)
}

// Handle: Read implementation.
func (h *FileHandle) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	var out []byte
	readErr := h.file.Read(uint64(off), uint64(len(dest)), func(data []byte) {
		out = make([]byte, len(data))
		copy(out, data)
	})
	if readErr != nil {
		return nil, syscall.EIO
	}
	return fuse.ReadResultData(out), 0
}

// Handle: Write implementation.
func (h *FileHandle) Write(ctx context.Context, data []byte, off int64) (uint32, syscall.Errno) {
	if err := h.file.Write(uint64(off), uint64(len(data)), func(dst []byte) {
		copy(dst, data)
	}); err != nil {
		return 0, syscall.EIO
	}
	return uint32(len(data)), 0
}

func main() {
	mountPoint := flag.String("mount_point", "/tmp/mnt", "mount point directory")
	flag.Parse()

	if err := os.MkdirAll(*mountPoint, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", *mountPoint, err)
	}

	backend := vfs.NewMemPathFS()
	if demoFile, err := backend.OpenOrCreate([]string{"xyz/hello.txt"}); err == nil {
		_ = demoFile.Truncate(0)
		_ = demoFile.Write(0, uint64(len("Hello from PathFS\n")), func(b []byte) { copy(b, []byte("Hello from PathFS\n")) })
	}

	root := &NodeDir{store: backend, path: nil}
	server, err := fs.Mount(*mountPoint, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			AllowOther: false,
			Name:       "pathfs",
			FsName:     path.Base(*mountPoint),
		},
	})
	if err != nil {
		log.Fatalf("mount: %v", err)
	}
	log.Printf("mounted PathFS(v2) at %s (Ctrl+C to unmount)", *mountPoint)
	server.Wait()
}
