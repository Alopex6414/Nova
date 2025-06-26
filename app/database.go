package app

import (
	"context"
	"database/sql"
	"encoding/json"
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
	// create single-choice question table
	sql = `CREATE TABLE IF NOT EXISTS single_choice (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		answers TEXT NOT NULL,
		standard_answer TEXT NOT NULL
	);`
	err = db.createQuestionSingleChoiceTable(sql)
	if err != nil {
		return err
	}
	// create multiple-choice question table
	sql = `CREATE TABLE IF NOT EXISTS multiple_choice (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		answers TEXT NOT NULL,
		standard_answers TEXT NOT NULL
	);`
	err = db.createQuestionMultipleChoiceTable(sql)
	if err != nil {
		return err
	}
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
	// create single-choice table
	if _, err := db.sqliteDB.Exec(sql); err != nil {
		return fmt.Errorf("create single-choice question table failed: %w", err)
	}
	return nil
}

func (db *DB) CreateQuestionSingleChoice(question *QuestionSingleChoice) (int64, error) {
	// execute single-choice sql
	query := `
	INSERT INTO single_choice (id, title, answers, standard_answer) 
	VALUES (?, ?, ?, ?)
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return 0, err
	}
	standardAnswer, err := json.Marshal(question.StandardAnswer)
	if err != nil {
		return 0, err
	}
	// perform insert single-choice
	result, err := db.sqliteDB.Exec(query, question.Id, question.Title, answers, standardAnswer)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("single-choice question already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) CreateQuestionSingleChoiceContext(ctx context.Context, question *QuestionSingleChoice) (int64, error) {
	// execute single-choice sql
	query := `
	INSERT INTO single_choice (id, title, answers, standard_answer) 
	VALUES (?, ?, ?, ?)
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return 0, err
	}
	standardAnswer, err := json.Marshal(question.StandardAnswer)
	if err != nil {
		return 0, err
	}
	// perform insert single-choice
	result, err := db.sqliteDB.ExecContext(ctx, query, question.Id, question.Title, answers, standardAnswer)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("single-choice question already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) QueryQuestionSingleChoice(id string) (*QuestionSingleChoice, error) {
	// query single-choice sql
	query := `
	SELECT id, title, answers, standard_answer
	FROM single_choice WHERE id = ?
	`
	// variables definition
	var answers []byte
	var standardAnswer []byte
	// execute query single-choice
	row := db.sqliteDB.QueryRow(query, id)
	question := &QuestionSingleChoice{}
	err := row.Scan(&question.Id, &question.Title, &answers, &standardAnswer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("single-choice question not found")
		}
		return nil, err
	}
	// unmarshal json slices & structure
	if err := json.Unmarshal(answers, &question.Answers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(standardAnswer, &question.StandardAnswer); err != nil {
		return nil, err
	}
	return question, nil
}

func (db *DB) QueryQuestionSingleChoiceContext(ctx context.Context, id string) (*QuestionSingleChoice, error) {
	// query single-choice sql
	query := `
	SELECT id, title, answers, standard_answer
	FROM single_choice WHERE id = ?
	`
	// variables definition
	var answers []byte
	var standardAnswer []byte
	// execute query single-choice
	row := db.sqliteDB.QueryRowContext(ctx, query, id)
	question := &QuestionSingleChoice{}
	err := row.Scan(&question.Id, &question.Title, &answers, &standardAnswer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("single-choice question not found")
		}
		return nil, err
	}
	// unmarshal json slices & structure
	if err := json.Unmarshal(answers, &question.Answers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(standardAnswer, &question.StandardAnswer); err != nil {
		return nil, err
	}
	return question, nil
}

func (db *DB) UpdateQuestionSingleChoice(question *QuestionSingleChoice) error {
	// update single-choice sql
	query := `
	UPDATE single_choice 
	SET title = ?, answers = ?, standard_answer = ?
	WHERE id = ?
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return err
	}
	standardAnswer, err := json.Marshal(question.StandardAnswer)
	if err != nil {
		return err
	}
	// execute update single-choice
	result, err := db.sqliteDB.Exec(query, question.Title, answers, standardAnswer, question.Id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("single-choice question not found")
	}
	return nil
}

func (db *DB) UpdateQuestionSingleChoiceContext(ctx context.Context, question *QuestionSingleChoice) error {
	// update single-choice sql
	query := `
	UPDATE single_choice 
	SET title = ?, answers = ?, standard_answer = ?
	WHERE id = ?
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return err
	}
	standardAnswer, err := json.Marshal(question.StandardAnswer)
	if err != nil {
		return err
	}
	// execute update single-choice
	result, err := db.sqliteDB.ExecContext(ctx, query, question.Title, answers, standardAnswer, question.Id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("single-choice question not found")
	}
	return nil
}

