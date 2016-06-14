package server

import (
	"database/sql"
	"testing"

	pb "github.com/antonikonovalov/benches/grpc/proto"
	_ "github.com/lib/pq"
)

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

func BenchmarkPsqlInsert_Parallel(b *testing.B) {
	db, err := sql.Open("postgres", `postgres://localhost/benches?sslmode=disable`)
	if err != nil {
		b.Fatalf("failed Open postgres: %v", err)
	}

	db.SetMaxOpenConns(16 * 3 / 2)
	db.SetMaxIdleConns(16)

	cleanup(db)
	err = makeTestTables(db)
	if err != nil {
		b.Fatalf("failed makeTestTables: %v", err)
	}

	//b.SetParallelism(4)
	b.RunParallel(func(pbt *testing.PB) {
		for pbt.Next() {
			resp := new(pb.MessageResponse)
			err := db.QueryRow(`INSERT INTO test_messages (msg) VALUES ($1) RETURNING id`, `aloxasxasd`).Scan(&resp.Id)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkPsqlInsert(b *testing.B) {
	db, err := sql.Open("postgres", `postgres://localhost/benches?sslmode=disable`)
	if err != nil {
		b.Fatalf("failed Open postgres: %v", err)
	}

	db.SetMaxOpenConns(16 * 3 / 2)
	db.SetMaxIdleConns(16)

	cleanup(db)
	err = makeTestTables(db)
	if err != nil {
		b.Fatalf("failed makeTestTables: %v", err)
	}

	for i := 0; i < b.N; i++ {
		resp := new(pb.MessageResponse)
		err := db.QueryRow(`INSERT INTO test_messages (msg) VALUES ($1) RETURNING id`, `aloxasxasd`).Scan(&resp.Id)
		if err != nil {
			b.Error(err)
		}
	}
}
