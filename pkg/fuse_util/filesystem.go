package fuse_util

import (
	"context"
	"errors"
	"slices"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
)

func NewFuseFileSystem(files FileStore) fuseutil.FileSystem {
	m := &memFS{
		files: files,
		inodePool: newInodePool(node{
			inode: fuseops.RootInodeID,
			path:  nil,
			file:  newDir(),
		}),
	}
	// populate the inodes for directories and files
	if err := files.Iterate(func(file File) bool {
		path := mustAttr(file).Path
		if len(path) == 0 {
			panic("file must not have empty path")
		}
		// populate the parent directories
		for i := 0; i < len(path); i++ {
			dirpath := path[:i]
			if _, ok := m.inodePool.getNodeFromPath(dirpath); !ok {
				m.inodePool.createNode(dirpath, newDir())
			}
		}
		// set the file inode
		m.inodePool.createNode(path, file)
		return true
	}); err != nil {
		panic(err)
	}

	m.inodePool.pathToNode.ReduceAll(mtimeReducer)
	return m
}

type memFS struct {
	fuseutil.NotImplementedFileSystem
	files     FileStore
	inodePool *inodePool
}

func (m *memFS) StatFS(ctx context.Context, op *fuseops.StatFSOp) error {
	return nil
}

// LookUpInode - get child info
func (m *memFS) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	parent, ok := m.inodePool.getNodeFromInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	path := append(slices.Clone(parent.path), op.Name)

	node, ok := m.inodePool.getNodeFromPath(path)
	if !ok {
		return fuse.ENOENT
	}

	op.Entry.Child = node.inode
	op.Entry.Attributes = getInodeAttributes(mustAttr(node.file))
	return nil
}

// GetInodeAttributes - get self info
func (m *memFS) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	node, ok := m.inodePool.getNodeFromInode(op.Inode)
	if !ok {
		return fuse.ENOENT
	}
	op.Attributes = getInodeAttributes(mustAttr(node.file))
	return nil
}

func (m *memFS) OpenDir(ctx context.Context, op *fuseops.OpenDirOp) error {
	return nil
}

// ReadDir - analogous to `ls`
func (m *memFS) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) error {
	node, ok := m.inodePool.getNodeFromInode(op.Inode)
	if !ok {
		return fuse.ENOENT
	}

	var entries []fuseutil.Dirent
	var offset fuseops.DirOffset = 1
	for name, child := range m.inodePool.list(node.path) {
		dtype := fuseutil.DT_Directory
		if !child.isDir() {
			dtype = fuseutil.DT_File
		}
		entries = append(entries, fuseutil.Dirent{
			Offset: offset,
			Inode:  child.inode,
			Name:   name,
			Type:   dtype,
		})
		offset += 1
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

// MkDir -
func (m *memFS) MkDir(ctx context.Context, op *fuseops.MkDirOp) error {
	parent, ok := m.inodePool.getNodeFromInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}

	path := append(slices.Clone(parent.path), op.Name)
	node, ok := m.inodePool.createNode(path, newDir())
	if !ok {
		return fuse.ENOENT
	}

	op.Entry = fuseops.ChildInodeEntry{
		Child:      node.inode,
		Attributes: getInodeAttributes(mustAttr(node.file)),
	}
	return nil
}

func (m *memFS) CreateFile(ctx context.Context, op *fuseops.CreateFileOp) error {
	parent, ok := m.inodePool.getNodeFromInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}

	path := append(slices.Clone(parent.path), op.Name)

	file, err := m.files.Create()
	if err != nil {
		return err
	}

	err = file.UpdateAttr(func(attr FileAttr) FileAttr {
		attr.Path = path
		return attr
	})
	if err != nil {
		return err
	}

	node, ok := m.inodePool.createNode(path, file)
	if !ok {
		return fuse.ENOENT
	}

	op.Handle = fuseops.HandleID(node.inode)
	op.Entry = fuseops.ChildInodeEntry{
		Child:      node.inode,
		Attributes: getInodeAttributes(mustAttr(node.file)),
	}
	return nil
}

// RmDir - rmdir
func (m *memFS) RmDir(ctx context.Context, op *fuseops.RmDirOp) error {
	parent, ok := m.inodePool.getNodeFromInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	path := append(slices.Clone(parent.path), op.Name)
	ok = m.inodePool.deleteNodeIf(path, func(n node) bool {
		return n.isDir()
	})
	if !ok {
		return fuse.ENOENT
	}
	return nil
}

// Unlink - rm
func (m *memFS) Unlink(ctx context.Context, op *fuseops.UnlinkOp) error {
	parent, ok := m.inodePool.getNodeFromInode(op.Parent)
	if !ok {
		return fuse.ENOENT
	}
	path := append(slices.Clone(parent.path), op.Name)

	var err error
	if ok := m.inodePool.deleteNodeIf(path, func(n node) bool {
		if n.isDir() {
			return false
		}
		id := mustAttr(n.file).ID
		err = m.files.Delete(id)
		if err != nil {
			return false
		}
		return true
	}); !ok {
		if err != nil {
			return err
		}
		return fuse.ENOENT
	}

	return nil
}

// TODO - support rename

func (m *memFS) Rename(ctx context.Context, op *fuseops.RenameOp) error {
	oldParent, ok := m.inodePool.getNodeFromInode(op.OldParent)
	oldName := op.OldName
	newParent, ok := m.inodePool.getNodeFromInode(op.NewParent)
	newName := op.NewName
	_ = []any{oldParent, oldName, newParent, newName, ok}
	return errors.New("not implemented yet")
}

func (m *memFS) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error {
	return nil
}

func (m *memFS) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	node, ok := m.inodePool.getNodeFromInode(op.Inode)
	if !ok || node.isDir() {
		return fuse.ENOENT
	}
	var err error
	op.BytesRead, err = node.file.Read(uint64(op.Offset), op.Dst)
	return err
}

func (m *memFS) WriteFile(ctx context.Context, op *fuseops.WriteFileOp) error {
	node, ok := m.inodePool.getNodeFromInode(op.Inode)
	if !ok || node.isDir() {
		return fuse.ENOENT
	}
	defer m.inodePool.pathToNode.ReducePartial(node.path, mtimeReducer)

	_, err := node.file.Write(uint64(op.Offset), op.Data)
	return err
}

// SetInodeAttributes -
func (m *memFS) SetInodeAttributes(ctx context.Context, op *fuseops.SetInodeAttributesOp) error {
	if op.Size != nil { // truncate
		node, ok := m.inodePool.getNodeFromInode(op.Inode)
		if !ok || node.isDir() {
			return fuse.ENOENT
		}

		defer m.inodePool.pathToNode.ReducePartial(node.path, mtimeReducer)

		return node.file.Trunc(*op.Size)
	}
	// other changes like mode are not supported yet
	return errors.New("not implemented yet")
}
