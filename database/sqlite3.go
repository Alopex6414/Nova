package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"modernc.org/sqlite"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Config struct {
	MaxOpenConns    int
	Debug           bool
	AutoCreateTable bool
	WALMode         bool
	CacheSize       int
	BusyTimeout     int
	Extensions      []string
}

// SqliteDB Sqlite3 Database
type SqliteDB struct {
	db        *sql.DB
	dsn       string
	config    *Config
	stmtCache sync.Map
	rwLock    sync.RWMutex
	metrics   *Metrics
	migration *MigrationMgr
}

// Metrics Performance
type Metrics struct {
	QueryCount       int64
	WriteCount       int64
	AvgQueryTime     time.Duration
	MaxQueryTime     time.Duration
	LastError        error
	LastErrorTime    time.Time
	ConnectionsInUse int
}

// MigrationMgr Migration Management
type MigrationMgr struct {
	mu         sync.Mutex
	versions   map[int]func(tx *sql.Tx) error
	currentVer int
}

func NewSqliteDB(dsn string, cfg *Config) (*SqliteDB, error) {
	// create sqlite3 database
	db, err := sql.Open("sqlite", buildDSN(dsn, cfg))
	if err != nil {
		return nil, fmt.Errorf("open database error: %w", err)
	}
	// configure settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(0)
	// create sqlite3 instance
	sdb := &SqliteDB{
		db:        db,
		dsn:       dsn,
		config:    cfg,
		metrics:   &Metrics{},
		migration: &MigrationMgr{versions: make(map[int]func(*sql.Tx) error)},
	}
	// initialize sqlite3 database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sdb.initDatabase(ctx); err != nil {
		return nil, err
	}
	// enable monitor process
	if cfg.Debug {
		go sdb.monitor()
	}
	return sdb, nil
}

// construct DSN with parameters
func buildDSN(base string, cfg *Config) string {
	var params []string
	if cfg.WALMode {
		params = append(params, "_journal_mode=WAL")
	}
	if cfg.CacheSize > 0 {
		params = append(params, fmt.Sprintf("_cache_size=%d", cfg.CacheSize))
	}
	if cfg.BusyTimeout > 0 {
		params = append(params, fmt.Sprintf("_busy_timeout=%d", cfg.BusyTimeout))
	}
	return base + "?" + strings.Join(params, "&")
}

// initialize sqlite3 database
func (s *SqliteDB) initDatabase(ctx context.Context) error {
	// execute context
	if _, err := s.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("enable foreign keys failed: %w", err)
	}
	// auto create migrate table
	if s.config.AutoCreateTable {
		if err := s.createMigrationTable(ctx); err != nil {
			return err
		}
	}
	return nil
}

// QueryWithRetry query database with retry
func (s *SqliteDB) QueryWithRetry(ctx context.Context, maxRetries int, query string, args ...interface{}) (*sql.Rows, error) {
	for i := 0; ; i++ {
		rows, err := s.Query(ctx, query, args...)
		if isRetryableError(err) && i < maxRetries {
			time.Sleep(time.Duration(i*100) * time.Millisecond)
			continue
		}
		return rows, err
	}
}

// InsertStruct insert structure
func (s *SqliteDB) InsertStruct(ctx context.Context, table string, data interface{}) (int64, error) {
	// structure reflect
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	fields := make([]string, 0, rv.NumField())
	values := make([]interface{}, 0, rv.NumField())
	// insert dataset
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		if tag := field.Tag.Get("db"); tag != "" {
			fields = append(fields, tag)
			values = append(values, rv.Field(i).Interface())
		}
	}
	// query sentence
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fields, ","),
		strings.Repeat("?,", len(fields)-1)+"?",
	)
	// execute query
	res, err := s.Exec(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// monitor process
func (s *SqliteDB) monitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := s.db.Stats()
		s.metrics.ConnectionsInUse = stats.InUse
		// record other metrics...
	}
}

// AddMigration migration methods
func (s *SqliteDB) AddMigration(version int, fn func(tx *sql.Tx) error) {
	s.migration.mu.Lock()
	defer s.migration.mu.Unlock()
	s.migration.versions[version] = fn
}

func (s *SqliteDB) RunMigrations(ctx context.Context) error {
	return s.WithTransaction(ctx, func(tx *sql.Tx) error {
		// migration logic...
		return nil
	})
}

// retryable error handle
func isRetryableError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == 5 // SQLITE_BUSY
	}
	return false
}

// CheckConnection health check for connection
func (s *SqliteDB) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

// QueryWithMetrics query with performance metrics
func (s *SqliteDB) QueryWithMetrics(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	defer func() {
		dur := time.Since(start)
		s.metrics.AvgQueryTime = (s.metrics.AvgQueryTime*time.Duration(s.metrics.QueryCount) + dur) /
			time.Duration(s.metrics.QueryCount+1)
		if dur > s.metrics.MaxQueryTime {
			s.metrics.MaxQueryTime = dur
		}
		s.metrics.QueryCount++
	}()
	return s.Query(ctx, query, args...)
}
