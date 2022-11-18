package main

import (
	"context"
	"io"
	"log"
	"net/http"

	pb "github.com/sjbarag/grpc-go-500-response-mwe/helloworld"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	mux := http.DefaultServeMux

	// Add a simple HTTP-only endpoint to confirm HTTP functionality.
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	}))

	// Add the Greeter service under the /greet/ path prefix.
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	mux.Handle("/greet/", s)

	log.Printf("server listening on :55123")
	log.Fatal(http.ListenAndServe(":55123", nil))
}
