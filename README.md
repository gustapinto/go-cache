# go-cache

A simple cache with ttl and background expiration process


## Example Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/gustapinto/go-cache"
)

func main() {
	st := cache.NewInMemoryStore()

	st.Set("FooKey", "FooValue", time.Now().Add(5*time.Second))

	value, exists := st.Get("FooKey")
	fmt.Printf("value=%s, exists=%v", value, exists)
}
```