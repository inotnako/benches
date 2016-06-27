package benchmark

import (
	"google.golang.org/grpc"
	"testing"

	"golang.org/x/net/context"

	pb "github.com/antonikonovalov/benches/grpc-max-streams/greetings"
	//"io"
)

type Streamer interface {
	Recv() (*pb.MsgResponse, error)
}

func BenchmarkGrpcStreams(b *testing.B) {
	conn, err := grpc.Dial(`0.0.0.0:4568`, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreetingsServiceClient(conn)
	out := Resiver(b)

	//b.SetParallelism(100)
	b.ResetTimer()
	b.RunParallel(func(pbt *testing.PB) {
		for pbt.Next() {
			stream, err := c.Talk(context.Background())
			if err != nil {
				b.Errorf(`error grpc: %s`, err)
			}
			out <- stream
		}
	})
}

func Resiver(b *testing.B) chan Streamer {
	streamers := []Streamer{}
	in := make(chan Streamer)
	go func() {
		for {
			select {
			case stream := <-in:
				streamers = append(streamers, stream)
			}
			for i := range streamers {
				_, err := streamers[i].Recv()
				if err != nil {
					b.Errorf(`error stream grpc: %s`, err)
				}
			}
		}
	}()
	return in
}
