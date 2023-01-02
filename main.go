package main

import "aboutMeStoreService/repository"

func main() {

	migrator := repository.MakeMigrator("sqlite3", "file:aboutMeDB.db", "./repository/migrations")

	migrator.UpToLastVersion()

	migrator.CloseConnection()
}
