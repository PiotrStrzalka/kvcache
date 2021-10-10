package kvcache

import (
	"errors"
	"log"
	"os"
	"time"
)

var expirationCheckerCycle = 5 * time.Second

type DataItem struct {
	data       []byte
	expiration time.Time
}

type KvCacheConfig struct {
	DefaultExpiryTime        time.Duration
	CacheSizeLimit           int
	DumpFilePath             string
	DeleteOldDataOnSizeLimit bool
}

type Kvcache struct {
	KvCacheConfig
	data      map[string]DataItem
	cacheSize int
	done      chan bool
}

func isValidPath(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return true
	}

	if err := os.WriteFile(path, nil, 0644); err != nil {
		return false
	}
	os.Remove(path)
	return true
}

func (k *Kvcache) updateDefaultConfig(cfg KvCacheConfig) {
	if cfg.CacheSizeLimit > 0 {
		k.CacheSizeLimit = cfg.CacheSizeLimit
	}

	if cfg.DefaultExpiryTime > 0 {
		k.DefaultExpiryTime = cfg.DefaultExpiryTime
	}

	if isValidPath(cfg.DumpFilePath) {
		k.DumpFilePath = cfg.DumpFilePath
	} else {
		log.Println("Cache path is not valid")
	}

	k.DeleteOldDataOnSizeLimit = cfg.DeleteOldDataOnSizeLimit
}

func NewCache(cfg KvCacheConfig) *Kvcache {
	kv := &Kvcache{
		data: make(map[string]DataItem),
		KvCacheConfig: KvCacheConfig{
			DefaultExpiryTime:        300 * time.Second,
			CacheSizeLimit:           1024 * 1024, // 1 MB
			DumpFilePath:             "",
			DeleteOldDataOnSizeLimit: true,
		},
	}
	kv.updateDefaultConfig(cfg)

	kv.done = make(chan bool)
	go func() {
		t := time.NewTicker(expirationCheckerCycle)
		for {
			select {
			case <-t.C:
				kv.removeExpiredData()
			case <-kv.done:
				t.Stop()
				return
			}
		}
	}()
	return kv
}

func StopCache(k *Kvcache) {
	k.done <- true
}

func (k *Kvcache) Set(key string, data []byte) error {
	return k.SetWithExpire(key, data, k.DefaultExpiryTime)
}

func (k *Kvcache) SetWithExpire(key string, data []byte, expire time.Duration) error {

	//TODO for muliple access some lock required
	oversize := (k.cacheSize + len(data)) > k.CacheSizeLimit

	if oversize && k.DeleteOldDataOnSizeLimit == false {
		return errors.New("No space in cache to save more data")
	}

	if oversize && k.DeleteOldDataOnSizeLimit == true {
		err := k.reclaimSpace(len(data))
		if err != nil {
			return err
		}
	}

	k.data[key] = DataItem{
		data:       data,
		expiration: time.Now().Add(expire),
	}
	k.cacheSize += len(data)

	log.Printf("Stored %q : %q , expire at: %v\n", key, string(data), k.data[key].expiration)
	return nil
}

func (k *Kvcache) Get(key string) ([]byte, error) {
	val, ok := k.data[key]
	if !ok {
		return nil, errors.New("Data not available")
	}

	log.Printf("Retrieved %q : %q\n", key, string(val.data))
	return val.data, nil
}

func (k *Kvcache) reclaimSpace(toReclaim int) error {
	if toReclaim > k.CacheSizeLimit {
		return errors.New("There is no enough space in the buffer")
	}

	var reclaimed int

	for reclaimed < toReclaim {
		reclaimed += k.deleteNearestExpirationData()
	}

	return nil
}

func (k *Kvcache) deleteNearestExpirationData() int {
	nearestToExpire := struct {
		time time.Time
		key  string
	}{
		time: time.Unix((1<<62)-1, 0),
		key:  "",
	}

	for key, val := range k.data {
		if nearestToExpire.time.After(val.expiration) {
			nearestToExpire.key = key
			nearestToExpire.time = val.expiration
		}
	}

	if nearestToExpire.key == "" {
		return 0
	}

	size := len(k.data[nearestToExpire.key].data)
	delete(k.data, nearestToExpire.key)
	log.Printf("Reclaiming space, data with key %s deleted, %d bytes gained\n", nearestToExpire.key, size)
	return size
}

func (k *Kvcache) removeExpiredData() {
	var deleted int
	now := time.Now()
	for key, val := range k.data {
		if now.After(val.expiration) {
			delete(k.data, key)
			deleted++
		}
	}

	log.Printf("Stale data clearing, removed items: %d\n", deleted)
}
