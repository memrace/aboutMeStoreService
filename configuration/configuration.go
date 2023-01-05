package configuration

type dbConnectionConfiguration struct {
	DriverName     string
	DataSourceName string
}

var DbConnectionConfiguration dbConnectionConfiguration = dbConnectionConfiguration{
	"sqlite3", "file:aboutMeDB.db",
}
