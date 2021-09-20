package main

import (
	"github.com/piotrstrzalka/kvcache/pkg/kvcache"
	rest "github.com/piotrstrzalka/kvcache/pkg/service"
)

func main() {
	//load some config

	cache := kvcache.NewCache()

	rest.Run(cache)
}
