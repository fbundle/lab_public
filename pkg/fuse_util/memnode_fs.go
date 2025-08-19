package fuse_util

import (
	"github.com/fbundle/go_util/pkg/fuse_util/node"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
)

func newFsFile(file node.File) nodefs.File {
	return &fsFile{
		File: nodefs.NewDefaultFile(),
		file: file,
	}
}

type fsFile struct {
	nodefs.File
	file node.File
}

func (f *fsFile) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	f.file.Read(int(off), dest)
	return nil, fuse.OK
}
func (f *fsFile) Write(data []byte, off int64) (written uint32, code fuse.Status) {
	n := f.file.Write(int(off), data)
	return uint32(n), fuse.OK
}
func (f *fsFile) Truncate(size uint64) fuse.Status {
	f.file.Truncate(int(size))
	return fuse.OK
}

func (f *fsFile) GetAttr(out *fuse.Attr) fuse.Status {
	*out = *f.file.Attr()
	return fuse.OK
}

func newFsNode(node node.Node) nodefs.Node {
	return &fsNode{
		Node: nodefs.NewDefaultNode(),
		node: node,
	}
}

type fsNode struct {
	nodefs.Node
	node node.Node
}

func (n *fsNode) Mkdir(name string, mode uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	child := newFsNode(node.NewDirNode())
	n.Inode().NewChild(name, true, child)
	return child.Inode(), fuse.OK
}

func (n *fsNode) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	child := n.Inode().RmChild(name)
	if child == nil {
		return fuse.ENOENT
	}
	return fuse.OK
}

func (n *fsNode) Rename(oldName string, newParent nodefs.Node, newName string, context *fuse.Context) (code fuse.Status) {
	child := n.Inode().RmChild(oldName)
	newParent.Inode().RmChild(newName)
	newParent.Inode().AddChild(newName, child)
	return fuse.OK
}

func (n *fsNode) Open(flags uint32, context *fuse.Context) (ifile nodefs.File, code fuse.Status) {
	file := n.node.File()
	return newFsFile(file), fuse.OK
}

func (n *fsNode) Create(name string, flags uint32, mode uint32, context *fuse.Context) (ifile nodefs.File, inode *nodefs.Inode, code fuse.Status) {
	file := node.NewFile()
	child := newFsNode(node.NewFileNode(file))
	n.Inode().NewChild(name, false, child)

	return newFsFile(file), child.Inode(), fuse.OK
}

func (n *fsNode) GetAttr(fi *fuse.Attr, file nodefs.File, context *fuse.Context) (code fuse.Status) {
	*fi = *n.node.Attr()
	return fuse.OK
}
