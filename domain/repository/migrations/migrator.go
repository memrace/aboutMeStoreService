package migrations

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Migrator struct {
	Db             *sql.DB
	migrationsPath string
}

func New(driverName string, dataSourceName string, migrationsPath string) *Migrator {
	db, err := sql.Open(driverName, dataSourceName)
	handleErr(err)
	return &Migrator{Db: db, migrationsPath: migrationsPath}
}

func (migrator *Migrator) UpToLastVersion() error {

	instance, err := sqlite3.WithInstance(migrator.Db, &sqlite3.Config{})

	handleErr(err)

	mSrc, err := (&file.File{}).Open(migrator.migrationsPath)

	handleErr(err)

	m, err := migrate.NewWithInstance("file", mSrc, "sqlite3", instance)

	handleErr(err)

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		handleErr(err)
	} else {
		err = nil
	}

	v, d, err := m.Version()
	if err != nil {
		log.Fatal(err)
	} else {
		print("Current version: ", v, " dirty: ", d)
	}
	return err
}

func (migrator *Migrator) Close() {
	migrator.Db.Close()
}
