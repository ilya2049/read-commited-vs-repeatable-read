package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

var (
	mode  = flag.String("mode", "read-committed", "execution mode of the transaction")
	reads = flag.Int("reads", 10, "count of reads in the transaction")
)

type Password struct {
	ID   int
	Hash string
}

func main() {
	flag.Parse()

	fmt.Fprintf(os.Stdout, "mode: %s\n", *mode)
	fmt.Fprintf(os.Stdout, "reads: %d\n", *reads)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	ctx, cancel := context.WithCancel(context.Background())

	const workers = 100

	go runUpdater(ctx, random)

	var wg sync.WaitGroup

	wg.Add(workers)

	start := time.Now()

	for i := 0; i < workers; i++ {
		go func() {
			runTx(*mode, *reads)

			wg.Done()
		}()
	}

	wg.Wait()
	cancel()

	fmt.Fprintf(os.Stdout, "elapsed: %dms\n", time.Since(start).Milliseconds())
}

func runTx(mode string, reads int) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)

		return
	}

	defer conn.Close(context.Background())

	tx, err := beginTx(mode, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open tx: %v\n", err)

		return
	}

	var passwords = []Password{}

	for i := 0; i < reads; i++ {
		err = pgxscan.Select(
			context.Background(), tx, &passwords, `select * from passwords;`)

		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to scan passwords: %v\n", err)

			return
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "unable to commit tx: %v\n", err)
	}
}

func beginTx(mode string, conn *pgx.Conn) (pgx.Tx, error) {
	switch mode {
	case "read-committed":
		return conn.BeginTx(context.Background(), pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		})
	case "repeatable-read":
		return conn.BeginTx(context.Background(), pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
	case "repeatable-read-read-only":
		return conn.BeginTx(context.Background(), pgx.TxOptions{
			IsoLevel:   pgx.RepeatableRead,
			AccessMode: pgx.ReadOnly,
		})
	}

	panic("unsupported tx mode")
}

func runUpdater(ctx context.Context, random *rand.Rand) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)

		return
	}

	defer conn.Close(context.Background())

	for {
		select {
		case <-ctx.Done():
			return
		default:
			tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to open tx: %v\n", err)

				return
			}

			id := random.Intn(10000)

			row := conn.QueryRow(ctx, `select hash from passwords where id = $1 for update`, id)
			var hash string

			if err := row.Scan(&hash); err != nil {
				fmt.Fprintf(os.Stderr, "unable to scan a password hash: %v\n", err)

				return
			}

			time.Sleep(50 * time.Millisecond)

			hash = RandStringRunes(25, random)

			_, err = conn.Exec(ctx, `
			update passwords set hash = $1
			where id = $2;
		`, hash, id)

			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to update a password: %v\n", err)

				return
			}

			if err := tx.Commit(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "unable to commit tx: %v\n", err)
			}
		}
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int, random *rand.Rand) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(b)
}
