package connections

import (
	"VoteAPI/controllers"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func streamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, stream)
}

func ConnectionGRPC() {
	const addr = "localhost:50004"
	conn, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}

	grpcServer := grpc.NewServer()
	controllers.RegisterVoteServiceServer(grpcServer)
	controllers.RegisterAuthServiceServer(grpcServer)

	fmt.Println("Starting gRPC server at: ", addr)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
