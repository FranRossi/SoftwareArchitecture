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
	pb "voter_api/proto/voteService"

	"encrypt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VoterServer struct {
	server     pb.VoteServiceServer
	jwtManager *jwt.Manager
	channel    chan VoteAndTime
}

type VoteAndTime struct {
	Vote         *pb.VoteRequest
	TimeFrontEnd time.Time
}

func RegisterServicesServer(grpcServer *grpc.Server, jwtManager *jwt.Manager) *VoterServer {
	voteServer := VoterServer{
		jwtManager: jwtManager,
		channel:    make(chan VoteAndTime),
	}
	pb.RegisterVoteServiceServer(grpcServer, &voteServer)
	return &voteServer
}

func (newVote *VoterServer) Vote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteReply, error) {
	message := "We received your vote, we will validate it shortly and send you a notification"
	var voteAndTime VoteAndTime
	voteAndTime.Vote = req
	voteAndTime.TimeFrontEnd = time.Now()
	newVote.channel <- voteAndTime
	return &pb.VoteReply{Message: message}, status.Errorf(codes.OK, "voted send for processing", nil)
}

func ActivateChannel(server *VoterServer) {
	for i := 0; i < 50000; i++ {
		go processVotes(server)
	}
}

func processVotes(server *VoterServer) {
	for {
		var vote VoteAndTime
		vote = <-server.channel
		processVoteAndSendEmail(vote.Vote, vote.TimeFrontEnd)
	}
}

func processVoteAndSendEmail(req *pb.VoteRequest, timeFrontEnd time.Time) {
	voteModel := domain.VoteModel{
		IdElection:  req.IdElection,
		IdVoter:     req.IdVoter,
		Circuit:     req.Circuit,
		IdCandidate: req.IdCandidate,
		Signature:   req.Signature,
	}
	encrypt.DecryptVote((*encrypt.VoteModel)(&voteModel))
	voteIdentification, err := processVote(timeFrontEnd, voteModel)
	if err != nil {
		go l.LogError(err.Error())
		fmt.Println(err.Error())
		go logic.SendCertificate(voteModel, voteIdentification, timeFrontEnd, err)
	}
	go l.LogInfo("Vote processed")
	go logic.SendCertificate(voteModel, voteIdentification, timeFrontEnd, nil)
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
	//if timeBackEnd.Sub(timeFrontEnd).Seconds() > 2 {
	//	err2 := logic.DeleteVote(voteModel)
	//	if err2 != nil {
	//		return "", fmt.Errorf("cannot delete vote that was processed over 2 seconds: %v", err2)
	//	}
	//	messageFailed := "vote cannot processed under 2 seconds"
	//	return "", fmt.Errorf(messageFailed)
	//} else {
	fmt.Println(timeBackEnd.Sub(timeFrontEnd).Seconds())
	voteIdentification, err2 := logic.StoreVoteInfo(voteModel.IdVoter, voteModel.IdElection, timeFrontEnd, timeBackEnd)
	if err2 != nil {
		return "", fmt.Errorf("cannot store vote info: %v", err2)
	}
	return voteIdentification, nil
	//}
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
