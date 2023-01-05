package configuration

type DbConnectionConfiguration struct {
	DriverName     string
	DataSourceName string
	MigrationsPath string
}

var DbConnConfiguration = DbConnectionConfiguration{
	"sqlite3", "file:aboutMeDB.db", "",
}

func DbTestConnectionConfiguration(mPath string) DbConnectionConfiguration {
	return DbConnectionConfiguration{
		"sqlite3", "file:test.db?cache=shared&mode=memory", mPath,
	}
}
