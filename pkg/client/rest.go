package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Rest struct {
	addr string
}

func NewRest(addr string) (*Rest, error) {
	return &Rest{addr: addr}, nil
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *Rest) Get(key string) ([]byte, error) {

	endpoint := "http://" + r.addr + "/get/" + "?key=" + key

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return buf.Bytes(), nil
}

func (r *Rest) Set(key string, value []byte) error {
	endpoint := "http://" + r.addr + "/set/"

	buf := bytes.NewBuffer(nil)
	kv := KeyValue{
		Key:   key,
		Value: string(value),
	}
	json.NewEncoder(buf).Encode(kv)

	//todo I don't like this raw bytes sending
	resp, err := http.Post(endpoint, "", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Wrong status code retrieved: %d", resp.StatusCode))
	}
	return nil
}
