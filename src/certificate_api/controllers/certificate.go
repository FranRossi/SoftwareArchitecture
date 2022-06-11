package controllers

import (
	jwt "auth"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	domain2 "voter_api/domain/user"
	domain "voter_api/domain/vote"
	"voter_api/logic"
	proto "voter_api/proto/authService"
	pb "voter_api/proto/voteService"
)

type VoterServer struct {
	server pb.VoteServiceServer
}

type AuthServer struct {
	server     proto.AuthServiceServer
	jwtManager *jwt.Manager
}

var certReply pb.CertificateReply

func RegisterServicesServer(grpcServer *grpc.Server, jwtManager *jwt.Manager) {
	voteServer := VoterServer{}
	pb.RegisterVoteServiceServer(grpcServer, &voteServer)

	server := AuthServer{jwtManager: jwtManager}
	proto.RegisterAuthServiceServer(grpcServer, &server)
}

func checkVoter(id, password string) (*domain2.User, error) {
	user, err := logic.FindVoter(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !logic.IsCorrectPassword(user, password) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}
	return user, nil
}

func certificate(ctx context.Context, req *pb.CertificateRequest) (*pb.CertificateReply, error) {
	idVoter, err := checkCertificate(req)
	if err != nil {
		return nil, err
	}
	//return &pb.VoteReply{Message: message}, status.Errorf(codes.OK, "vote stored")

	// rabbitmq test
	api_voter.sendCertificate(idVoter)
	api_voter.PrintCertificate()
}

func checkCertificate(req *pb.CertificateRequest) (int, error) {
	user, err := logic.FindVoter(req.GetIdVoter())
	if err != nil {
		return nil, "", status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	return req.GetIdVoter(), nil
}
