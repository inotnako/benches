package main

import (
	"database/sql"
	_ "github.com/lib/pq"

	"flag"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/antonikonovalov/benches/grpc/proto"
	"github.com/antonikonovalov/benches/grpc/server"
	"github.com/kavu/go_reuseport"
	"net/http"
	_ "net/http/pprof"

	_ "golang.org/x/net/trace"
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

	lis, err := reuseport.NewReusablePortListener("tcp4", *addr)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	go func() {
		grpclog.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

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
