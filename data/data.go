package main

import (
	"log"
	"push/meta"

	"golang.org/x/net/context"
)

type Data struct {
}

func (*Data) Online(ctx context.Context, req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	log.Printf("user:%s login ip:%s", req.UserId, req.IP)
	header := &meta.ResponseHeader{Code: 0, Msg: "success"}
	return &meta.DataOnlineResponse{Header: header}, nil

}

func (*Data) Offline(ctx context.Context, req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	log.Printf("user:%s offline", req.UserId)
	return &meta.DataOfflineResponse{}, nil

}
