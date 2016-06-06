package main

import (
	"net/http"

	"database/sql"
	_ "github.com/lib/pq"

	"flag"
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

	http.HandleFunc(`/create`, func(rw http.ResponseWriter, r *http.Request) {
		_, err := db.Exec(`INSERT INTO test_messages (msg) VALUES ($1)`, `hello dolly!`)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})

	if err := http.ListenAndServeTLS(":4568", "server.pem", "server.key", nil); err != nil {
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
