package greetings

import (
	"testing"

	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"sync"
	"time"
)

type server struct{}

var paylaod = []byte(`loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload
loadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadloadload`)

func (s *server) Talk(req *MsgRequest, stream GreetingsService_TalkServer) error {

	for i := 0; i < int(req.Count); i++ {
		err := stream.Send(&MsgResponse{
			Msg:     req.Msg,
			Payload: paylaod,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func startTestSrv() (func(), error) {
	lis, err := net.Listen("tcp", ":50045")
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	RegisterGreetingsServiceServer(s, &server{})
	up := make(chan struct{}, 0)
	go func() {
		close(up)
		if err = s.Serve(lis); err != nil {
			// for add information
			panic(err)
		}
	}()
	<-up

	return s.GracefulStop, nil
}

func TestBlockByReadFromStreams(t *testing.T) {
	stop, err := startTestSrv()
	if err != nil {
		t.Fatal(err)
	}
	defer stop()

	conn, err := grpc.Dial("localhost:50045", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	cli := NewGreetingsServiceClient(conn)

	startFirst := make(chan *struct{}, 0)
	wg := &sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// set default timeout
			timeout := 5 * time.Second

			if id == 1 {
				timeout = 10 * time.Second
			} else {
				<-startFirst
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			// open stream for reading
			wannaGot := 10000
			stream, err := cli.Talk(ctx, &MsgRequest{Msg: fmt.Sprintf(`msg-%d-issss`, id), Count: int32(wannaGot)})
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("open stream from id=%d", id)

			var counter int
			// not read - just wait
			if id == 1 {
				// close first before start other streams open
				close(startFirst)
				t.Logf("start sleep %s on task-%d", timeout.String(), id)
				time.Sleep(timeout)
			} else {
				// try read data from streams
				for {
					_, err := stream.Recv()
					if err == io.EOF {
						if counter != wannaGot {
							t.Errorf("expected counter eq %d, got - %d", wannaGot, counter)
						}
						t.Logf("got count msg %d n task-%d", counter, id)
						return
					}
					if err != nil {
						t.Error(err)
						return
					}
					counter++
					//t.Logf("got msg %s on task-%d", gotMsg.Msg, id)
				}
			}
		}(i)
	}

	wg.Wait()
}
