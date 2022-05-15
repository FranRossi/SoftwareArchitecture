package api_voter

import (
	"239850_221025_219401/api-voter/data-access/repository"
	"239850_221025_219401/api-voter/proto"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var vote proto.VoteReply

type voterServer struct {
	proto.UnimplementedVoteServiceServer
}

func (newVote *voterServer) Vote(ctx context.Context, req *proto.VoteRequest) (*proto.VoteReply, error) {

  // rabbitmq test
  sendCertificate(name)
	PrintCertificate()

	idVoter := req.GetId()
	message, err := checkVoter(idVoter)
	vote := &proto.VoteReply{Message: message}
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

func Connection() {
	const addr = "localhost:50004"
	conn, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}

	grpcServer := grpc.NewServer()

	server := voterServer{}

	proto.RegisterVoteServiceServer(grpcServer, &server)

	fmt.Println("Starting gRPC server at: ", addr)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
