package cache

import (
	"sync"
	"time"
)

type Entry struct {
	ExpiresAt time.Time
	Value     any
}

type Store struct {
	cacheMap *sync.Map
}

func NewStore() *Store {
	store := &Store{}
	if store.cacheMap == nil {
		store.cacheMap = &sync.Map{}
	}

	go func() {
		for {
			store.cacheMap.Range(func(key, value any) bool {
				entry, ok := value.(Entry)
				if !ok {
					store.Del(key)
					return true
				}

				if entry.ExpiresAt.Before(time.Now()) {
					store.Del(key)
				}

				return true
			})

			time.Sleep(500 * time.Millisecond)
		}
	}()

	return store
}

func (c *Store) Set(key, value any, expiresAt time.Time) {
	entry := Entry{
		ExpiresAt: expiresAt,
		Value:     value,
	}

	c.cacheMap.Store(key, entry)
}

func (c *Store) Del(key any) {
	c.cacheMap.Delete(key)
}

func (c *Store) Get(key any) (any, bool) {
	value, ok := c.cacheMap.Load(key)
	if !ok {
		return nil, false
	}

	entry, ok := value.(Entry)
	if !ok {
		c.Del(key)
		return nil, false
	}

	if entry.ExpiresAt.Before(time.Now()) {
		c.Del(key)
		return nil, false
	}

	return entry.Value, true
}
