package api

import (
	"embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

var (
	//go:embed schema/sqlite
	schemaFS embed.FS

	logger = log.New(os.Stdout, "DB: ", log.LstdFlags|log.Lshortfile)
)

func GetDB() (*sqlx.DB, error) {
	dbLoc := "app.sqlite"
	if userSuppliedLoc, ok := os.LookupEnv("DB_LOCATION"); ok {
		dbLoc = userSuppliedLoc
	}

	reInitializeDb := false
	if envReinitFlag, ok := os.LookupEnv(("DB_REINIT")); ok {
		reInitializeDb = strings.ToLower(envReinitFlag) == "true"
	}

	db, err := sqlx.Connect("sqlite", "file:"+dbLoc+"?mode=rwc&_foreign_keys=true&_journal_mode=WAL")
	if err != nil {
		logger.Printf("Failed to connect to database at %s", dbLoc)
		return nil, err
	}

	if _, err := os.Stat(dbLoc); errors.Is(err, os.ErrNotExist) || reInitializeDb {
		logMsg := "No database found at %s, attempting to create..."
		if reInitializeDb {
			logMsg = "Reinitializing database at %s as requested..."
		}
		logger.Printf(logMsg, dbLoc)
		if err = initializeDb(db); err != nil {
			logger.Printf("Could not initialize new database at %s", dbLoc)
			return nil, err
		}
		logger.Printf("...finished creating database at %s", dbLoc)
	}

	logger.Printf("using database at %s", dbLoc)
	return db, nil

}

func initializeDb(db *sqlx.DB) error {
	files, err := schemaFS.ReadDir("schema/sqlite")
	if err != nil {
		return err
	}

	for _, file := range files {

		filename := "schema/sqlite/" + file.Name()

		content, err := schemaFS.ReadFile(filename)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute script %s: %w", filename, err)
		}
	}

	return nil
}
