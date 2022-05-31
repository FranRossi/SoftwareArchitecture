package connections

import (
	jwt "auth"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	"voter_api/controllers"
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

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func ConnectionGRPC() {
	const addr = "localhost:50004"
	conn, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}

	jwtManager := jwt.NewJWTManager(secretKey, tokenDuration)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
	)
	controllers.RegisterServicesServer(grpcServer, jwtManager)

	fmt.Println("Starting gRPC server at: ", addr)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
