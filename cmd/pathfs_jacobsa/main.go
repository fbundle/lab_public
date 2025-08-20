package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/fbundle/go_util/pkg/vfs"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
)

// pathKey joins a path slice with '/'
func pathKey(p []string) string { return strings.Join(p, "/") }

// split and sanitize a path string into components
func splitClean(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	parts := strings.Split(p, "/")
	res := make([]string, 0, len(parts))
	for _, s := range parts {
		if s == "." || s == ".." || s == "" {
			continue
		}
		res = append(res, s)
	}
	return res
}

type nodeType int

const (
	nodeTypeFile nodeType = iota
	nodeTypeDir
)

type pathFuseFS struct {
	fuseutil.NotImplementedFileSystem

	store vfs.PathFS

	uid uint32
	gid uint32

	// inode bookkeeping
	nextInode    fuseops.InodeID
	pathToInode  map[string]fuseops.InodeID
	inodeToPath  map[fuseops.InodeID][]string
	explicitDirs map[string]struct{}
}

func newPathFuseFS(store vfs.PathFS, uid, gid uint32) *pathFuseFS {
	fs := &pathFuseFS{
		store:        store,
		uid:          uid,
		gid:          gid,
		nextInode:    fuseops.InodeID(fuseops.RootInodeID + 1),
		pathToInode:  make(map[string]fuseops.InodeID),
		inodeToPath:  make(map[fuseops.InodeID][]string),
		explicitDirs: make(map[string]struct{}),
	}
	fs.pathToInode[""] = fuseops.RootInodeID
	fs.inodeToPath[fuseops.RootInodeID] = nil
	return fs
}

func (p *pathFuseFS) ensureInodeFor(path []string) fuseops.InodeID {
	key := pathKey(path)
	if id, ok := p.pathToInode[key]; ok {
		return id
	}
	id := p.nextInode
	p.nextInode++
	p.pathToInode[key] = id
	p.inodeToPath[id] = append([]string{}, path...)
	return id
}

// examinePath determines whether a child path exists as file and/or has deeper entries (dir)
func (pfs *pathFuseFS) examinePath(childPath []string) (isFile bool, isDir bool, f vfs.File) {
	isFile = false
	isDir = false
	var theFile vfs.File
	pfs.store.Walk(func(p []string, file vfs.File) bool {
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
			theFile = file
			return true
		}
		isDir = true
		return true
	})
	return isFile, isDir, theFile
}

func (pfs *pathFuseFS) StatFS(ctx context.Context, op *fuseops.StatFSOp) error { return nil }

func (pfs *pathFuseFS) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	parentPath, ok := pfs.inodeToPath[op.Parent]
	if !ok {
		return fuse.ENOENT
	}
	childPath := append(append([]string{}, parentPath...), op.Name)
	isFile, isDir, file := pfs.examinePath(childPath)

	if !isFile && !isDir {
		return fuse.ENOENT
	}

	var childInode fuseops.InodeID
	childInode = pfs.ensureInodeFor(childPath)
	op.Entry.Child = childInode
	if isFile {
		op.Entry.Attributes = fuseops.InodeAttributes{Mode: 0o777, Nlink: 1}
		op.Entry.Attributes.Size = file.Length()
	} else {
		op.Entry.Attributes = fuseops.InodeAttributes{Mode: os.ModeDir | 0o777, Nlink: 1}
	}
	// Set ownership if supported by library (Uid/Gid are present in struct in jacobsa/fuse)
	op.Entry.Attributes.Uid = pfs.uid
	op.Entry.Attributes.Gid = pfs.gid
	return nil
}

func (pfs *pathFuseFS) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	path, ok := pfs.inodeToPath[op.Inode]
	if !ok {
		return fuse.ENOENT
	}
	isFile, isDir, file := pfs.examinePath(path)
	if !isFile && !isDir {
		// root or explicit empty dir
		if op.Inode == fuseops.RootInodeID || pfs.isExplicitDir(path) {
			op.Attributes = fuseops.InodeAttributes{Mode: os.ModeDir | 0o777, Nlink: 1}
			op.Attributes.Uid = pfs.uid
			op.Attributes.Gid = pfs.gid
			return nil
		}
		return fuse.ENOENT
	}
	if isFile {
		op.Attributes = fuseops.InodeAttributes{Mode: 0o777, Nlink: 1}
		op.Attributes.Size = file.Length()
	} else {
		op.Attributes = fuseops.InodeAttributes{Mode: os.ModeDir | 0o777, Nlink: 1}
	}
	op.Attributes.Uid = pfs.uid
	op.Attributes.Gid = pfs.gid
	return nil
}

