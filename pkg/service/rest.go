package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/piotrstrzalka/kvcache/pkg/kvcache"
)

func init() {
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/set/", setHandler)
}

var cache *kvcache.Kvcache

func Run(c *kvcache.Kvcache) {
	host := "0.0.0.0:5010"

	cache = c

	s := http.Server{
		Addr:         host,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.DefaultServeMux,
	}

	log.Printf("Listening on: %s\n", host)
	s.ListenAndServe()
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d, err := cache.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write(d)
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	var kv KeyValue

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received data %v", string(buf))

	err = json.NewDecoder(bytes.NewBuffer(buf)).Decode(&kv)
	if err != nil {
		log.Printf("wrong data received, %q", kv)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = cache.Set(kv.Key, []byte(kv.Value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
