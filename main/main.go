package main

import (
	"log"

	"github.com/piotrstrzalka/kvcache/pkg/kvcache"
)

func main() {
	cache := kvcache.NewCache(kvcache.KvCacheConfig{DumpFilePath: "./cache-dumped"})

	cache.Set("pio", []byte("programmer"))

	if val, err := cache.Get("pio"); err == nil {
		log.Printf("Data from cache: %s\n ", string(val))
	}

	//todo check if cache is deleted by GC after following line
	cache.Stop()

	cache2 := kvcache.NewCache(kvcache.KvCacheConfig{DumpFilePath: "./cache-dumped"})
	val, err := cache2.Get("pio")
	if err != nil {
		log.Fatalf("Cannot fetch pio from cache2")
	}

	log.Println("pio from cache2: ", val)

	cache2.Stop()
}
