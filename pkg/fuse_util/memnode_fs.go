package fuse_util

import (
	"github.com/fbundle/go_util/pkg/fuse_util/node"
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
