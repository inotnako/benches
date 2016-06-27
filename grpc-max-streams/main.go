package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/antonikonovalov/benches/grpc-max-streams/greetings"
	"github.com/antonikonovalov/benches/grpc-max-streams/service"
	// _ "net/http/pprof"
)

var (
	addr = flag.String(`addr`, `localhost:4568`, `address for listen service`)
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpc.EnableTracing = false

	// Set up a connection to the lookupd services

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	// go http.ListenAndServe(fmt.Sprintf(":%d", 36663), nil)
	pb.RegisterGreetingsServiceServer(grpcServer, service.New())
	grpcServer.Serve(lis)
}
