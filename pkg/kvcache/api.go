package kvcache

import "time"

type KvApi interface {
	Set(key string, value []byte) error
	SetWithExpiry(key string, value []byte, expire time.Duration)
	Get(key string) ([]byte, error)
	Delete(key string) error
	Exists(key string) bool
}
