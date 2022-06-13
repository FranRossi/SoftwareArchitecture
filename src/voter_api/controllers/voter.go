package controllers

import (
	jwt "auth"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"voter_api/controllers/validation"
	"voter_api/domain"
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

var voteReply pb.VoteReply

func RegisterServicesServer(grpcServer *grpc.Server, jwtManager *jwt.Manager) {
	voteServer := VoterServer{}
	pb.RegisterVoteServiceServer(grpcServer, &voteServer)

	server := AuthServer{jwtManager: jwtManager}
	proto.RegisterAuthServiceServer(grpcServer, &server)
}

func (server *AuthServer) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	//user, err := checkVoter(request.GetId(), request.GetPassword())
	//if err != nil {
	//	return nil, err
	//}
	//
	//token, err := server.jwtManager.Generate(user.Username, user.Id, user.Role)
	//if err != nil {
	//	return nil, status.Errorf(codes.Internal, "cannot generate access token")
	//}
	//
	//res := &proto.LoginResponse{AccessToken: token}
	//return res, nil
	return &proto.LoginResponse{AccessToken: "Falso Token 1234"}, nil
}

//func checkVoter(id, password string) (*domain2.User, error) {
//	user, err := logic.FindVoter(id)
//	if err != nil {
//		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
//	}
//
//	if user == nil || !logic.IsCorrectPassword(user, password) {
//		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
//	}
//	return user, nil
//}

func (newVote *VoterServer) Vote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteReply, error) {
	timeFrontEnd := time.Now()
	voteModel := &domain.VoteModel{
		IdElection:  req.GetIdElection(),
		IdVoter:     req.GetIdVoter(),
		Circuit:     req.GetCircuit(),
		IdCandidate: req.GetIdCandidate(),
		Signature:   req.GetSignature(),
	}
	failed := verifyVote(voteModel)
	if failed != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid vote")
	}
	err := logic.StoreVote(voteModel)
	if err != nil {
		return &pb.VoteReply{Message: "Error"}, status.Errorf(codes.Internal, "cannot store vote: %v", err)
	}
	timeBackEnd := time.Now()
	if timeBackEnd.Sub(timeFrontEnd).Seconds() > 2 {
		logic.DeleteVote(voteModel)
		messageFailed := "vote cannot processed"
		return &pb.VoteReply{Message: messageFailed}, status.Errorf(codes.ResourceExhausted, "vote failed")
	} else {
		voteIdentification, err2 := logic.StoreVoteInfo(req.GetIdVoter(), req.GetIdElection(), timeFrontEnd, timeBackEnd)
		if err2 != nil {
			return &pb.VoteReply{Message: "Error"}, status.Errorf(codes.Internal, "cannot store vote info: %v", err)
		}
		go logic.SendCertificateSMS(voteModel, voteIdentification, timeFrontEnd)
		message := "voted correctly"
		return &pb.VoteReply{Message: message}, status.Errorf(codes.OK, "vote stored")
	}
}

func verifyVote(vote *domain.VoteModel) error {
	publicKeyPEM := validation.ReadKeyFromFile("./controllers/validation/pubkey.pem")
	publicKey := validation.ExportPEMStrToPubKey(publicKeyPEM)
	candidate := []byte(vote.IdCandidate)
	msgHash := sha256.New()
	msgHash.Write(candidate)
	msgHashSBytes := msgHash.Sum(nil)
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSBytes, vote.Signature, nil)
	if err != nil {
		return fmt.Errorf("verification failed")
	}
	return nil
}
