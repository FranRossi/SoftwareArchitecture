package api_voter

import (
	"239850_221025_219401/api_voter/proto"
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
	name := req.GetName()
	//fmt.Println(name)
	sendCertificate(name)
	vote := &proto.VoteReply{Message: "Anduvo 1+1"}
	return vote, nil
}

const addr = "localhost:50004"

func Connection() {
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
