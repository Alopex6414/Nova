package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
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

// 核心增强功能实现

// ExecWithRetry 带重试机制的 Exec
func (s *SQLite3) ExecWithRetry(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error
	for attempt := 0; attempt < s.config.MaxRetryAttempts; attempt++ {
		result, err = s.db.ExecContext(s.ctx, query, args...)
		if err == nil {
			return result, nil
		}
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.Code, sqlite3.ErrBusy) {
				time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
				continue
			}
		}
		break
	}
	return nil, fmt.Errorf("执行失败（重试 %d 次）: %w", s.config.MaxRetryAttempts, err)
}

// QueryScan 查询并自动扫描到结构体切片
func (s *SQLite3) QueryScan(dest interface{}, query string, args ...interface{}) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return ErrInvalidPointer
	}
	sliceVal := destVal.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("目标必须是指向切片的指针")
	}
	rows, err := s.db.QueryContext(s.ctx, query, args...)
	if err != nil {
		return fmt.Errorf("查询失败: %w", err)
	}
	defer rows.Close()

	elementType := sliceVal.Type().Elem()
	for rows.Next() {
		elem := reflect.New(elementType).Elem()
		fields, err := s.getFieldAddresses(elem)
		if err != nil {
			return err
		}

		if err := rows.Scan(fields...); err != nil {
			return fmt.Errorf("扫描失败: %w", err)
		}

		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	return rows.Err()
}

// getFieldAddresses 通过结构体标签获取字段地址
func (s *SQLite3) getFieldAddresses(v reflect.Value) ([]interface{}, error) {
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("需要结构体类型")
	}

	t := v.Type()
	columns := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		columns[tag] = v.Field(i).Addr().Interface()
	}

	rows, err := s.db.QueryContext(s.ctx, "SELECT * FROM table LIMIT 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(colNames))
	for i, name := range colNames {
		addr, ok := columns[name]
		if !ok {
			return nil, fmt.Errorf("缺少对应字段: %s", name)
		}
		result[i] = addr
	}

	return result, nil
}

// Backup 在线热备份
func (s *SQLite3) Backup(dest string) error {
	conn, err := s.db.Conn(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Raw(func(destConn interface{}) error {
		destDb, ok := destConn.(*sqlite3.SQLiteConn)
		if !ok {
			return fmt.Errorf("无效的目标连接")
		}

		srcDb, err := sqlite3.Open(s.config.Path)
		if err != nil {
			return err
		}
		defer srcDb.Close()

		backup, err := destDb.Backup("main", srcDb, "main")
		if err != nil {
			return err
		}

		_, err = backup.Step(-1)
		if err != nil {
			return err
		}

		return backup.Finish()
	})
}

// RegisterHook 注册数据库操作钩子
func (s *SQLite3) RegisterHook(hook Hook) {
	sqlite3.RegisterUpdateHook(s.db, func(op int, db string, table string, rowid int64) {
		var operation string
		switch op {
		case sqlite3.SQLITE_INSERT:
			operation = "INSERT"
		case sqlite3.SQLITE_UPDATE:
			operation = "UPDATE"
		case sqlite3.SQLITE_DELETE:
			operation = "DELETE"
		}
		hook(operation, fmt.Sprintf("%s.%s", db, table), nil)
	})
}

// AnalyzePerformance 性能分析工具
func (s *SQLite3) AnalyzePerformance() map[string]interface{} {
	stats := s.db.Stats()
	result := map[string]interface{}{
		"connection_stats": stats,
		"page_size":        s.querySingleInt("PRAGMA page_size"),
		"cache_size":       s.querySingleInt("PRAGMA cache_size"),
		"schema_size":      s.querySingleInt("SELECT SUM(pgsize) FROM dbstat"),
	}

	if s.config.EnableMetrics {
		result["metrics"] = s.metrics
	}

	return result
}

// 辅助方法
func (s *SQLite3) querySingleInt(query string) int {
	var result int
	row := s.db.QueryRowContext(s.ctx, query)
	_ = row.Scan(&result)
	return result
}

// Close 安全关闭
func (s *SQLite3) Close() error {
	s.cancel()

	// 清理预处理缓存
	for _, stmt := range s.stmtCache.cache {
		_ = stmt.Close()
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("关闭数据库失败: %w", err)
	}

	if s.config.AutoVacuum {
		if err := os.Remove(s.config.Path + "-wal"); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("清理WAL文件失败: %w", err)
		}
	}

	return nil
}

// JSON 类型支持
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("类型断言失败")
	}
	return json.Unmarshal(b, &j)
}

// 自动迁移工具
func (s *SQLite3) AutoMigrate(models ...interface{}) error {
	for _, model := range models {
		stmt, err := generateCreateTable(model)
		if err != nil {
			return err
		}

		if _, err := s.Exec(stmt); err != nil {
			return fmt.Errorf("迁移失败: %w", err)
		}
	}
	return nil
}

func generateCreateTable(model interface{}) (string, error) {
	// 实现基于反射的DDL生成
	// ...
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
