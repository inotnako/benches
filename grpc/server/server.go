package server

import (
	"database/sql"
	pb "github.com/antonikonovalov/benches/grpc/proto"
	"golang.org/x/net/context"
	"io"
	"time"
)

func New(db *sql.DB) pb.CreatorServer {
	return &server{db}
}

type server struct {
	db *sql.DB
}

func (s *server) Create(_ context.Context, msg *pb.MessageRequest) (*pb.MessageResponse, error) {
	resp := new(pb.MessageResponse)
	err := s.db.QueryRow(`INSERT INTO test_messages (msg) VALUES ($1) RETURNING id`, msg.Msg).Scan(&resp.Id)
	resp.Time = time.Now().String()
	return resp, err
}

type Request struct {
	m *pb.MessageRequest
	e error
}

func (s *server) CreateByStream(stream pb.Creator_CreateByStreamServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		resp := new(pb.MessageResponse)
		err = s.db.QueryRow(`INSERT INTO test_messages (msg) VALUES ($1) RETURNING id`, msg.Msg).Scan(&resp.Id)
		if err != nil {
			return err
		}
		resp.Time = time.Now().String()

		if err = stream.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) CreateByStream2(stream pb.Creator_CreateByStreamServer) error {
	requests := RequestMessages(stream)
	in, errChan, out := s.Processor()
	for {
		select {
		case req := <-requests:
			print(`<- `)
			if req.e == io.EOF {
				return nil
			}
			if req.e != nil {
				return req.e
			}
			in <- req.m
			print(` | `)
		case e := <-errChan:
			return e
		case resp := <-out:
			println(`- > `)
			err := stream.Send(resp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *server) Processor() (chan *pb.MessageRequest, chan error, chan *pb.MessageResponse) {
	in := make(chan *pb.MessageRequest, 0)
	errC := make(chan error, 0)
	out := make(chan *pb.MessageResponse, 0)

	go func() {
		var i int64
		for _ = range in {
			//time.Sleep(10 * time.Microsecond)
			i++
			out <- &pb.MessageResponse{Id: (1000000000000000004 + i)}
			//select {
			//case <-in:
			//println(`PR IN:`)
			//resp := new(pb.MessageResponse)

			//resp.Id = time.Now().UnixNano() + 123234345
			//err := s.db.QueryRow(`INSERT INTO test_messages (msg) VALUES ($1) RETURNING id`, msg.Msg).Scan(&resp.Id)
			//if err != nil {
			//	//println(`PR ERR:`)
			//	errC <- err
			//	close(errC)
			//	close(out)
			//	close(in)
			//	return
			//} else {
			//println(`PR OUT:`)

			//}
			//}
		}
		close(errC)
		close(out)
		close(in)
		return
	}()

	return in, errC, out
}

func RequestMessages(stream pb.Creator_CreateByStreamServer) chan *Request {
	messages := make(chan *Request, 0)

	go func() {
		for {
			msg, err := stream.Recv()
			messages <- &Request{msg, err}
			if err != nil {
				close(messages)
				return
			}
		}
	}()

	return messages
}
