package main

import (
	"github.com/piotrstrzalka/kvcache/pkg/kvcache"
	"github.com/piotrstrzalka/kvcache/pkg/service"
)

func main() {
	//load some config

	cache := kvcache.NewCache(kvcache.KvCacheConfig{})

	// service.RunRest(cache)
	service.RunGrpc(cache)
}