func (pfs *pathFuseFS) isExplicitDir(path []string) bool {
	_, ok := pfs.explicitDirs[pathKey(path)]
	return ok
}

func (pfs *pathFuseFS) OpenDir(ctx context.Context, op *fuseops.OpenDirOp) error { return nil }

func (pfs *pathFuseFS) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) error {
	path, ok := pfs.inodeToPath[op.Inode]
	if !ok {
		return fuse.ENOENT
	}

	childrenSet := make(map[string]nodeType)
	// from files
	pfs.store.Walk(func(p []string, f vfs.File) bool {
		if len(p) < len(path)+1 {
			return true
		}
		for i := range path {
			if p[i] != path[i] {
				return true
			}
		}
		name := p[len(path)]
		t := nodeTypeFile
		if len(p) > len(path)+1 {
			t = nodeTypeDir
		}
		if prev, ok := childrenSet[name]; !ok || prev != nodeTypeDir {
			childrenSet[name] = t
		}
		return true
	})
	// from explicit dirs
	for key := range pfs.explicitDirs {
		parts := splitClean(key)
		if len(parts) == len(path)+1 {
			match := true
			for i := range path {
				if parts[i] != path[i] {
					match = false
					break
				}
			}
			if match {
				childrenSet[parts[len(path)]] = nodeTypeDir
			}
		}
	}

	// Sort for stable output
	var names []string
	for n := range childrenSet {
		names = append(names, n)
	}
	sort.Strings(names)

	// Honor offset
	if int(op.Offset) < len(names) {
		names = names[op.Offset:]
	} else {
		names = nil
	}

	for _, name := range names {
		childPath := append(append([]string{}, path...), name)
		id := pfs.ensureInodeFor(childPath)
		var typ fuseutil.DirentType
		if childrenSet[name] == nodeTypeDir {
			typ = fuseutil.DT_Directory
		} else {
			typ = fuseutil.DT_File
		}
		de := fuseutil.Dirent{Inode: id, Name: name, Type: typ}
		n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], de)
		if n == 0 {
			break
		}
		op.BytesRead += n
	}
	return nil
}

func (pfs *pathFuseFS) MkDir(ctx context.Context, op *fuseops.MkDirOp) error {
	parentPath, ok := pfs.inodeToPath[op.Parent]
	if !ok {
		return fuse.ENOENT
	}
	child := append(append([]string{}, parentPath...), op.Name)
	pfs.explicitDirs[pathKey(child)] = struct{}{}
	childInode := pfs.ensureInodeFor(child)
	op.Entry.Child = childInode
	op.Entry.Attributes = fuseops.InodeAttributes{Mode: os.ModeDir | 0o777, Nlink: 1}
	op.Entry.Attributes.Uid = pfs.uid
	op.Entry.Attributes.Gid = pfs.gid
	return nil
}

func (pfs *pathFuseFS) CreateFile(ctx context.Context, op *fuseops.CreateFileOp) error {
	parentPath, ok := pfs.inodeToPath[op.Parent]
	if !ok {
		return fuse.ENOENT
	}
	child := append(append([]string{}, parentPath...), op.Name)
	file, err := pfs.store.OpenOrCreate(child)
	if err != nil {
		return fuse.EIO
	}
	_ = file.Truncate(0)
	childInode := pfs.ensureInodeFor(child)
	op.Entry.Child = childInode
	op.Entry.Attributes = fuseops.InodeAttributes{Mode: 0o777, Nlink: 1}
	op.Entry.Attributes.Uid = pfs.uid
	op.Entry.Attributes.Gid = pfs.gid
	op.Handle = fuseops.HandleID(childInode)
	return nil
}

func (pfs *pathFuseFS) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error {
	// Use inode as handle for simplicity
	op.Handle = fuseops.HandleID(op.Inode)
	return nil
}

