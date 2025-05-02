package main

import (
	"database/sql"
)

const MaxOpenConnects = 1

type SQLite3 struct {
	db *sql.DB
}

func NewSQLite3(path string) (*SQLite3, error) {
	// open sqlite3 database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// setup database configure
	db.SetMaxOpenConns(MaxOpenConnects)
	return &SQLite3{db}, nil
}

func (s *SQLite3) Close() error {
	return s.db.Close()
}

func (s *SQLite3) CreateTable(name string) error {
	query := "CREATE TABLE IF NOT EXISTS `" + name + "` (" + ")`"
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return err
}
