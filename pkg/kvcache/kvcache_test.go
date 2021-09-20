package kvcache

import (
	"bytes"
	"testing"
)

func TestGet(t *testing.T) {
	c := NewCache()
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
