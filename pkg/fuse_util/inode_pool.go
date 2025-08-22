package fuse_util

import (
	"sync"

	"github.com/fbundle/go_util/pkg/fuse_util/trie"
	"github.com/jacobsa/fuse/fuseops"
)

func newInodePool(rootNode node) *inodePool {
	return &inodePool{
		mu:       sync.RWMutex{},
		maxInode: rootNode.inode,
		inodeToNode: map[fuseops.InodeID]node{
			rootNode.inode: rootNode,
		},
		pathToNode: trie.New[string, node](rootNode),
	}
}

func newNode(inode fuseops.InodeID, path []string, file File) node {
	return node{
		inode: inode,
		path:  path,
		file:  file,
	}
}

type node struct {
	inode fuseops.InodeID
	path  []string

	file File
}

func (n node) isDir() bool {
	return mustAttr(n.file).IsDir
}

type inodePool struct {
	mu          sync.RWMutex
	inodeToNode map[fuseops.InodeID]node
	pathToNode  *trie.Trie[string, node]
	maxInode    fuseops.InodeID
}

func (p *inodePool) list(prefix []string) func(yield func(name string, child node) bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.pathToNode.List(prefix)
}

func (p *inodePool) getNodeFromPath(path []string) (node, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.pathToNode.Load(path)
}

func (p *inodePool) getNodeFromInode(inode fuseops.InodeID) (node, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	n, ok := p.inodeToNode[inode]
	return n, ok
}

func (p *inodePool) createNode(path []string, file File) (n node, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	n = newNode(p.maxInode+1, path, file)

	p.maxInode = max(p.maxInode, n.inode)

	ok = p.pathToNode.Insert(path, n)
	if !ok {
		return n, false
	}
	p.inodeToNode[p.maxInode] = n
	return n, true
}

func (p *inodePool) deleteNodeIf(path []string, filters ...func(node) bool) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	n, ok := p.pathToNode.Load(path)
	if !ok {
		return false
	}
	for _, filter := range filters {
		if !filter(n) {
			return false
		}
	}
	ok = p.pathToNode.Delete(path)
	if !ok {
		return false
	}
	delete(p.inodeToNode, n.inode)
	return true
}
