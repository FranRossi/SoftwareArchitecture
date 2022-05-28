package controllers

import (
	"VoteAPI/data_access/repository"
	"VoteAPI/proto/voteService"
	"context"
	"google.golang.org/grpc"
	"log"
)

type VoterServer struct {
	server voteService.VoteServiceServer
}

var voteReply voteService.VoteReply

func RegisterVoteServiceServer(grpcServer *grpc.Server) {
	server := VoterServer{}
	voteService.RegisterVoteServiceServer(grpcServer, &server)
}

func (newVote *VoterServer) Vote(ctx context.Context, req *voteService.VoteRequest) (*voteService.VoteReply, error) {

	idVoter := req.GetId()
	//passwordVoter := req.GetPassword()

	//token, err := JwtAuth(idVoter, passwordVoter)

	//voting, aca irian las go routines no?
	message, err := checkVoter(idVoter)
	vote := &voteService.VoteReply{Message: message}

	// rabbitmq test
	//api_voter.sendCertificate(idVoter)
	//api_voter.PrintCertificate()

	return vote, err
}

func checkVoter(idVoter string) (string, error) {
	_, err := repository.CheckVoterId(idVoter)
	if err != nil {
		log.Fatal(err)
		return "Error chequeando", err
	}
	return "Chequeo bien el id", err
}