func (db *DB) DeleteQuestionSingleChoice(id string) error {
	// update single-choice sql
	query := `DELETE FROM single_choice WHERE id = ?`
	// execute delete single-choice
	result, err := db.sqliteDB.Exec(query, id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("single-choice question not found")
	}
	return nil
}

func (db *DB) DeleteQuestionSingleChoiceContext(ctx context.Context, id string) error {
	// update single-choice sql
	query := `DELETE FROM single_choice WHERE id = ?`
	// execute delete single-choice
	result, err := db.sqliteDB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("single-choice question not found")
	}
	return nil
}

func (db *DB) QueryQuestionsSingleChoice() ([]*QuestionSingleChoice, error) {
	// query single-choice questions
	query := `
	SELECT id, title, answers, standard_answer
	FROM single_choice
	`
	// execute query single-choice questions
	rows, err := db.sqliteDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// fetch single-choice questions from database
	var questions []*QuestionSingleChoice
	for rows.Next() {
		// variables definition
		var answers []byte
		var standardAnswer []byte
		// query single-choice question
		question := &QuestionSingleChoice{}
		if err := rows.Scan(&question.Id, &question.Title, &answers, &standardAnswer); err != nil {
			return nil, err
		}
		// unmarshal json slices & structure
		if err := json.Unmarshal(answers, &question.Answers); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(standardAnswer, &question.StandardAnswer); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

func (db *DB) QueryQuestionsSingleChoiceContext(ctx context.Context) ([]*QuestionSingleChoice, error) {
	// query single-choice questions
	query := `
	SELECT id, title, answers, standard_answer
	FROM single_choice
	`
	// execute query single-choice questions
	rows, err := db.sqliteDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// fetch single-choice questions from database
	var questions []*QuestionSingleChoice
	for rows.Next() {
		// variables definition
		var answers []byte
		var standardAnswer []byte
		// query single-choice question
		question := &QuestionSingleChoice{}
		if err := rows.Scan(&question.Id, &question.Title, &answers, &standardAnswer); err != nil {
			return nil, err
		}
		// unmarshal json slices & structure
		if err := json.Unmarshal(answers, &question.Answers); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(standardAnswer, &question.StandardAnswer); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

func (db *DB) createQuestionMultipleChoiceTable(sql string) error {
	// create multiple-choice table
	if _, err := db.sqliteDB.Exec(sql); err != nil {
		return fmt.Errorf("create multiple-choice question table failed: %w", err)
	}
	return nil
}

func (db *DB) CreateQuestionMultipleChoice(question *QuestionMultipleChoice) (int64, error) {
	// execute multiple-choice sql
	query := `
	INSERT INTO multiple_choice (id, title, answers, standard_answers) 
	VALUES (?, ?, ?, ?)
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return 0, err
	}
	standardAnswers, err := json.Marshal(question.StandardAnswers)
	if err != nil {
		return 0, err
	}
	// perform insert multiple-choice
	result, err := db.sqliteDB.Exec(query, question.Id, question.Title, answers, standardAnswers)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("multiple-choice question already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) CreateQuestionMultipleChoiceContext(ctx context.Context, question *QuestionMultipleChoice) (int64, error) {
	// execute multiple-choice sql
	query := `
	INSERT INTO multiple_choice (id, title, answers, standard_answers) 
	VALUES (?, ?, ?, ?)
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return 0, err
	}
	standardAnswers, err := json.Marshal(question.StandardAnswers)
	if err != nil {
		return 0, err
	}
	// perform insert multiple-choice
	result, err := db.sqliteDB.ExecContext(ctx, query, question.Id, question.Title, answers, standardAnswers)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return 0, fmt.Errorf("multiple-choice question already exists")
			}
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) QueryQuestionMultipleChoice(id string) (*QuestionMultipleChoice, error) {
	// query multiple-choice sql
	query := `
	SELECT id, title, answers, standard_answers
	FROM multiple_choice WHERE id = ?
	`
	// variables definition
	var answers []byte
	var standardAnswers []byte
	// execute query multiple-choice
	row := db.sqliteDB.QueryRow(query, id)
	question := &QuestionMultipleChoice{}
	err := row.Scan(&question.Id, &question.Title, &answers, &standardAnswers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("multiple-choice question not found")
		}
		return nil, err
	}
	// unmarshal json slices & structure
	if err := json.Unmarshal(answers, &question.Answers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(standardAnswers, &question.StandardAnswers); err != nil {
		return nil, err
	}
	return question, nil
}

func (db *DB) QueryQuestionMultipleChoiceContext(ctx context.Context, id string) (*QuestionMultipleChoice, error) {
	// query multiple-choice sql
	query := `
	SELECT id, title, answers, standard_answers
	FROM multiple_choice WHERE id = ?
	`
	// variables definition
	var answers []byte
	var standardAnswers []byte
	// execute query multiple-choice
	row := db.sqliteDB.QueryRowContext(ctx, query, id)
	question := &QuestionMultipleChoice{}
	err := row.Scan(&question.Id, &question.Title, &answers, &standardAnswers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("multiple-choice question not found")
		}
		return nil, err
	}
	// unmarshal json slices & structure
	if err := json.Unmarshal(answers, &question.Answers); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(standardAnswers, &question.StandardAnswers); err != nil {
		return nil, err
	}
	return question, nil
}

func (db *DB) UpdateQuestionMultipleChoice(question *QuestionMultipleChoice) error {
	// update multiple-choice sql
	query := `
	UPDATE multiple_choice 
	SET title = ?, answers = ?, standard_answers = ?
	WHERE id = ?
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return err
	}
	standardAnswers, err := json.Marshal(question.StandardAnswers)
	if err != nil {
		return err
	}
	// execute update multiple-choice
	result, err := db.sqliteDB.Exec(query, question.Title, answers, standardAnswers, question.Id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("multiple-choice question not found")
	}
	return nil
}

func (db *DB) UpdateQuestionMultipleChoiceContext(ctx context.Context, question *QuestionMultipleChoice) error {
	// update multiple-choice sql
	query := `
	UPDATE multiple_choice 
	SET title = ?, answers = ?, standard_answers = ?
	WHERE id = ?
	`
	// marshal json slices & structure
	answers, err := json.Marshal(question.Answers)
	if err != nil {
		return err
	}
	standardAnswers, err := json.Marshal(question.StandardAnswers)
	if err != nil {
		return err
	}
	// execute update multiple-choice
	result, err := db.sqliteDB.ExecContext(ctx, query, question.Title, answers, standardAnswers, question.Id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("multiple-choice question not found")
	}
	return nil
}

func (db *DB) DeleteQuestionMultipleChoice(id string) error {
	// update multiple-choice sql
	query := `DELETE FROM multiple_choice WHERE id = ?`
	// execute delete multiple-choice
	result, err := db.sqliteDB.Exec(query, id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("multiple-choice question not found")
	}
	return nil
}

func (db *DB) DeleteQuestionMultipleChoiceContext(ctx context.Context, id string) error {
	// update multiple-choice sql
	query := `DELETE FROM multiple_choice WHERE id = ?`
	// execute delete multiple-choice
	result, err := db.sqliteDB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("multiple-choice question not found")
	}
	return nil
}

func (db *DB) QueryQuestionsMultipleChoice() ([]*QuestionMultipleChoice, error) {
	// query multiple-choice questions
	query := `
	SELECT id, title, answers, standard_answers
	FROM multiple_choice
	`
	// execute query multiple-choice questions
	rows, err := db.sqliteDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// fetch multiple-choice questions from database
	var questions []*QuestionMultipleChoice
	for rows.Next() {
		// variables definition
		var answers []byte
		var standardAnswers []byte
		// query multiple-choice question
		question := &QuestionMultipleChoice{}
		if err := rows.Scan(&question.Id, &question.Title, &answers, &standardAnswers); err != nil {
			return nil, err
		}
		// unmarshal json slices & structure
		if err := json.Unmarshal(answers, &question.Answers); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(standardAnswers, &question.StandardAnswers); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

func (db *DB) QueryQuestionsMultipleChoiceContext(ctx context.Context) ([]*QuestionMultipleChoice, error) {
	// query multiple-choice questions
	query := `
	SELECT id, title, answers, standard_answers
	FROM multiple_choice
	`
	// execute query multiple-choice questions
	rows, err := db.sqliteDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// fetch multiple-choice questions from database
	var questions []*QuestionMultipleChoice
	for rows.Next() {
		// variables definition
		var answers []byte
		var standardAnswers []byte
		// query multiple-choice question
		question := &QuestionMultipleChoice{}
		if err := rows.Scan(&question.Id, &question.Title, &answers, &standardAnswers); err != nil {
			return nil, err
		}
		// unmarshal json slices & structure
		if err := json.Unmarshal(answers, &question.Answers); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(standardAnswers, &question.StandardAnswers); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}
