package kvstore

import (
	"encoding/json"
	"sync"
)

type KVStore interface {
	Watch(uuid string) (watchCh <-chan bool, cancel func())
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	Commit(commands ...string)
	Get(key string) Entry
	Snapshot() map[string]Entry
}

func New() KVStore {
	return &kvstore{
		mu:         sync.RWMutex{},
		store:      map[string]Entry{},
		watchChMap: sync.Map{},
	}
}

type Entry struct {
	Version uint64 `json:"version"`
	Value   string `json:"value"`
}

type kvstore struct {
	mu    sync.RWMutex
	store map[string]Entry
	// private
	watchChMap sync.Map // map[uuid]chan bool
}

func (s *kvstore) Snapshot() map[string]Entry {
	out := map[string]Entry{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.store {
		out[k] = v
	}
	return out
}
func (s *kvstore) Get(key string) Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, loaded := s.store[key]; loaded {
		return e
	}
	return Entry{}
}

func (s *kvstore) Watch(uuid string) (watchCh <-chan bool, cancel func()) {
	ch := make(chan bool, 1)
	s.watchChMap.Store(uuid, ch)
	return ch, func() {
		s.watchChMap.Delete(uuid)
	}
}

func (s *kvstore) Marshal() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.store)
}
func (s *kvstore) Unmarshal(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(b, &s.store)
}

func (s *kvstore) Commit(commands ...string) {
	for _, command := range commands {
		c := &Command{}
		if err := c.Decode(command); err != nil {
			continue
		}
		var watchCh chan bool = nil
		if val, loaded := s.watchChMap.LoadAndDelete(c.Uuid); loaded {
			watchCh = val.(chan bool)
		}

		s.mu.Lock()
		var version uint64 = 0
		if e, loaded := s.store[c.Key]; loaded {
			version = e.Version
		}
		consistent := c.Version >= version+1
		if consistent {
			switch c.Operation {
			case OpSet:
				s.store[c.Key] = Entry{
					Version: c.Version,
					Value:   c.Value,
				}
			case OpDel:
				delete(s.store, c.Key)
			}
		}
		s.mu.Unlock()

		if watchCh != nil {
			watchCh <- consistent
			close(watchCh)
		}
	}
}
