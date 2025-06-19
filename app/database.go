package app

import (
	"context"
	"database/sql"
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

func (db *DB) Close() error {
	// close sqlite3 database
	if err := db.sqliteDB.Close(); err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateTables() error {
	// create user table
	sql := `
		CREATE TABLE IF NOT EXISTS users (
		user_id TEXT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		email TEXT,
		address TEXT,
		company TEXT
	);`
	err := db.createUserTable(sql)
	if err != nil {
		return err
	}
	// create question single-choice table
	sql = `CREATE TABLE IF NOT EXISTS single_choice (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		answers TEXT NOT NULL,
		standard_answer TEXT NOT NULL,
	);`
	err = db.createQuestionSingleChoiceTable(sql)
	return nil
}

func (db *DB) createUserTable(sql string) error {
	// create user table
	if _, err := db.sqliteDB.Exec(sql); err != nil {
		return fmt.Errorf("create user table failed: %w", err)
	}
	return nil
}

func (db *DB) CreateUser(user *User) (int64, error) {
	// execute user sql
	query := `
	INSERT INTO users (user_id, username, password, phone_number, email, address, company) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	// perform insert user
	result, err := db.sqliteDB.Exec(query, user.UserId, user.Username, user.Password, user.PhoneNumber, user.Email, user.Address, user.Company)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("user already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) CreateUserContext(ctx context.Context, user *User) (int64, error) {
	// create user sql
	query := `
	INSERT INTO users (user_id, username, password, phone_number, email, address, company) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	// execute create user
	result, err := db.sqliteDB.ExecContext(ctx, query, user.UserId, user.Username, user.Password, user.PhoneNumber, user.Email, user.Address, user.Company)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("user already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) QueryUser(userId string) (*User, error) {
	// query user sql
	query := `
	SELECT user_id, username, password, phone_number, email, address, company
	FROM users WHERE user_id = ?
	`
	// execute query user
	row := db.sqliteDB.QueryRow(query, userId)
	user := &User{}
	err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.PhoneNumber, &user.Email, &user.Address, &user.Company)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (db *DB) QueryUserContext(ctx context.Context, userId string) (*User, error) {
	// query user sql
	query := `
	SELECT user_id, username, password, phone_number, email, address, company
	FROM users WHERE user_id = ?
	`
	// execute query user
	row := db.sqliteDB.QueryRowContext(ctx, query, userId)
	user := &User{}
	err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.PhoneNumber, &user.Email, &user.Address, &user.Company, &user.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (db *DB) UpdateUser(user *User) error {
	// update user sql
	query := `
	UPDATE users 
	SET username = ?, password = ?, phone_number = ?, email = ?, address = ?, company = ?
	WHERE user_id = ?
	`
	// execute update user
	result, err := db.sqliteDB.Exec(query, user.Username, user.Password, user.PhoneNumber, user.Email, user.Address, user.Company, user.UserId)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (db *DB) UpdateUserContext(ctx context.Context, user *User) error {
	// update user sql
	query := `
	UPDATE users 
	SET username = ?, password = ?, phone_number = ?, email = ?, address = ?, company = ?
	WHERE user_id = ?
	`
	// execute update user
	result, err := db.sqliteDB.ExecContext(ctx, query, user.Username, user.Password, user.PhoneNumber, user.Email, user.Address, user.Company, user.UserId)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (db *DB) DeleteUser(userId string) error {
	// update user sql
	query := `DELETE FROM users WHERE user_id = ?`
	// execute delete user
	result, err := db.sqliteDB.Exec(query, userId)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (db *DB) DeleteUserContext(ctx context.Context, userId string) error {
	// update user sql
	query := `DELETE FROM users WHERE user_id = ?`
	// execute delete user
	result, err := db.sqliteDB.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (db *DB) QueryUsers() ([]*User, error) {
	// query users
	query := `
	SELECT user_id, username, password, phone_number, email, address, company
	FROM users
	`
	// execute query users
	rows, err := db.sqliteDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// fetch users from database
	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.UserId, &user.Username, &user.Password, &user.PhoneNumber, &user.Email, &user.Address, &user.Company); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *DB) QueryUsersContext(ctx context.Context) ([]*User, error) {
	// query users
	query := `
	SELECT user_id, username, password, phone_number, email, address, company
	FROM users
	`
	// execute query users
	rows, err := db.sqliteDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// fetch users from database
	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.UserId, &user.Username, &user.Password, &user.PhoneNumber, &user.Email, &user.Address, &user.Company); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *DB) createQuestionSingleChoiceTable(sql string) error {
	// create question table
	if _, err := db.sqliteDB.Exec(sql); err != nil {
		return fmt.Errorf("create question single-choice table failed: %w", err)
	}
	return nil
}
