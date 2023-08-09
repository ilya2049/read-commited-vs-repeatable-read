package main

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID   int
	Name string
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)

		return
	}

	defer conn.Close(context.Background())

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open tx: %v\n", err)

		return
	}

	var users = []User{}

	err = pgxscan.Select(
		context.Background(), tx, &users, `
		select id, name from users
		where active = true
		order by id;
	`)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to scan users: %v\n", err)

		return
	}

	for _, user := range users {
		fmt.Println(user.ID, user.Name)
	}

	var total int

	err = pgxscan.Get(
		context.Background(), tx, &total, `
		select count(*) from users where active is true;
	`)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to scan user total quantity: %v\n", err)

		return
	}

	fmt.Println("total: ", total)

	if err := tx.Commit(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "unable to commit tx: %v\n", err)
	}
}
