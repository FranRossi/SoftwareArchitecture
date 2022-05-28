package controllers

import (
	"VoteAPI/logic"
	"VoteAPI/proto/authService"
	pb "VoteAPI/proto/voteService"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VoterServer struct {
	server pb.VoteServiceServer
}

type AuthServer struct {
	server     proto.AuthServiceServer
	jwtManager *logic.JWTManager
}

var voteReply pb.VoteReply

func RegisterServicesServer(grpcServer *grpc.Server, jwtManager *logic.JWTManager) {
	voteServer := VoterServer{}
	pb.RegisterVoteServiceServer(grpcServer, &voteServer)

	server := AuthServer{jwtManager: jwtManager}
	proto.RegisterAuthServiceServer(grpcServer, &server)
}

func (server *AuthServer) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user, err := logic.RegisterUser(request.GetId(), request.GetUsername(), request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	user.Token = token
	_, err = logic.StoreUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot store user on data base: %v", err)
	}
	res := &proto.RegisterResponse{AccessToken: token}
	return res, nil
}

func (newVote *VoterServer) Vote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteReply, error) {
	user, err := logic.CheckVoter(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	if user == nil || !logic.IsCorrectPassword(user, req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	message := user.Username + " voted correctly"
	vote := &pb.VoteReply{Message: message}

	// rabbitmq test
	//api_voter.sendCertificate(idVoter)
	//api_voter.PrintCertificate()

	return vote, err
}
