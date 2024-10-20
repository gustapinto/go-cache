package cache

import (
	"time"
)

// Store Defines a interface for the implementation of concrete cache objects (aka. stores)
// in this library, it is strongly advised to use your own consumer-side interfaces
// instead of this
type Store interface {
	// Set Inserts or update a key-value entry into the store cache layer
	Set(key, value any, expiresAt time.Time)

	// Del Deletes a key-value entry from the store cache layer
	Del(key any)

	// Get Retuns a key-value entry from the store cache layer
	Get(key any) (value any, exists bool)
}
