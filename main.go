package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbHelloWorld "github.com/iamrajiv/helloworld-grpc-gateway/proto/helloworld"
)

// Server struct representing our service implementation
type server struct{}

// SayHello is the implementation of the SayHello method defined in the proto file
func (*server) SayHello(_ context.Context, in *pbHelloWorld.HelloRequest) (*pbHelloWorld.HelloReply, error) {
	return &pbHelloWorld.HelloReply{Message: in.Name + " world"}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pbHelloWorld.RegisterGreeterServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// Create a new ServeMux for the gRPC-Gateway
	gwmux := runtime.NewServeMux()
	// Register the Greeter service with the gRPC-Gateway
	err = pbHelloWorld.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// Create a new HTTP server for the gRPC-Gateway
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())

}
