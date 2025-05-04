package main

import (
	"database/sql"
	"fmt"
	"time"
)

type SQLite3 struct {
	db       *sql.DB
	dbPath   string
	maxConns int
	timeout  time.Duration
}

func NewSQLite3DB(dbPath string) *SQLite3 {
	// create and return SQLite3 database
	return &SQLite3{
		dbPath:   dbPath,
		maxConns: 5,
		timeout:  5 * time.Second,
	}
}

func (s *SQLite3) Connect() error {
	// connect to SQLite3 database
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return fmt.Errorf("error opening SQLite3 database: %s", err)
	}
	// configure SQLite3 database pool
	db.SetMaxOpenConns(s.maxConns)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(s.timeout)
	// verify SQLite3 database connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging SQLite3 database: %s", err)
	}
	// return SQLite3 database
	s.db = db
	return nil
}

func (s *SQLite3) Close() error {
	// close SQLite3 database
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLite3) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *SQLite3) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *SQLite3) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *SQLite3) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *SQLite3) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := s.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
