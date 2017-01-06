package main

import (
	"log"
	"push/meta"

	"golang.org/x/net/context"
)

type Session struct {
}

func (*Session) Register(ctx context.Context, req *meta.SessionRegisterRequest) (*meta.SessionRegisterResponse, error) {
	log.Printf("user:%s login ip:%s", req.UserId, req.IP)
	return &meta.SessionRegisterResponse{}, nil
}

func (*Session) Unregister(ctx context.Context, req *meta.SessionUnregisterRequest) (*meta.SessionUnregisterResponse, error) {
	return &meta.SessionUnregisterResponse{}, nil
}
