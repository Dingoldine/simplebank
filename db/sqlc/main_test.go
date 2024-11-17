package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simplebank?sslmode=disable"
)

var testQueries *Queries
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error

	testPool, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Cannot connect using pgxpool", err)
	}
	defer testPool.Close()

	testQueries = New(testPool)

	os.Exit(m.Run())
}
