package fs

import "github.com/jacobsa/fuse/fuseops"

type InodePool struct {
	lastInode  fuseops.InodeID
	inodeToKey map[fuseops.InodeID]string
	keyToInode map[string]fuseops.InodeID
}

func NewInodePool() *InodePool {
	pool := &InodePool{
		lastInode:  fuseops.RootInodeID,
		inodeToKey: make(map[fuseops.InodeID]string),
		keyToInode: make(map[string]fuseops.InodeID),
	}
	pool.inodeToKey[fuseops.RootInodeID] = ""
	pool.keyToInode[""] = fuseops.RootInodeID
	return pool
}

func (p *InodePool) GetInodeFromKey(key string) fuseops.InodeID {
	if ino, ok := p.keyToInode[key]; ok {
		return ino
	}
	ino := p.lastInode + 1
	p.lastInode++

	p.keyToInode[key] = ino
	p.inodeToKey[ino] = key
	return ino
}

func (p *InodePool) GetKeyFromInode(ino fuseops.InodeID) (string, bool) {
	key, ok := p.inodeToKey[ino]
	return key, ok
}
