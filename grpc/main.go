package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/antonikonovalov/benches/grpc/proto"
	"github.com/antonikonovalov/benches/grpc/server"
)

var (
	addr = flag.String(`addr`, `0.0.0.0:4569`, `binding address`)
)

func main() {
	flag.Parse()

	// init database
	db, err := sql.Open("postgres", `postgres://localhost/benches?sslmode=disable`)
	if err != nil {
		grpclog.Fatalf("failed Open postgres: %v", err)
	}

	db.SetMaxOpenConns(16 * 3 / 2)
	db.SetMaxIdleConns(16)

	cleanup(db)
	err = makeTestTables(db)
	if err != nil {
		grpclog.Fatalf("failed makeTestTables: %v", err)
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCreatorServer(grpcServer, server.New(db))
	grpcServer.Serve(lis)
}

func makeTestTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS  test_messages (
			id serial PRIMARY KEY,
			msg text NOT NULL
		)`,
	)
	return err
}

func cleanup(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS test_messages")
}
