package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL,
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT,
	repeat VARCHAR(128)
);

CREATE INDEX idx_date ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {

	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}

	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil

}

