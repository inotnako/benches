package main

import (
	"net/http"

	"database/sql"
	_ "github.com/lib/pq"

	"flag"
	"github.com/kavu/go_reuseport"
	"log"
)

var (
	addr = flag.String(`addr`, `0.0.0.0:4567`, `binding address`)
)

func main() {
	flag.Parse()

	// init database
	db, err := sql.Open("postgres", `postgres://localhost/benches?sslmode=disable`)
	if err != nil {
		log.Panic(err)
	}

	db.SetMaxOpenConns(16 * 3 / 2)
	db.SetMaxIdleConns(16)

	cleanup(db)
	err = makeTestTables(db)
	if err != nil {
		log.Panic(err)
	}

	lis, err := reuseport.NewReusablePortListener("tcp4", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	http.HandleFunc(`/create`, func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			rw.WriteHeader(http.StatusNotFound)
		}

		_, err := db.Exec(`INSERT INTO test_messages (msg) VALUES ($1)`, `hello dolly!`)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Panic(err)
	}
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
