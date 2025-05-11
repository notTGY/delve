package memory

import (
	"context"
	"database/sql"
	"time"
  "errors"

	_ "modernc.org/sqlite"
)

func init() {
  migrate()
}

type Dialog struct {
  timestamp string
  UserMessage string
  AssistantMessage string
}

func Query(tgid int64) ([]Dialog, error) {
  dialog := []Dialog{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open(
		"sqlite",
		"./db/sqlite.db?_pragma=busy_timeout(5000)",
	)
	if err != nil {
    return dialog, err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
    return dialog, err
	}
	defer tx.Rollback()

  _, err = tx.ExecContext(
    ctx,
    `INSERT OR IGNORE INTO users (tgid) VALUES (?)`,
    tgid,
  )
	if err != nil {
    return dialog, err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id FROM users WHERE tgid = ?`,
    tgid,
	)
	defer rows.Close()
	if err != nil {
    return dialog, err
	}

	if err := rows.Err(); err != nil {
    return dialog, err
	}

  user_id := -1
  for rows.Next() {
		if err := rows.Scan(&user_id); err != nil {
      return dialog, err
		}
    break
  }
  if user_id == -1 {
    return dialog, nil
  }

	rows, err = tx.QueryContext(
		ctx,
		`SELECT
      user_message,
      assistant_message,
      timestamp
    FROM history WHERE user_id = ?`,
    user_id,
	)
	defer rows.Close()
	if err != nil {
    return dialog, err
	}

	if err := rows.Err(); err != nil {
    return dialog, err
	}

  for rows.Next() {
    var timestamp, UserMessage, AssistantMessage string
		if err := rows.Scan(
      &UserMessage,
      &AssistantMessage,
      &timestamp,
    ); err != nil {
      return dialog, err
		}
    dialog = append(dialog, Dialog{
      timestamp,
      UserMessage,
      AssistantMessage,
    })
  }
  tx.Commit()
  return dialog, nil
}

func Save(tgid int64, UserMessage, AssistantMessage string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open(
		"sqlite",
		"./db/sqlite.db?_pragma=busy_timeout(5000)",
	)
	if err != nil {
    return err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
    return err
	}
	defer tx.Rollback()

  _, err = tx.ExecContext(
    ctx,
    `INSERT OR IGNORE INTO users (tgid) VALUES (?)`,
    tgid,
  )
	if err != nil {
    return err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id FROM users WHERE tgid = ?`,
    tgid,
	)
	defer rows.Close()
	if err != nil {
    return err
	}

	if err := rows.Err(); err != nil {
    return err
	}

  user_id := -1
  for rows.Next() {
		if err := rows.Scan(&user_id); err != nil {
      return err
		}
    break
  }
  if user_id == -1 {
    return errors.New("Failed to create user")
  }

  _, err = tx.ExecContext(
    ctx,
    `INSERT INTO history (
      user_id,
      user_message,
      assistant_message
    ) VALUES (?, ?, ?)`,
    user_id,
    UserMessage,
    AssistantMessage,
  )
	if err != nil {
    return err
	}
  tx.Commit()
  return nil
}
