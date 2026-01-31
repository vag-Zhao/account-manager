package cache

import "time"

// ICache defines the interface for caching operations
type ICache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, duration time.Duration)
	Delete(key string)
	Flush()
}
