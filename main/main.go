package main

import (
	"fmt"

	"github.com/piotrstrzalka/kvcache/pkg/kvcache"
)

func main() {
	cache := kvcache.NewCache()

	cache.Set("pio", []byte("programmer"))

	if val, err := cache.Get("pio"); err == nil {
		fmt.Printf("Data from cache: %s", string(val))
	}
}
