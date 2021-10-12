package client

import (
	"context"
	"errors"
	"time"

	"github.com/piotrstrzalka/kvcache/pkg/model"
	"google.golang.org/grpc"
)

type Grpc struct {
	addr   string
	client model.KvCacheClient
}

func NewGrpc(addr string) (*Grpc, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	//defer conn.Close()
	grpc := &Grpc{addr: addr}
	grpc.client = model.NewKvCacheClient(conn)

	return grpc, nil
}

func (r *Grpc) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := r.client.Get(ctx, &model.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func (r *Grpc) Set(key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := r.client.Set(ctx, &model.SetRequest{Key: key, Value: value})
	if err != nil {
		return err
	}
	if res.Result == false {
		return errors.New("Problem with setting data")
	}
	return nil
}
