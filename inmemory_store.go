package cache

import (
	"sync"
	"time"
)

type inMemoryStoreEntry struct {
	expiresAt time.Time
	value     any
}

// InMemoryStore A in memory key-value based cache store, implements [Store]
type InMemoryStore struct {
	cacheMap *sync.Map
}

var _ Store = (*InMemoryStore)(nil)

// NewInMemoryStore Creates a new usabled [InMemoryStore]
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		cacheMap: &sync.Map{},
	}
}

func (s *InMemoryStore) Set(key, value any, expiresAt time.Time) {
	entry := inMemoryStoreEntry{
		expiresAt: expiresAt,
		value:     value,
	}

	s.cacheMap.Store(key, entry)
}

func (s *InMemoryStore) Del(key any) {
	s.cacheMap.Delete(key)
}

func (s *InMemoryStore) Get(key any) (any, bool) {
	value, ok := s.cacheMap.Load(key)
	if !ok {
		return nil, false
	}

	entry, ok := value.(inMemoryStoreEntry)
	if !ok {
		// Remove non expirableEntry entries to avoid type colisions
		s.Del(key)
		return nil, false
	}

	if entry.expiresAt.Before(time.Now()) {
		s.Del(key)
		return nil, false
	}

	return entry.value, true
}
