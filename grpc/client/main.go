package main

import (
	"flag"
	"fmt"
	pb "github.com/antonikonovalov/benches/grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

var (
	address = flag.String("addr", "0.0.0.0:4569", "addr of grpc service")
	msg     = flag.String("msg", "default message", "message")
)

func init() {
	flag.Parse()
}

func main() {

	// Set up first connection to the server.
	conn1, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn1.Close()

	// Set up second connection to the server.
	conn2, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		os.Exit(1)
	}
	defer conn2.Close()

	// first grpc client
	client1 := pb.NewCreatorClient(conn1)
	// second grpc client on the same conn
	client2 := pb.NewCreatorClient(conn1)
	// third client on another conn
	client3 := pb.NewCreatorClient(conn2)

	stream1, err := client1.CreateByStream(context.TODO())
	if err != nil {
		log.Fatalf("cant get first stream %v", err)
	}
	stream2, err := client2.CreateByStream(context.TODO())
	if err != nil {
		log.Fatalf("cant get second stream %v", err)
	}
	stream3, err := client3.CreateByStream(context.TODO())
	if err != nil {
		log.Fatalf("cant get third stream %v", err)
	}

	for i := 0; i < 100; i++ {
		log.Println("send: ", i)
		if err = stream1.Send(&pb.MessageRequest{Msg: fmt.Sprintf("stream1: number: %d", i)}); err != nil {
			log.Println(err)
		}
		if err = stream2.Send(&pb.MessageRequest{Msg: fmt.Sprintf("stream2: number: %d", i)}); err != nil {
			log.Println(err)
		}
		if err = stream3.Send(&pb.MessageRequest{Msg: fmt.Sprintf("stream3: number: %d", i)}); err != nil {
			log.Println(err)
		}
	}

	for i := 0; i < 100; i++ {
		log.Println("recv: ", i)
		resp, err := stream1.Recv()
		log.Println("stream1: response: ", resp)
		if err != nil {
			log.Println(err)
		}
		resp, err = stream2.Recv()
		log.Println("stream2: response: ", resp)
		if err != nil {
			log.Println(err)
		}
		resp, err = stream3.Recv()
		log.Println("stream3: response: ", resp)
		if err != nil {
			log.Println(err)
		}
	}
}
