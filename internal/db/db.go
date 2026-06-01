package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schemaCreateTable = `
CREATE TABLE scheduler(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
`
	schemaCreateIndex = `
	CREATE INDEX idx_date ON scheduler (date);
`
)

var DB *sql.DB

// A function for defining a database
func Init(dbFile string) error {
	install := false
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if install {
		_, err = DB.Exec(schemaCreateTable)
		if err != nil {
			DB.Close()
			return err
		}
		_, err = DB.Exec(schemaCreateIndex)
		if err != nil {
			DB.Close()
			return err
		}
	}

	return nil
}
