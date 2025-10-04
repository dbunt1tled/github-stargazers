package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3" //nolint:blank-imports // base db class
)

type Storage struct {
	db *sql.DB
}

func New(ctx context.Context, dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS stars (
			id INTEGER PRIMARY KEY,
			repo TEXT NOT NULL,
			user TEXT NOT NULL,
			date TEXT NOT NULL,
            UNIQUE(date,user,repo)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}
func (s *Storage) Add(ctx context.Context, repo, date string, users []string) error {
	var (
		tx   *sql.Tx
		err  error
		stmt *sql.Stmt
	)
	tx, err = s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err = tx.PrepareContext(ctx, "INSERT OR IGNORE INTO stars (repo, date, user) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	for _, u := range users {
		if _, err = stmt.ExecContext(ctx, repo, date, u); err != nil {
			return err
		}
	}
	err = stmt.Close() //nolint:sqlclosecheck // stupid linter can't understand check err with defer
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Storage) GetPreviousDate(ctx context.Context, repo string, today string) (string, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT MAX(date) FROM stars WHERE repo = ? AND date < ? ORDER BY date DESC LIMIT 1`,
		repo,
		today,
	)
	var date *string
	if err := row.Scan(&date); err != nil {
		return "", err
	}
	if date == nil {
		return "", nil
	}
	return *date, nil
}

func (s *Storage) Diff(ctx context.Context, repo string, today string, prev string) ([]string, []string, error) {
	var (
		added, removed []string
		rows           *sql.Rows
	)
	stmt, err := s.db.PrepareContext( //nolint:sqlclosecheck // stupid linter
		ctx,
		`SELECT user 
		 FROM stars 
		 WHERE repo = ? AND date = ? AND user NOT IN (SELECT user FROM stars WHERE repo = ? AND date = ?)`,
	)
	if err != nil {
		return nil, nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	rows, err = stmt.QueryContext(ctx, repo, today, repo, prev)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, nil, rows.Err()
		}
		var u string
		if err = rows.Scan(&u); err != nil {
			return nil, nil, err
		}
		added = append(added, u)
	}
	err = rows.Close() //nolint:sqlclosecheck // stupid linter can't understand defer Close for next rows
	if err != nil {
		return nil, nil, err
	}

	rows, err = stmt.QueryContext(ctx, repo, prev, repo, today) //nolint:sqlclosecheck // stupid linter
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		if rows.Err() != nil {
			return nil, nil, rows.Err()
		}
		var u string
		if err = rows.Scan(&u); err != nil {
			return nil, nil, err
		}
		removed = append(removed, u)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	return added, removed, nil
}

func (s *Storage) GetStargazers(ctx context.Context) ([]string, error) {
	row := s.db.QueryRowContext(ctx, `SELECT MAX(date) FROM stars LIMIT 1`)
	var (
		date       *string
		u          string
		stargazers []string
	)
	if err := row.Scan(&date); err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	rows, err := s.db.QueryContext( //nolint:sqlclosecheck // stupid linter
		ctx,
		`SELECT DISTINCT user FROM stars WHERE date = ?`,
		*date,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		if rows.Err() != nil {
			continue
		}
		if err = rows.Scan(&u); err != nil {
			continue
		}
		stargazers = append(stargazers, u)
	}
	return stargazers, nil
}
