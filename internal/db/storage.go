package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
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
func (s *Storage) Add(repo, date string, users []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO stars (repo, date, user) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	for _, u := range users {
		if _, err := stmt.Exec(repo, date, u); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) GetPreviousDate(repo string, today string) (string, error) {
	row := s.db.QueryRow("SELECT MAX(date) FROM stars WHERE repo = ? AND date < ? ORDER BY date DESC LIMIT 1", repo, today)
	var date *string
	if err := row.Scan(&date); err != nil {
		return "", err
	}
	if date == nil {
		return "", nil
	}
	return *date, nil
}

func (s *Storage) Diff(repo string, today string, prev string) ([]string, []string, error) {
	var (
		added, removed []string
		rows           *sql.Rows
	)
	stmt, err := s.db.Prepare(`
		SELECT user FROM stars 
		WHERE repo = ? AND date = ? 
		AND user NOT IN (
			SELECT user FROM stars WHERE repo = ? AND date = ?
		)`)
	if err != nil {
		return nil, nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	rows, err = stmt.Query(repo, today, repo, prev)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, nil, err
		}
		added = append(added, u)
	}
	err = rows.Close()
	if err != nil {
		return nil, nil, err
	}

	rows, err = stmt.Query(repo, prev, repo, today)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, nil, err
		}
		removed = append(removed, u)
	}
	err = rows.Close()
	if err != nil {
		return nil, nil, err
	}
	return added, removed, nil
}
