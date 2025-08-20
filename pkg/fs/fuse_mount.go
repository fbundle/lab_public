package fs

import (
	"context"
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
)

type memPathFS struct {
	files map[string]File
}

const (
	PathSeparator = "/"
)

// MountMemPathFS mounts an in-memory PathFS at mountpoint using jacobsa/fuse.
// Directories are represented by nil File entries. Files are 0666 and directories 0777.
// UID/GID use the current process. allow_other is not enabled.
func MountMemPathFS(mountpoint string) error {
	if mountpoint == "" {
		return errors.New("mountpoint is required")
	}
	if err := os.MkdirAll(mountpoint, 0o755); err != nil {
		return err
	}
	mp := &memPathFS{files: make(map[string]File)}
	server := fuseutil.NewFileSystemServer(newMemFS(mp))
	mfs, err := fuse.Mount(mountpoint, server, &fuse.MountConfig{ReadOnly: false})
	if err != nil {
		return err
	}
	return mfs.Join(context.Background())
}

// jacobsa/fuse-backed FS implementation mapping to memPathFS

type memFS struct {
	fuseutil.NotImplementedFileSystem
	fs          *memPathFS
	inodeToPath map[fuseops.InodeID]Path
	pathToInode map[string]fuseops.InodeID
	nextInode   fuseops.InodeID
}

func newMemFS(p *memPathFS) *memFS {
	m := &memFS{
		fs:          p,
		inodeToPath: make(map[fuseops.InodeID]Path),
		pathToInode: make(map[string]fuseops.InodeID),
		nextInode:   fuseops.InodeID(2),
	}
	m.inodeToPath[fuseops.RootInodeID] = nil
	m.pathToInode[""] = fuseops.RootInodeID
	return m
}

func (m *memFS) inodeForPath(pth Path) fuseops.InodeID {
	key := strings.Join(pth, PathSeparator)
	if ino, ok := m.pathToInode[key]; ok {
		return ino
	}
	ino := m.nextInode
	m.nextInode++
	m.pathToInode[key] = ino
	m.inodeToPath[ino] = append(Path{}, pth...)
	return ino
}

func (m *memFS) pathForInode(inode fuseops.InodeID) (Path, bool) {
	pth, ok := m.inodeToPath[inode]
	return pth, ok
}

func (m *memFS) StatFS(ctx context.Context, op *fuseops.StatFSOp) error { return nil }

func (m *memFS) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	parent, ok := m.pathForInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	child := slices.Clone(parent)
	child = append(child, op.Name)
	key, _ := pathToKey(child)
	if file, ok := m.fs.files[key]; ok {
		ino := m.inodeForPath(child)
		op.Entry.Child = ino
		if file == nil {
			op.Entry.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: os.ModeDir | 0o777}
		} else {
			op.Entry.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: 0o666, Size: file.Size()}
		}
		return nil
	}
	prefix := key + PathSeparator
	for k := range m.fs.files {
		if strings.HasPrefix(k, prefix) {
			ino := m.inodeForPath(child)
			op.Entry.Child = ino
			op.Entry.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: os.ModeDir | 0o777}
			return nil
		}
	}
	return fuse.ENOENT
}

func (m *memFS) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	pth, ok := m.pathForInode(op.Inode)
	if !ok {
		return fuse.ENOENT
	}
	if len(pth) == 0 {
		op.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: os.ModeDir | 0o777}
		return nil
	}
	key := strings.Join(pth, PathSeparator)
	if file, ok := m.fs.files[key]; ok {
		if file == nil {
			op.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: os.ModeDir | 0o777}
		} else {
			op.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: 0o666, Size: file.Size()}
		}
		return nil
	}
	prefix := key + PathSeparator
	for k := range m.fs.files {
		if strings.HasPrefix(k, prefix) {
			op.Attributes = fuseops.InodeAttributes{Nlink: 1, Mode: os.ModeDir | 0o777}
			return nil
		}
	}
	return fuse.ENOENT
}

func (m *memFS) OpenDir(ctx context.Context, op *fuseops.OpenDirOp) error { return nil }

