package service

import (
	pb "github.com/antonikonovalov/benches/grpc-max-streams/greetings"
	"strconv"
	"sync/atomic"
)

func New() pb.GreetingsServiceServer {
	return new(server)
}

type server struct{}

var ops uint64 = 0

func (s *server) Talk(conversation pb.GreetingsService_TalkServer) error {
	count := atomic.AddUint64(&ops, 1)
	if count%20 == 0 {
		print(", ", count, "\n")
	} else {
		print(", ", count)
	}

	for i := 0; i <= 100000; i++ {
		//time.Sleep(100 * time.Millisecond)
		err := conversation.Send(&pb.MsgResponse{`test-test-` + strconv.Itoa(i)})
		if err != nil {
			println(`ERR`, err.Error())
			return err
		}
		//println(`sended msg - `, i)
	}

	return nil
}
