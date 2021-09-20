package client

import "log"

type Grpc struct {
	addr string
}

func NewGrpc(addr string) *Grpc {
	return &Grpc{addr: addr}
}

func (r *Grpc) Get(key string) ([]byte, error) {
	log.Fatal("Not implemented")
	return nil, nil
}

func (r *Grpc) Set(key string, vue []byte) error {
	log.Fatal("Not plemented")
	return nil
}
