package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dsn      = "postgresql://root:rootcc@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
