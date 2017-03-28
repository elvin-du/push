package main

import (
//	"hscore/log"
	"push/meta"

	"golang.org/x/net/context"
)

type Session struct {
//	ClientId string
//	UserId   string
}

func (*Session) Online(ctx context.Context, req *meta.SessionOnlineRequest) (*meta.SessionOnlineResponse, error) {
	return &meta.SessionOnlineResponse{}, nil
}

func (*Session) Offline(ctx context.Context, req *meta.SessionOfflineRequest) (*meta.SessionOfflineResponse, error) {
	return &meta.SessionOfflineResponse{}, nil
}

func (*Session) Update(ctx context.Context, req *meta.SessionUpdateRequest) (*meta.SessionUpdateResponse, error) {
	return &meta.SessionUpdateResponse{}, nil
}

func (*Session) Info(ctx context.Context, req *meta.SessionInfoRequest) (*meta.SessionInfoResponse, error) {
	return &meta.SessionInfoResponse{}, nil
}
