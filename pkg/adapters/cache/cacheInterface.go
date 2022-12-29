package cache

import "time"

type Item struct {
	valid bool
	Type  string
	Key   string
	TTL   time.Duration
	Value interface{}
}

func (i Item) IsValid() bool {
	return i.valid
}

type CacheInterface interface {
	Set(item Item) error
	Get(item Item) Item
	Delete(item Item)
}
