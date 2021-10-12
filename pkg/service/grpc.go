package service

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/piotrstrzalka/kvcache/pkg/kvcache"
	"github.com/piotrstrzalka/kvcache/pkg/model"
	"google.golang.org/grpc"
)

type server struct {
	model.UnimplementedKvCacheServer
}

func RunGrpc(c *kvcache.Kvcache) error {
	lis, err := net.Listen("tcp", ":5010")
	if err != nil {
		return err
	}

	cache = c

	s := grpc.NewServer()
	model.RegisterKvCacheServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *server) Set(ctx context.Context, in *model.SetRequest) (*model.SetReply, error) {
	log.Printf("Received set request k: %v, v: %v\n", in.Key, in.Value)
	if err := cache.Set(in.Key, in.Value); err != nil {
		return &model.SetReply{Result: false}, err
	}
	return &model.SetReply{Result: true}, nil
}

func (s *server) SetWithExpire(ctx context.Context, in *model.SetRequest) (*model.SetReply, error) {
	log.Printf("Received set with expire request k: %v, v: %v, time: %v\n", in.Key, in.Value, time.Duration(in.ExpiryTime))
	if err := cache.SetWithExpire(in.Key, in.Value, time.Duration(in.ExpiryTime)); err != nil {
		return &model.SetReply{Result: false}, err
	}
	return &model.SetReply{Result: true}, nil
}

func (s *server) Get(ctx context.Context, in *model.GetRequest) (*model.GetResponse, error) {
	log.Printf("Received get request k: %v\n", in.Key)
	buf, err := cache.Get(in.Key)
	if err != nil {
		return nil, err
	}
	return &model.GetResponse{Key: in.Key, Value: buf}, nil
}
