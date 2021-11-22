package kvcache

import (
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/braintree/manners"
)

type DataItem struct {
	Data       []byte    `json:"data"`
	Expiration time.Time `json:"expiration"`
}

type KvCacheConfig struct {
	DefaultExpiryTime        time.Duration
	ExpirationCheckerCycle   time.Duration
	CacheSizeLimit           int
	DumpFilePath             string
	DeleteOldDataOnSizeLimit bool
	CacheStatus              bool
}

type Kvcache struct {
	KvCacheConfig
	data      map[string]DataItem
	mutex     sync.RWMutex
	cacheSize int
	done      chan bool
	monitor   *manners.GracefulServer
}

var (
	eCacheSize                *expvar.Int
	eCacheSizeOccupied        *expvar.Int
	eCacheSizeOccupiedPercent *expvar.Int
)

//NewCache inits new Kvcache instance, validates and applies configuration from
//input struct cfg. Depending on settings starts Cache status server.
//Returns pointer to initialized instance of Kvcache
func NewCache(cfg KvCacheConfig) *Kvcache {
	kv := &Kvcache{
		data: make(map[string]DataItem),
		KvCacheConfig: KvCacheConfig{
			DefaultExpiryTime:        87600 * time.Hour,
			ExpirationCheckerCycle:   5 * time.Second,
			CacheSizeLimit:           1024 * 1024, // 1 MB
			DumpFilePath:             "",
			DeleteOldDataOnSizeLimit: true,
		},
	}
	kv.updateDefaultConfig(cfg)

	if kv.DumpFilePath != "" {
		if err := kv.fetchDataFromDisk(); err != nil {
			log.Printf("Cannot fetch data from disk: %v\n", err.Error())
		}
	}

	if cfg.CacheStatus == true {
		kv.startHealthMonitor()
	}

	kv.done = make(chan bool)
	go func() {
		t := time.NewTicker(kv.ExpirationCheckerCycle)
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

//Stop stores cache data on disk
//stops cache status server
//and closes cache clearing goroutine
func (k *Kvcache) Stop() error {
	if k.DumpFilePath != "" {
		if err := k.dumpDataToDisk(); err != nil {
			return errors.New(fmt.Sprintf("Cannot dump data to the disk %v", err.Error()))
		}
	}

	if k.monitor != nil {
		k.monitor.Close()
	}

	k.done <- true
	return nil
}

//Set stores data in cache to expire after default time
func (k *Kvcache) Set(key string, data []byte) error {
	return k.SetWithExpire(key, data, k.DefaultExpiryTime)
}

//Set stores data in cache to expire after specific time
func (k *Kvcache) SetWithExpire(key string, data []byte, expire time.Duration) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()
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

	k.cacheSize += len(data)

	if _, exists := k.data[key]; exists == true {
		k.cacheSize -= len(k.data[key].Data)
	}

	k.data[key] = DataItem{
		Data:       data,
		Expiration: time.Now().Add(expire),
	}

	log.Printf("Stored %q : %q , expire at: %v\n", key, string(data), k.data[key].Expiration)
	return nil
}

//Get fetches data from cache
func (k *Kvcache) Get(key string) ([]byte, error) {
	k.mutex.RLock()
	val, ok := k.data[key]
	k.mutex.RUnlock()
	if !ok {
		return nil, errors.New("Data not available")
	}

	log.Printf("Retrieved %q : %q\n", key, string(val.Data))
	return val.Data, nil
}

func isValidPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
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

	if cfg.ExpirationCheckerCycle > 0 {
		k.ExpirationCheckerCycle = cfg.ExpirationCheckerCycle
	}

	if isValidPath(cfg.DumpFilePath) {
		k.DumpFilePath = cfg.DumpFilePath
	} else {
		log.Println("Cache path is not valid")
	}

	k.DeleteOldDataOnSizeLimit = cfg.DeleteOldDataOnSizeLimit
}

func (k *Kvcache) startHealthMonitor() {
	eCacheSize = expvar.NewInt("cachememory.size")
	eCacheSizeOccupied = expvar.NewInt("cachememory.occupied")
	eCacheSizeOccupiedPercent = expvar.NewInt("cachememory.occupiedpercentage")

	eCacheSize.Set(int64(k.CacheSizeLimit))
	eCacheSizeOccupied.Set(0)
	eCacheSizeOccupiedPercent.Set(0)
	log.Println("status available at http://localhost:5011/metrics")

	host := "0.0.0.0:5011"

	s := manners.NewWithServer(&http.Server{
		Addr:           host,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})

	go s.ListenAndServe()
	k.monitor = s
}

func (k *Kvcache) dumpDataToDisk() error {
	f, err := os.Create(k.DumpFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(k.data); err != nil {
		return err
	}
	return nil
}

func (k *Kvcache) fetchDataFromDisk() error {
	f, err := os.Open(k.DumpFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&k.data); err != nil {
		return err
	}
	return nil
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
		if nearestToExpire.time.After(val.Expiration) {
			nearestToExpire.key = key
			nearestToExpire.time = val.Expiration
		}
	}

	if nearestToExpire.key == "" {
		return 0
	}

	size := len(k.data[nearestToExpire.key].Data)
	delete(k.data, nearestToExpire.key)
	log.Printf("Reclaiming space, data with key %s deleted, %d bytes gained\n", nearestToExpire.key, size)
	return size
}

func (k *Kvcache) removeExpiredData() {
	var deleted int
	now := time.Now()
	for key, val := range k.data {
		if now.After(val.Expiration) {
			delete(k.data, key)
			deleted++
		}
	}

	if k.CacheStatus == true {
		eCacheSizeOccupiedPercent.Set(int64(k.cacheSize * 100 / k.CacheSizeLimit))
		eCacheSizeOccupied.Set(int64(k.cacheSize))
	}

	log.Printf("Stale data clearing, removed items: %d\n", deleted)
}
