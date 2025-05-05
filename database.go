package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SQLite3 database structure
type SQLite3 struct {
	db     *sql.DB
	dbPath string
	logger Logger
	ctx    context.Context
	cancel context.CancelFunc
}

// Logger logger interface
type Logger interface {
	Log(query string, args ...interface{})
}

// Config database configuration
type Config struct {
	DBPath       string
	MaxOpenConns int           // maximum open connections
	MaxIdleConns int           // maximum idle connections
	BusyTimeout  time.Duration // busy timeout
	EnableWAL    bool          // enable WAL mode
	ForeignKeys  bool          // enable foreign keys
	TraceLogger  Logger        // SQLite3 trace logger
}

func NewSQLite3DB(dbPath string) *SQLite3 {
	// create and return SQLite3 database
	return &SQLite3{
		dbPath: dbPath,
	}
}

func (s *SQLite3) Connect() error {
	// connect to SQLite3 database
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return fmt.Errorf("error opening SQLite3 database: %s", err)
	}
	// configure SQLite3 database pool
	db.SetMaxIdleConns(1)
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
