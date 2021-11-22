package kvcache_test

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	kv "github.com/piotrstrzalka/kvcache/pkg/kvcache"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c := kv.NewCache(kv.KvCacheConfig{
		DefaultExpiryTime: 5 * time.Minute,
	})
	key := "height"
	data := []byte("128")

	c.Set(key, data)

	r, err := c.Get("height")
	if err != nil {
		t.Fatalf("Error while reading data: %s", err.Error())
	}

	if cmp := bytes.Compare(data, r); cmp != 0 {
		t.Errorf("Received value is different tham stored, stored: %q , received: %q", string(data), string(r))
	}
}

func TestDefaultExpiration(t *testing.T) {
	c := kv.NewCache(kv.KvCacheConfig{
		DefaultExpiryTime:      50 * time.Millisecond,
		ExpirationCheckerCycle: 5 * time.Millisecond,
	})

	err := c.Set("tkey", []byte("tdata"))
	if err != nil {
		t.Fatal(err)
	}
	//how to mock time??
	_, err = c.Get("tkey")
	if err != nil {
		t.Fatalf("Data shall be still available")
	}

	time.Sleep(60 * time.Millisecond)
	_, err = c.Get("tkey")
	if err == nil {
		t.Fatalf("Data shall be deleted by this time")
	}
}

func TestSizeLimitWithDataDeletion(t *testing.T) {
	c := kv.NewCache(kv.KvCacheConfig{CacheSizeLimit: 64, DeleteOldDataOnSizeLimit: true})
	for i := 0; i < 10; i++ {
		c.SetWithExpire(fmt.Sprintf("data%d", i), []byte("eightbyt"), (time.Duration(i) * time.Hour))
	}

}

func TestSizeLimitWithOutDataDeletion(t *testing.T) {
	c := kv.NewCache(kv.KvCacheConfig{CacheSizeLimit: 64, DeleteOldDataOnSizeLimit: false})
	for i := 0; i < 8; i++ {
		c.SetWithExpire(fmt.Sprintf("data%d", i), []byte("eightbyt"), (time.Duration(i) * time.Hour))
	}

	err := c.SetWithExpire("data", []byte("eightbyt"), time.Hour)
	if err == nil {
		t.Error("Data should not fit into cache")
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestStoringData(t *testing.T) {
	dataToStore := "programmer"
	path := "./test-" + RandStringRunes(5)
	c := kv.NewCache(kv.KvCacheConfig{DumpFilePath: path})
	if err := c.Set("pio", []byte(dataToStore)); err != nil {
		log.Fatalf("Cannot store the data, %v\n", err.Error())
	}

	assert.Equal(t, nil, c.Stop())

	defer os.Remove(path)

	cache2 := kv.NewCache(kv.KvCacheConfig{DumpFilePath: path})
	val, err := cache2.Get("pio")
	if err != nil {
		t.Fatalf("Cannot fetch pio from cache2")
	}

	if string(val) != dataToStore {
		t.Fatalf("Bad data fetch after store, wanted: %s, got: %s\n", dataToStore, string(val))
	}
}

//TestRaceDetector should be run with -race (eg. go test -race -run=TestRaceDetector )
func TestRaceDetector(t *testing.T) {
	starter := make(chan bool)
	var wg sync.WaitGroup

	c := kv.NewCache(kv.KvCacheConfig{})

	for i := 0; i < 100; i++ {
		loci := i
		go func() {
			wg.Add(1)
			defer wg.Done()
			<-starter
			key := fmt.Sprintf("pio%d", loci)
			value := []byte(fmt.Sprintf("data%d", loci))
			c.Set(key, value)

			fetched, err := c.Get(key)
			assert.Equal(t, nil, err)
			assert.Equal(t, value, fetched)
		}()
	}

	close(starter)
	wg.Wait()
}