func (m *memFS) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) error {
	base, ok := m.pathForInode(op.Inode)
	if !ok {
		return fuse.ENOENT
	}
	depth := len(base)
	baseKey := strings.Join(base, PathSeparator)
	seen := make(map[string]fuseutil.Dirent)
	idx := 1
	for k, file := range m.fs.files {
		parts := strings.Split(k, PathSeparator)
		if depth > 0 {
			if strings.Join(parts[:min(depth, len(parts))], PathSeparator) != baseKey {
				continue
			}
		}
		if depth >= len(parts) {
			continue
		}
		name := parts[depth]
		dtype := fuseutil.DT_Directory
		if len(parts) == depth+1 && file != nil {
			dtype = fuseutil.DT_File
		}
		if _, ok := seen[name]; !ok {
			child := append(Path{}, base...)
			child = append(child, name)
			seen[name] = fuseutil.Dirent{Offset: fuseops.DirOffset(idx), Inode: m.inodeForPath(child), Name: name, Type: dtype}
			idx++
		}
	}
	entries := make([]fuseutil.Dirent, 0, len(seen))
	for _, e := range seen {
		entries = append(entries, e)
	}
	if op.Offset > fuseops.DirOffset(len(entries)) {
		return nil
	}
	entries = entries[op.Offset:]
	for _, e := range entries {
		n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], e)
		if n == 0 {
			break
		}
		op.BytesRead += n
	}
	return nil
}

func (m *memFS) MkDir(ctx context.Context, op *fuseops.MkDirOp) error {
	parent, ok := m.pathForInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	child := slices.Clone(parent)
	child = append(child, op.Name)
	key, _ := pathToKey(child)
	m.fs.files[key] = nil
	ino := m.inodeForPath(child)
	op.Entry = fuseops.ChildInodeEntry{Child: ino, Attributes: fuseops.InodeAttributes{Nlink: 1, Mode: os.ModeDir | 0o777}}
	return nil
}

func (m *memFS) CreateFile(ctx context.Context, op *fuseops.CreateFileOp) error {
	parent, ok := m.pathForInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	child := slices.Clone(parent)
	child = append(child, op.Name)
	key, _ := pathToKey(child)
	m.fs.files[key] = NewMemFile()
	ino := m.inodeForPath(child)
	op.Handle = fuseops.HandleID(ino)
	op.Entry = fuseops.ChildInodeEntry{Child: ino, Attributes: fuseops.InodeAttributes{Nlink: 1, Mode: 0o666}}
	return nil
}

func (m *memFS) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error { return nil }

func (m *memFS) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	pth, ok := m.pathForInode(op.Inode)
	if !ok {
		return fuse.ENOENT
	}
	key := strings.Join(pth, PathSeparator)
	file, ok := m.fs.files[key]
	if !ok || file == nil {
		return fuse.ENOENT
	}
	bytesRead := 0
	if err := file.Read(uint64(op.Offset), uint64(len(op.Dst)), func(src []byte) { copy(op.Dst, src); bytesRead = len(src) }); err != nil {
		return err
	}
	op.BytesRead = bytesRead
	return nil
}

func (m *memFS) WriteFile(ctx context.Context, op *fuseops.WriteFileOp) error {
	pth, ok := m.pathForInode(op.Inode)
	if !ok {
		return fuse.ENOENT
	}
	key := strings.Join(pth, PathSeparator)
	file, ok := m.fs.files[key]
	if !ok || file == nil {
		return fuse.ENOENT
	}
	if err := file.Write(uint64(op.Offset), uint64(len(op.Data)), func(dst []byte) { copy(dst, op.Data) }); err != nil {
		return err
	}
	return nil
}

func (m *memFS) SetInodeAttributes(ctx context.Context, op *fuseops.SetInodeAttributesOp) error {
	if op.Size != nil {
		pth, ok := m.pathForInode(op.Inode)
		if !ok {
			return fuse.ENOENT
		}
		key := strings.Join(pth, PathSeparator)
		file, ok := m.fs.files[key]
		if !ok || file == nil {
			return fuse.ENOENT
		}
		if err := file.Truncate(uint64(*op.Size)); err != nil {
			return err
		}
	}
	return nil
}

func (m *memFS) Unlink(ctx context.Context, op *fuseops.UnlinkOp) error {
	parent, ok := m.pathForInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	child := slices.Clone(parent)
	child = append(child, op.Name)
	key, _ := pathToKey(child)
	delete(m.fs.files, key)
	return nil
}

func (m *memFS) RmDir(ctx context.Context, op *fuseops.RmDirOp) error {
	parent, ok := m.pathForInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	child := slices.Clone(parent)
	child = append(child, op.Name)
	key, _ := pathToKey(child)
	prefix := key + PathSeparator
	for k := range m.fs.files {
		if strings.HasPrefix(k, prefix) {
			return fuse.ENOTEMPTY
		}
	}
	delete(m.fs.files, key)
	return nil
}
