package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
	"os"
	"reflect"
	"strings"
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
	// connection options string
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

// ExecWithRetry Exec with Retry
func (s *SQLite3) ExecWithRetry(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error
	// execute query with retry
	for attempt := 0; attempt < s.config.MaxRetryAttempts; attempt++ {
		// execute get result
		result, err = s.db.ExecContext(s.ctx, query, args...)
		if err == nil {
			return result, nil
		}
		// handle execute error
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.Code, sqlite3.ErrBusy) {
				time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
				continue
			}
		}
		break
	}
	return nil, fmt.Errorf("execute failed (retry %d): %w", s.config.MaxRetryAttempts, err)
}

// QueryScan query and auto scan structure slice
func (s *SQLite3) QueryScan(dest interface{}, query string, args ...interface{}) error {
	destVal := reflect.ValueOf(dest)
	// type should not be pointer or nil value
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return ErrInvalidPointer
	}
	sliceVal := destVal.Elem()
	// type should not be a slice
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("target must be a slice pointer")
	}
	// query context
	rows, err := s.db.QueryContext(s.ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	// get element type
	elementType := sliceVal.Type().Elem()
	for rows.Next() {
		// get element
		elem := reflect.New(elementType).Elem()
		fields, err := s.getFieldAddresses(elem)
		if err != nil {
			return err
		}
		if err := rows.Scan(fields...); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		sliceVal.Set(reflect.Append(sliceVal, elem))
	}
	return rows.Err()
}

// getFieldAddresses get field through structure label
func (s *SQLite3) getFieldAddresses(v reflect.Value) ([]interface{}, error) {
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type must be a struct")
	}
	// get structure type
	t := v.Type()
	columns := make(map[string]interface{}, v.NumField())
	// iterate fields in structure
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		columns[tag] = v.Field(i).Addr().Interface()
	}
	// query context
	rows, err := s.db.QueryContext(s.ctx, "SELECT * FROM table LIMIT 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// get column names
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// return results
	result := make([]interface{}, len(colNames))
	for i, name := range colNames {
		addr, ok := columns[name]
		if !ok {
			return nil, fmt.Errorf("failed cause lack of: %s", name)
		}
		result[i] = addr
	}
	return result, nil
}

// Close database safety
func (s *SQLite3) Close() error {
	s.cancel()
	// clean cache
	for _, stmt := range s.stmtCache.cache {
		_ = stmt.Close()
	}
	// close database
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	// clean WAL file
	if s.config.AutoVacuum {
		if err := os.Remove(s.config.Path + "-wal"); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to clean WAL file: %w", err)
		}
	}
	return nil
}

// Transaction event handle
func (s *SQLite3) Transaction(fn func(tx *sql.Tx) error) error {
	// database start transaction
	tx, err := s.db.BeginTx(s.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// raise exception and rollback
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	// callback function for rollback
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed rollback: %v (original error: %w)", rbErr, err)
		}
		return err
	}
	// transaction commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed commit: %w", err)
	}
	return nil
}

// JSON type interface
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Test() {
	// create database configure
	cfg := DefaultConfig("sqlite3.db")
	cfg.EnableTrace = true
	cfg.EnableMetrics = true
	// create sqlite3 database
	db, err := NewSQLite3DB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// define data structure
	type User struct {
		ID        int       `db:"id"`
		Name      string    `db:"name"`
		Metadata  JSON      `db:"metadata"`
		CreatedAt time.Time `db:"created_at"`
	}
	// insert data
	user := User{
		Name:      "John Doe",
		Metadata:  JSON{"department": "IT"},
		CreatedAt: time.Now(),
	}
	// begin transaction
	err = db.Transaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO users (name, metadata, created_at) VALUES (?, ?, ?)",
			user.Name,
			user.Metadata,
			user.CreatedAt,
		)
		return err
	})
	// error handle
	if err != nil {
		log.Fatal("failed to transaction:", err)
	}
}
