package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)

		return
	}

	defer conn.Close(context.Background())

	for i := 0; i < 10000; i++ {
		_, err = conn.Exec(context.Background(), `update users set active = not active where id = 3`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to execute command: %v\n", err)
		}

		fmt.Println("HA-HA!")
	}
}
