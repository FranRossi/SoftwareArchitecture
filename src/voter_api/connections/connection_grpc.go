package connections

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	l "own_logger"
)

func streamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	fmt.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, stream)
}

var connection net.Listener

func ConnectionGRPC() *grpc.Server {
	address := os.Getenv("grpc_address")
	conn, err := net.Listen("tcp", address)
	connection = conn
	if err != nil {
		l.LogError("tcp connection error: " + err.Error())
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),

	)
	fmt.Println("Starting gRPC server at: ", address)
	return grpcServer
}

func ServeGRPC(server *grpc.Server) {
	if err := server.Serve(connection); err != nil {
		l.LogError("gRPC server error: " + err.Error())
		log.Fatal(err)
	}
}
