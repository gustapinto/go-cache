package cache

import (
	"time"
)

type Entry struct {
	ExpiresAt time.Time
	Value     any
}

type Store interface {
	Set(key any, value any, expiresAt time.Time)

	Del(key any)

	Get(key any) (value any, exists bool)
}
