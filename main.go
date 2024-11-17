package main

import (
	"context"
	"fmt"
	"log"
	db "simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *db.Queries
var Pool *pgxpool.Pool

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simplebank?sslmode=disable"
)

func main() {
	fmt.Println("simplebank")
	var err error
	Pool, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Cannot connect using pgxpool", err)
	}
	err = Pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Connection to db failed: ", err)
	}
	defer Pool.Close()
}
