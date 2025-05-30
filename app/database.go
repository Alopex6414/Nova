package app

import (
	"errors"
	"fmt"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
	. "nova/database"
)

type DB struct {
	sqliteDB *SQLiteDB
}

func NewDB(dbPath string) (*DB, error) {
	// create sqlite3 database
	sqliteDB, err := NewSQLiteDB(dbPath)
	if err != nil {
		return nil, err
	}
	return &DB{sqliteDB}, err
}

func (db *DB) CreateTables() error {
	// create user table
	sql := `
		CREATE TABLE IF NOT EXISTS users (
		userId TEXT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		phoneNumber TEXT NOT NULL,
		email TEXT,
		address TEXT,
		company TEXT
	);`
	err := db.createTable(sql)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) createTable(sql string) error {
	// create table
	if _, err := db.sqliteDB.Exec(sql); err != nil {
		return fmt.Errorf("create table failed: %w", err)
	}
	return nil
}

func (db *DB) CreateUser(user *User) (int64, error) {
	query := `
	INSERT INTO users (userId, username, password, phoneNumber, email, address, company) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	// perform insert user
	result, err := db.sqliteDB.Exec(query, user.UserId, user.Username, user.Password, user.PhoneNumber, user.Email, user.Address, user.Company)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("email already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}
