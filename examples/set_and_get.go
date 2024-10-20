package main

import (
	"fmt"
	"time"

	"github.com/gustapinto/go-cache"
)

var store cache.Store

func main() {
	store = cache.NewInMemoryStore()

	store.Set("FooKey", "FooValue", time.Now().Add(5*time.Second))

	value, exists := store.Get("FooKey")
	fmt.Printf("value=%s, exists=%v\n", value, exists)
}
