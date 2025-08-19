package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// --- In-memory file ---
type MemFile struct {
	data []byte
}

func NewMemFile() *MemFile {
	return &MemFile{}
}

func (f *MemFile) Size() int {
	return len(f.data)
}

func (f *MemFile) Read(offset int, buf []byte) int {
	if offset >= len(f.data) {
		return 0
	}
	n := copy(buf, f.data[offset:])
	return n
}

func (f *MemFile) Write(offset int, buf []byte) int {
	end := offset + len(buf)
	if end > len(f.data) {
		newData := make([]byte, end)
		copy(newData, f.data)
		f.data = newData
	}
	copy(f.data[offset:], buf)
	return len(buf)
}

func (f *MemFile) Truncate(size int) {
	if size < len(f.data) {
		f.data = f.data[:size]
	} else if size > len(f.data) {
		f.data = append(f.data, make([]byte, size-len(f.data))...)
	}
}

// --- In-memory node (file or directory) ---
type MemNode struct {
	file     *MemFile
	children map[string]*MemNode
}

func NewDirNode() *MemNode {
	return &MemNode{children: make(map[string]*MemNode)}
}

func NewFileNode() *MemNode {
	return &MemNode{file: NewMemFile()}
}

// --- FUSE wrapper ---
type nodeWrapper struct {
	fs.Inode
	n *MemNode
}

// Lookup file or dir
func (n *nodeWrapper) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child, ok := n.n.children[name]
	if !ok {
		return nil, syscall.ENOENT
	}
	var mode uint32 = fuse.S_IFREG
	if child.file == nil {
		mode = fuse.S_IFDIR
	}
	return n.NewInode(ctx, &nodeWrapper{n: child}, fs.StableAttr{Mode: mode}), 0
}

// Read directory
func (n *nodeWrapper) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	entries := make([]fuse.DirEntry, 0, len(n.n.children))
	for name, child := range n.n.children {
		mode := uint32(fuse.S_IFREG)
		if child.file == nil {
			mode = fuse.S_IFDIR
		}
		entries = append(entries, fuse.DirEntry{Name: name, Mode: mode})
	}
	return fs.NewListDirStream(entries), 0
}

// Get attributes
func (n *nodeWrapper) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if n.n.file != nil {
		out.Size = uint64(n.n.file.Size())
		out.Mode = fuse.S_IFREG | 0644
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return 0
}

// Open file
func (n *nodeWrapper) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	if n.n.file != nil {
		return &fileWrapper{f: n.n.file}, fuse.FOPEN_KEEP_CACHE, 0
	}
	return nil, 0, syscall.EISDIR
}

// Create file
func (n *nodeWrapper) Create(ctx context.Context, name string, flags uint32, mode uint32) (fs.FileHandle, *fs.Inode, uint32, syscall.Errno) {
	child := NewFileNode()
	n.n.children[name] = child
	ino := n.NewInode(ctx, &nodeWrapper{n: child}, fs.StableAttr{Mode: fuse.S_IFREG})
	return &fileWrapper{f: child.file}, ino, fuse.FOPEN_KEEP_CACHE, 0
}

// Make directory
func (n *nodeWrapper) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child := NewDirNode()
	n.n.children[name] = child
	ino := n.NewInode(ctx, &nodeWrapper{n: child}, fs.StableAttr{Mode: fuse.S_IFDIR})
	return ino, 0
}

// Remove file
func (n *nodeWrapper) Unlink(ctx context.Context, name string) syscall.Errno {
	delete(n.n.children, name)
	return 0
}

// Remove directory
func (n *nodeWrapper) Rmdir(ctx context.Context, name string) syscall.Errno {
	delete(n.n.children, name)
	return 0
}

// Set attributes (truncate)
func (n *nodeWrapper) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if n.n.file == nil {
		return syscall.EPERM
	}
	if int(in.Size) != n.n.file.Size() {
		n.n.file.Truncate(int(in.Size))
		out.Size = in.Size
	}
	return 0
}

// --- File handle wrapper ---
type fileWrapper struct {
	fs.FileHandle
	f *MemFile
}

func (f *fileWrapper) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n := f.f.Read(int(off), dest)
	return fuse.ReadResultData(dest[:n]), 0
}

func (f *fileWrapper) Write(ctx context.Context, data []byte, off int64) (uint32, syscall.Errno) {
	n := f.f.Write(int(off), data)
	return uint32(n), 0
}

func (f *fileWrapper) Flush(ctx context.Context) syscall.Errno               { return 0 }
func (f *fileWrapper) Fsync(ctx context.Context, flags uint32) syscall.Errno { return 0 }
func (f *fileWrapper) Release(ctx context.Context) syscall.Errno             { return 0 }

// --- Mount ---
func main() {
	root := NewDirNode()
	rootInode := &nodeWrapper{n: root}

	mnt := "mnt"
	_ = os.MkdirAll(mnt, 0755)

	zeroDur := time.Duration(0)
	opts := &fs.Options{
		MountOptions: fuse.MountOptions{
			Debug: true,
		},
		EntryTimeout: &zeroDur,
		AttrTimeout:  &zeroDur,
	}

	server, err := fs.Mount(mnt, rootInode, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mounted at", mnt)
	server.Wait()
}
