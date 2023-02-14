package connec

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	var err error

	databaseUrl := "postgres://postgres:achyarbagus17@localhost:5432/Personal-Web"
	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unnable to connect to Database: %v", err)
		os.Exit(1)
	}
	fmt.Println("succes connect to database")
}
