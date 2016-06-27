package benchmark

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"testing"

	pb "github.com/antonikonovalov/benches/grpc-max-streams/greetings"
	//"io"
)

type Streamer interface {
	Recv() (*pb.MsgResponse, error)
	CloseSend() error
}

func BenchmarkGrpcStreams(b *testing.B) {
	conn, err := grpc.Dial(`0.0.0.0:4568`, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreetingsServiceClient(conn)
	out, done, exited := Resiver(b)

	//b.SetParallelism(100)
	b.ResetTimer()
	b.RunParallel(func(pbt *testing.PB) {
		// println(`NNNN`, b.N)
		for pbt.Next() {
			// println(`i`, i)

			stream, err := c.Talk(context.Background())
			if err != nil {
				b.Errorf(`error grpc: %s`, err)
			}
			out <- stream
		}
	})

	close(done)

	<-exited
}

func Resiver(b *testing.B) (chan Streamer, chan struct{}, chan struct{}) {
	streamers := []Streamer{}
	in := make(chan Streamer)
	done := make(chan struct{})
	exited := make(chan struct{})
	go func() {
		for {
			select {
			case stream := <-in:
				streamers = append(streamers, stream)
			case <-done:
				println(`streams count:`, len(streamers))

				for i := range streamers {
					var err error

					_, err = streamers[i].Recv()
					if err != nil {
						println(`error stream grpc: %s`, err)
					}

					if err != io.EOF {
						if err := streamers[i].CloseSend(); err != nil {
							println(`error stream %d on close: %s`, i, err)
						}
					}
				}
				exited <- struct{}{}
				return
			}
		}
	}()
	return in, done, exited
}
