package kvcache

import (
	"errors"
	"log"
)

type Kvcache struct {
	data map[string][]byte
}

func NewCache() *Kvcache {
	return &Kvcache{
		data: make(map[string][]byte),
	}
}

func (k *Kvcache) Set(key string, data []byte) error {
	k.data[key] = data
	log.Printf("Stored %q : %q\n", key, string(data))
	return nil
}

func (k *Kvcache) Get(key string) ([]byte, error) {
	val, ok := k.data[key]
	if !ok {
		return nil, errors.New("Data not available")
	}

	log.Printf("Retrieved %q : %q\n", key, string(val))
	return val, nil
}
