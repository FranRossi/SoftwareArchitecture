package controllers

import (
	jwt "auth"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"time"
	"voter_api/controllers/validation"
	"voter_api/domain"
	"voter_api/logic"
	proto "voter_api/proto/authService"
	pb "voter_api/proto/voteService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	message := "We received your vote, we will validate it shortly and send you a notification"
	go processVoteAndSendEmail(timeFrontEnd, req)
	return &pb.VoteReply{Message: message}, status.Errorf(codes.OK, "voted send for processing", nil)
}

func processVoteAndSendEmail(timeFrontEnd time.Time, req *pb.VoteRequest) {
	voteModel := &domain.VoteModel{
		IdElection:  req.GetIdElection(),
		IdVoter:     req.GetIdVoter(),
		Circuit:     req.GetCircuit(),
		IdCandidate: req.GetIdCandidate(),
		Signature:   req.GetSignature(),
	}
	voteIdentification, err := processVote(timeFrontEnd, voteModel)
	if err != nil {
		logic.SendCertificate(voteModel, voteIdentification, timeFrontEnd, err)
	}
	logic.SendCertificate(voteModel, voteIdentification, timeFrontEnd, nil)
}

func processVote(timeFrontEnd time.Time, voteModel *domain.VoteModel) (string, error) {
	failed := verifySignatureVote(voteModel)
	if failed != nil {
		_ = fmt.Errorf("signature verification failed: %v", failed)
		return "", failed
	}
	err := logic.StoreVote(voteModel)
	if err != nil {
		_ = fmt.Errorf("cannot store vote: %v", err)
		return "", err
	}
	timeBackEnd := time.Now()
	if timeBackEnd.Sub(timeFrontEnd).Seconds() > 2 {
		err2 := logic.DeleteVote(voteModel)
		if err2 != nil {
			_ = fmt.Errorf("cannot delete vote: %v", err2)
			return "", err2
		}
		messageFailed := "vote cannot processed under 2 seconds"
		return "", fmt.Errorf(messageFailed)
	} else {
		voteIdentification, err2 := logic.StoreVoteInfo(voteModel.IdVoter, voteModel.IdElection, timeFrontEnd, timeBackEnd)
		if err2 != nil {
			return "", fmt.Errorf("cannot store vote info: %v", err)
		}
		return voteIdentification, nil
	}
}

func verifySignatureVote(vote *domain.VoteModel) error {
	publicKeyPEM := validation.ReadKeyFromFile("./controllers/validation/pubkey.pem")
	publicKey := validation.ExportPEMStrToPubKey(publicKeyPEM)
	voter := []byte(vote.IdVoter)
	msgHash := sha256.New()
	msgHash.Write(voter)
	msgHashSBytes := msgHash.Sum(nil)
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSBytes, vote.Signature, nil)
	if err != nil {
		return fmt.Errorf("verification failed")
	}
	return nil
}
