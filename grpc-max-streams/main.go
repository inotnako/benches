package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"fmt"
	pb "gitlab.nct/msg/grpc-exp/greetings"
	"gitlab.nct/msg/grpc-exp/service"
	"net/http"

	_ "net/http/pprof"
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

	// Set up a connection to the lookupd services

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go http.ListenAndServe(fmt.Sprintf(":%d", 36663), nil)
	pb.RegisterGreetingsServiceServer(grpcServer, service.New())
	grpcServer.Serve(lis)
}
