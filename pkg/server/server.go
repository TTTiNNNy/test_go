package server

import (
	"context"

	. "challenge/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestRPCServer struct {
}

func (TestRPCServer) MakeShortLink(context.Context, *Link) (*Link, error) {
	return &Link{Data: "qwerty"}, nil
}
func (TestRPCServer) StartTimer(*Timer, ChallengeService_StartTimerServer) error {
	return status.Errorf(codes.Unimplemented, "method StartTimer not implemented")
}
func (TestRPCServer) ReadMetadata(context.Context, *Placeholder) (*Placeholder, error) {
	return &Placeholder{Data: "DataServerString"}, nil
}
