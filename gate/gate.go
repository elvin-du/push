package main

import (
	"log"
	"push/meta"

	"golang.org/x/net/context"
)

type Gate struct {
}

func (*Gate) Push(ctx context.Context, req *meta.GatePushRequest) (*meta.GatePushResponse, error) {
	log.Printf("userId:%s Msg:%s", req.UserId, req.Msg)
	return &meta.GatePushResponse{}, nil
}

