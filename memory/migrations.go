package memory

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func migrate() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open(
		"sqlite",
		"./db/sqlite.db?_pragma=busy_timeout(5000)",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS users (
      id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
			tgid text NOT NULL
		)`,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS history (
      id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
      timestamp integer DEFAULT (current_timestamp),
      user_id integer NOT NULL,
      user_message text NOT NULL,
      assistant_message text NOT NULL
		)`,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE UNIQUE INDEX IF NOT EXISTS users_tgid_unique ON users (tgid)`,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	tx.Commit()
}
