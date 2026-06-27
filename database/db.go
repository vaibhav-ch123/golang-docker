package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SSLMODE string

const (
	SSLModeEnable  SSLMODE = "enable"
	SSLModeDisable SSLMODE = "disable"
)

var (
	Todo *sqlx.DB
)

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMODE) error {

	connectstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	DB, err := sqlx.Open("postgres", connectstr)

	if err != nil {
		return nil
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	Todo = DB
	return MigrateUp(DB)
}

func MigrateUp(db *sqlx.DB) error {

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../database/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func ShutDownDatabase() error {
	return Todo.Close()
}

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := Todo.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transation: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}

func ReplaceSQL(old string, oldPattern string) string{

	count := strings.Count(old, oldPattern)
	for i := 1 ; i<=count ; i++ {
		old = strings.Replace(old, oldPattern, "$"+strconv.Itoa(i), 1)
	}
	return old
}
