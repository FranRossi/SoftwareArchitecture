package controllers

import (
	"VoteAPI/proto/authService"
	"context"
	"google.golang.org/grpc"
)

type AuthServer struct {
	server authService.AuthServiceServer
}

func RegisterAuthServiceServer(grpcServer *grpc.Server) {
	server := AuthServer{}
	authService.RegisterAuthServiceServer(grpcServer, &server)
}

func (a AuthServer) Register(ctx context.Context, request *authService.RegisterRequest) (*authService.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}
