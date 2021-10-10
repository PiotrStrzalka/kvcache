package kvcache

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	c := NewCache(KvCacheConfig{
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
	expirationCheckerCycle = 5 * time.Millisecond
	c := NewCache(KvCacheConfig{DefaultExpiryTime: 50 * time.Millisecond})

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
	expirationCheckerCycle = 5 * time.Second
}

func TestSizeLimitWithDataDeletion(t *testing.T) {
	c := NewCache(KvCacheConfig{CacheSizeLimit: 64, DeleteOldDataOnSizeLimit: true})
	for i := 0; i < 10; i++ {
		c.SetWithExpire(fmt.Sprintf("data%d", i), []byte("eightbyt"), (time.Duration(i) * time.Hour))
	}

}

func TestSizeLimitWithOutDataDeletion(t *testing.T) {
	c := NewCache(KvCacheConfig{CacheSizeLimit: 64, DeleteOldDataOnSizeLimit: false})
	for i := 0; i < 8; i++ {
		c.SetWithExpire(fmt.Sprintf("data%d", i), []byte("eightbyt"), (time.Duration(i) * time.Hour))
	}

	err := c.SetWithExpire("data", []byte("eightbyt"), time.Hour)
	if err == nil {
		t.Error("Data should not fit into cache")
	}
}
