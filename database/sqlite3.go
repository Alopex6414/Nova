package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"modernc.org/sqlite"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Config
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

// Performance Metrics
type Metrics struct {
	QueryCount       int64
	WriteCount       int64
	AvgQueryTime     time.Duration
	MaxQueryTime     time.Duration
	LastError        error
	LastErrorTime    time.Time
	ConnectionsInUse int
}

// Migration Management
type MigrationMgr struct {
	mu         sync.Mutex
	versions   map[int]func(tx *sql.Tx) error
	currentVer int
}

func NewSqliteDB(dsn string, cfg *Config) (*SqliteDB, error) {
	// 初始化数据库连接
	db, err := sql.Open("sqlite", buildDSN(dsn, cfg))
	if err != nil {
		return nil, fmt.Errorf("open database error: %w", err)
	}
	// 连接设置
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(0)
	// 创建增强实例
	sdb := &SqliteDB{
		db:        db,
		dsn:       dsn,
		config:    cfg,
		metrics:   &Metrics{},
		migration: &MigrationMgr{versions: make(map[int]func(*sql.Tx) error)},
	}
	// 初始化数据库
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sdb.initDatabase(ctx); err != nil {
		return nil, err
	}
	// 启动监控协程
	if cfg.Debug {
		go sdb.monitor()
	}
	return sdb, nil
}

// 构建带参数的DSN
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

// 初始化数据库
func (s *SqliteDB) initDatabase(ctx context.Context) error {
	// 执行基础设置
	if _, err := s.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("enable foreign keys failed: %w", err)
	}

	// 自动创建迁移表
	if s.config.AutoCreateTable {
		if err := s.createMigrationTable(ctx); err != nil {
			return err
		}
	}
	return nil
}

// 带重试机制的查询
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

// 结构体插入（简化ORM功能）
func (s *SqliteDB) InsertStruct(ctx context.Context, table string, data interface{}) (int64, error) {
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	fields := make([]string, 0, rv.NumField())
	values := make([]interface{}, 0, rv.NumField())

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		if tag := field.Tag.Get("db"); tag != "" {
			fields = append(fields, tag)
			values = append(values, rv.Field(i).Interface())
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(fields, ","),
		strings.Repeat("?,", len(fields)-1)+"?",
	)

	res, err := s.Exec(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// 监控协程
func (s *SqliteDB) monitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := s.db.Stats()
		s.metrics.ConnectionsInUse = stats.InUse
		// 记录其他指标...
	}
}

// 迁移管理方法
func (s *SqliteDB) AddMigration(version int, fn func(tx *sql.Tx) error) {
	s.migration.mu.Lock()
	defer s.migration.mu.Unlock()
	s.migration.versions[version] = fn
}

func (s *SqliteDB) RunMigrations(ctx context.Context) error {
	return s.WithTransaction(ctx, func(tx *sql.Tx) error {
		// 迁移逻辑...
		return nil
	})
}

// 错误处理增强
func isRetryableError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == 5 // SQLITE_BUSY
	}
	return false
}

// 连接池健康检查
func (s *SqliteDB) CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

// 带性能监控的查询
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

// 示例使用
func main() {
	cfg := &Config{
		MaxOpenConns:    1,
		Debug:           true,
		AutoCreateTable: true,
		WALMode:         true,
		CacheSize:       2000,
		BusyTimeout:     5000,
	}

	db, err := NewSqliteDB("file:demo.db", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 添加迁移
	db.AddMigration(1, func(tx *sql.Tx) error {
		_, err := tx.Exec(`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE
		)`)
		return err
	})

	// 运行迁移
	if err := db.RunMigrations(context.Background()); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// 插入结构体
	type User struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	user := User{Name: "Alice", Email: "alice@example.com"}
	id, err := db.InsertStruct(context.Background(), "users", user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted ID:", id)
}
