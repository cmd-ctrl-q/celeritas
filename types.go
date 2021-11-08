package celeritas

import "database/sql"

type initPaths struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string

	// does it persist between browswer closes
	persist string

	// is the cookie encrypted
	secure string

	// the domain the cookie is associated with
	domain string
}

type databaseConfig struct {
	dsn      string
	database string
}

type Database struct {
	// DataType is the database type
	DataType string

	// Pool is the connection pool
	Pool *sql.DB
}

type redisConfig struct {
	host     string
	password string

	// prefix will add a prefix to a key to prevent name conflicts between multiple applications
	prefix string
}
