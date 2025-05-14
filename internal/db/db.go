package db

import (
	"fmt"
	"log"
	"os"

	"github.com/LoaltyProgramm/to-do-app/internal/config"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//дописать создание таблицы в бд
const (
	schema = `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "", 
		title VARCHAR(32) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat CHAR(8) NOT NULL DEFAULT ""
	);
	CREATE INDEX idx_data ON scheduler (date);`
)

var DB *sqlx.DB

func InitDB(cfg *config.Config) error {
	dbFile := cfg.DatabaseFile

	var err error

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	if install {
		
		DB, err = sqlx.Open("sqlite", dbFile)
		if err != nil {
			return fmt.Errorf("db is not open: %w", err)
		}

		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("error create exec: %w", err)
		}

		log.Printf("a database with the scheduler table has been created, the path to the database: %w", dbFile)
	}

	DB, err = sqlx.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("db is not open: %w", err)
	}

	return nil
}
