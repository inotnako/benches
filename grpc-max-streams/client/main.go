package main

import (
	"github.com/antonikonovalov/benches/grpc-max-streams/greetings"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"time"
)

func main() {
	// Set up a connection to the lookupd services
	conn, err := grpc.Dial(`localhost:4568`, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("failed to listen lookupd: %v", err)
	}

	defer conn.Close()
	ctx := context.Background()
	stream, err := greetings.NewGreetingsServiceClient(conn).Talk(ctx)
	if err != nil {
		panic(err)
	}

	for i := 0; ; i++ {
		msg, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		println(`iter:`, i)
		println(`msg:`, msg.Msg)
		if i == 10 {
			time.Sleep(time.Second * 3)
		}
	}
}
