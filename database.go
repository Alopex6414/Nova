package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// DefaultConfig Initialization
func DefaultConfig(path string) Config {
	return Config{
		Path:             path,
		MaxConnections:   10,
		BusyTimeout:      5 * time.Second,
		WALMode:          true,
		ForeignKeys:      true,
		AutoVacuum:       true,
		CacheSize:        -2000, // 2MB
		JournalMode:      "wal",
		SyncMode:         "normal",
		MaxRetryAttempts: 3,
	}
}

func NewSQLite3DB(cfg Config) (*SQLite3, error) {
	// construct connection string
	dsn := fmt.Sprintf(
		"%s?_busy_timeout=%d&_foreign_keys=%d&_auto_vacuum=%d&_cache_size=%d&_journal_mode=%s&_sync=%s",
		cfg.Path,
		int(cfg.BusyTimeout.Seconds()*1000),
		boolToInt(cfg.ForeignKeys),
		boolToInt(cfg.AutoVacuum),
		cfg.CacheSize,
		cfg.JournalMode,
		cfg.SyncMode,
	)
	if cfg.WALMode {
		dsn += "&_journal_mode=WAL"
	}
	// connect to SQLite3 database
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to SQLite3 database failed: %w", err)
	}
	// configure database pool
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections / 2)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)
	// cancel function
	ctx, cancel := context.WithCancel(context.Background())
	return &SQLite3{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
		config: cfg,
		stmtCache: &StmtCache{
			cache: make(map[string]*sql.Stmt),
		},
	}, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
