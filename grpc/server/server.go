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

func (s *server) CreateByStream2(stream pb.Creator_CreateByStreamServer) error {
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

func (s *server) CreateByStream(stream pb.Creator_CreateByStreamServer) error {
	in, errChan, out := s.processor()
	go dispatcher(stream, in, errChan)
	for {
		select {
		case e := <-errChan:
			if e == io.EOF {
				return nil
			}
			return e
		case resp := <-out:
			if resp != nil {
				err := stream.Send(resp)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *server) processor() (chan *pb.MessageRequest, chan error, chan *pb.MessageResponse) {
	w := 50
	in := make(chan *pb.MessageRequest)
	errC := make(chan error)
	out := make(chan *pb.MessageResponse)

	// start workers
	for i := 0; i < w; i++ {
		go func() {
			for msg := range in {
				resp := new(pb.MessageResponse)
				err := s.db.QueryRow(`INSERT INTO test_messages (msg) VALUES ($1) RETURNING id`, msg.Msg).Scan(&resp.Id)
				if err != nil {
					errC <- err
					close(errC)
					close(out)
					close(in)
					return
				} else {
					out <- resp
				}

			}
		}()
	}

	return in, errC, out
}

func dispatcher(stream pb.Creator_CreateByStreamServer, in chan *pb.MessageRequest, errColl chan error) {

	for {
		msg, err := stream.Recv()
		if err != nil {
			errColl <- err
			return
		} else {
			in <- msg
		}
	}

}
