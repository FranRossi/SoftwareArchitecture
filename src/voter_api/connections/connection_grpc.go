package connections

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
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

var connection net.Listener

func ConnectionGRPC() *grpc.Server {
	address := os.Getenv("grpc_address")
	conn, err := net.Listen("tcp", address)
	connection = conn
	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
	)

	fmt.Println("Starting gRPC server at: ", address)
	return grpcServer
}

func ServeGRPC(server *grpc.Server) {
	if err := server.Serve(connection); err != nil {
		log.Fatal(err)
	}
}
