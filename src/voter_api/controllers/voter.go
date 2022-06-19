package controllers

import (
	jwt "auth/jwt"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	l "own_logger"
	"time"
	"voter_api/controllers/validation"
	"voter_api/domain"
	"voter_api/logic"
	proto "voter_api/proto/authService"
	pb "voter_api/proto/voteService"

	"encrypt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VoterServer struct {
	server     pb.VoteServiceServer
	jwtManager *jwt.Manager
}

type AuthServer struct {
	server     proto.AuthServiceServer
	jwtManager *jwt.Manager
}

var voteReply pb.VoteReply

func RegisterServicesServer(grpcServer *grpc.Server, jwtManager *jwt.Manager) {
	voteServer := VoterServer{jwtManager: jwtManager}
	pb.RegisterVoteServiceServer(grpcServer, &voteServer)
}

func (newVote *VoterServer) Vote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteReply, error) {
	timeFrontEnd := time.Now()
	message := "We received your vote, we will validate it shortly and send you a notification"
	go processVoteAndSendEmail(timeFrontEnd, req)
	return &pb.VoteReply{Message: message}, status.Errorf(codes.OK, "voted send for processing", nil)
}

func processVoteAndSendEmail(timeFrontEnd time.Time, req *pb.VoteRequest) {
	voteModel := domain.VoteModel{
		IdElection:  req.GetIdElection(),
		IdVoter:     req.GetIdVoter(),
		Circuit:     req.GetCircuit(),
		IdCandidate: req.GetIdCandidate(),
		Signature:   req.GetSignature(),
	}
	encrypt.DecryptVote((*encrypt.VoteModel)(&voteModel))
	voteIdentification, err := processVote(timeFrontEnd, voteModel)
	if err != nil {
		l.LogError(err.Error())
		fmt.Println(err.Error())
		logic.SendCertificate(voteModel, voteIdentification, timeFrontEnd, err)
	}
	l.LogInfo("Vote processed")
	logic.SendCertificate(voteModel, voteIdentification, timeFrontEnd, nil)
}

func processVote(timeFrontEnd time.Time, voteModel domain.VoteModel) (string, error) {
	failed := verifySignatureVote(voteModel)
	if failed != nil {
		return "", failed
	}
	err := logic.StoreVote(voteModel)
	if err != nil {
		return "", err
	}
	timeBackEnd := time.Now()
	if timeBackEnd.Sub(timeFrontEnd).Seconds() > 2 {
		err2 := logic.DeleteVote(voteModel)
		if err2 != nil {
			return "", fmt.Errorf("cannot delete vote that was processed over 2 seconds: %v", err2)
		}
		messageFailed := "vote cannot processed under 2 seconds and was deleted"
		return "", fmt.Errorf(messageFailed)
	} else {
		voteIdentification, err2 := logic.StoreVoteInfo(voteModel.IdVoter, voteModel.IdElection, timeFrontEnd, timeBackEnd)
		if err2 != nil {
			return "", fmt.Errorf("cannot store vote info: %v", err2)
		}
		return voteIdentification, nil
	}
}

func verifySignatureVote(vote domain.VoteModel) error {
	publicKeyPEM := validation.ReadKeyFromFile("./controllers/validation/pubkey.pem")
	publicKey := validation.ExportPEMStrToPubKey(publicKeyPEM)
	voter := []byte(vote.IdVoter)
	msgHash := sha256.New()
	msgHash.Write(voter)
	msgHashSBytes := msgHash.Sum(nil)
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSBytes, vote.Signature, nil)
	if err != nil {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}
