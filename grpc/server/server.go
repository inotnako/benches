package server

import (
	"database/sql"
	pb "github.com/antonikonovalov/benches/grpc/proto"
	"golang.org/x/net/context"
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
	return resp, err
}
