package main

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository/migrations"
)

func main() {

	migrator := migrations.NewMigrator(
		configuration.DbConnectionConfiguration.DriverName,
		configuration.DbConnectionConfiguration.DataSourceName,
		"./repository/migrations")

	migrator.UpToLastVersion()

	migrator.Close()
}
