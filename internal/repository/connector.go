package repository

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnectorNoDbErr = errors.New("no available database connection")

type DbConnector interface {
	Connect(file string) error
	InitializeDb() error
	Exec(query string, params ...any) (sql.Result, error)
	QueryOne(query string, params ...any) *sql.Row
	QueryMany(query string, params ...any) (*sql.Rows, error)
}

type DbConnection struct {
	db *sql.DB
}

var _ DbConnector = &DbConnection{}

func (conn *DbConnection) Connect(file string) error {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return err
	}

	conn.db = db

	return nil
}

func (conn *DbConnection) InitializeDb() error {
	if conn.db == nil {
		return DbConnectorNoDbErr
	}

	if _, err := conn.db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return err
	}

	tx, beginErr := conn.db.Begin()

	if beginErr != nil {
		return beginErr
	}

	defer tx.Rollback()

	if _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS "task" ("id" INTEGER PRIMARY KEY ASC, "name" TEXT UNIQUE)`); err != nil {
		return err
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS "name_idx" ON "task"("name");`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS "time_range" (
			"id" INTEGER PRIMARY KEY ASC,
			"task_fk" INTEGER NOT NULL,
			"start" INTEGER,
			"end" INTEGER,
			CONSTRAINT task_fk
				FOREIGN KEY (task_fk)
				REFERENCES "task"("id")
				ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS "task_idx" ON "time_range"("task_fk");`); err != nil {
		return err
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		return commitErr
	}

	return nil
}

func (conn *DbConnection) Exec(query string, params ...any) (sql.Result, error) {
	return conn.db.Exec(query, params...)
}

func (conn *DbConnection) QueryOne(query string, params ...any) *sql.Row {
	return conn.db.QueryRow(query, params...)
}

func (conn *DbConnection) QueryMany(query string, params ...any) (*sql.Rows, error) {
	return conn.db.Query(query, params...)
}
