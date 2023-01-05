package main

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository/migrations"
)

func main() {

	migrator := migrations.New(
		configuration.DbConnConfiguration.DriverName,
		configuration.DbConnConfiguration.DataSourceName,
		"./repository/migrations")

	migrator.UpToLastVersion()

	migrator.Close()
}