func (pfs *pathFuseFS) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	path, ok := pfs.inodeToPath[op.Inode]
	if !ok {
		return fuse.ENOENT
	}
	isFile, _, file := pfs.examinePath(path)
	if !isFile {
		return fuse.ENOENT
	}
	var out []byte
	if err := file.Read(uint64(op.Offset), uint64(len(op.Dst)), func(b []byte) { out = append([]byte{}, b...) }); err != nil {
		return fuse.EIO
	}
	copy(op.Dst, out)
	op.BytesRead = len(out)
	return nil
}

func (pfs *pathFuseFS) WriteFile(ctx context.Context, op *fuseops.WriteFileOp) error {
	path, ok := pfs.inodeToPath[op.Inode]
	if !ok {
		return fuse.ENOENT
	}
	isFile, _, file := pfs.examinePath(path)
	if !isFile {
		// If file doesn't exist yet but a handle points here, create it.
		var err error
		file, err = pfs.store.OpenOrCreate(path)
		if err != nil {
			return fuse.EIO
		}
	}
	if err := file.Write(uint64(op.Offset), uint64(len(op.Data)), func(dst []byte) { copy(dst, op.Data) }); err != nil {
		return fuse.EIO
	}
	return nil
}

func (pfs *pathFuseFS) SetInodeAttributes(ctx context.Context, op *fuseops.SetInodeAttributesOp) error {
	path, ok := pfs.inodeToPath[op.Inode]
	if !ok {
		return fuse.ENOENT
	}
	if op.Size != nil {
		sz := *op.Size
		isFile, _, file := pfs.examinePath(path)
		if !isFile {
			return fuse.ENOENT
		}
		if err := file.Truncate(uint64(sz)); err != nil {
			return fuse.EIO
		}
	}
	return nil
}

func (pfs *pathFuseFS) Unlink(ctx context.Context, op *fuseops.UnlinkOp) error {
	parentPath, ok := pfs.inodeToPath[op.Parent]
	if !ok {
		return fuse.ENOENT
	}
	child := append(append([]string{}, parentPath...), op.Name)
	if err := pfs.store.Unlink(child); err != nil {
		return fuse.ENOENT
	}
	return nil
}

func main() {
	mountPoint := flag.String("mount_point", "/tmp/mnt_jacobsa", "mount point")
	uidFlag := flag.Int("uid", -1, "ownership UID to report (default: current user or $SUDO_UID when running as root)")
	gidFlag := flag.Int("gid", -1, "ownership GID to report (default: current group or $SUDO_GID when running as root)")
	readOnly := flag.Bool("read_only", false, "mount read-only")
	debug := flag.Bool("debug", false, "enable fuse debug logs")
	flag.Parse()

	if err := os.MkdirAll(*mountPoint, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", *mountPoint, err)
	}

	uid := os.Getuid()
	gid := os.Getgid()
	if uid == 0 {
		if v := os.Getenv("SUDO_UID"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				uid = n
			}
		}
		if v := os.Getenv("SUDO_GID"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				gid = n
			}
		}
	}
	if *uidFlag >= 0 {
		uid = *uidFlag
	}
	if *gidFlag >= 0 {
		gid = *gidFlag
	}

	backend := vfs.NewMemPathFS()
	if f, err := backend.OpenOrCreate([]string{"xyz", "hello.txt"}); err == nil {
		_ = f.Truncate(0)
		_ = f.Write(0, uint64(len("Hello from PathFS\n")), func(b []byte) { copy(b, []byte("Hello from PathFS\n")) })
	}

	fsrv := newPathFuseFS(backend, uint32(uid), uint32(gid))

	cfg := &fuse.MountConfig{
		ReadOnly: *readOnly,
	}
	if *debug {
		cfg.DebugLogger = log.New(os.Stderr, "fuse: ", 0)
	}

	abs, _ := filepath.Abs(*mountPoint)
	server, err := fuse.Mount(abs, fuseutil.NewFileSystemServer(fsrv), cfg)
	if err != nil {
		log.Fatalf("mount: %v", err)
	}
	log.Printf("mounted pathfs_jacobsa at %s", abs)
	if err := server.Join(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Join: %v\n", err)
	}
}
