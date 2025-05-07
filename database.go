package main

import (
	""
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// define SQLite3 error type
var (
	ErrDBNotOpen      = errors.New("database not open")
	ErrInvalidPointer = errors.New("invalid pointer type")
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrLockTimeout    = errors.New("database lock timeout")
)

// SQLite3 database type structure
type SQLite3 struct {
	db        *sql.DB
	ctx       context.Context
	cancel    context.CancelFunc
	config    Config
	metrics   Metrics
	stmtCache *StmtCache
}

// Config database configure
type Config struct {
	Path             string
	MaxConnections   int           `json:"max_connections"`
	BusyTimeout      time.Duration `json:"busy_timeout"`
	WALMode          bool          `json:"wal_mode"`
	ForeignKeys      bool          `json:"foreign_keys"`
	AutoVacuum       bool          `json:"auto_vacuum"`
	CacheSize        int           `json:"cache_size"`
	JournalMode      string        `json:"journal_mode"`
	SyncMode         string        `json:"sync_mode"`
	EnableTrace      bool          `json:"enable_trace"`
	EnableMetrics    bool          `json:"enable_metrics"`
	MaxRetryAttempts int           `json:"max_retry_attempts"`
}

// Metrics database metrics
type Metrics struct {
	QueryCount     int64
	WriteCount     int64
	ErrorCount     int64
	RetryCount     int64
	AvgQueryTime   time.Duration
	MaxQueryTime   time.Duration
	LastBackupTime time.Time
}

// StmtCache pre-handling cache
type StmtCache struct {
	cache map[string]*sql.Stmt
}

// Hook function type
type Hook func(operation string, query string, args []interface{})

func NewSQLite3DB(cfg Config) *SQLite3 {
	return &SQLite3{}
}
