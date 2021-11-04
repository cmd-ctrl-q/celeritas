package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// setup populates the Celeritas type
func setup() {
	err := godotenv.Load()
	if err != nil {
		exitGraceFully(err)
	}

	// get root path
	path, err := os.Getwd()
	if err != nil {
		exitGraceFully(err)
	}

	cel.RootPath = path
	cel.DB.DataType = os.Getenv("DATABASE_TYPE")
}

// getDSN returns a dsn in the correct format
func getDSN() string {
	dbType := cel.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}
	return "mysql://" + cel.BuildDSN()
}

func showHelp() {
	color.Yellow(`Available commands:

	help			- show the help commands
	version			- print application version
	migrate 		- run all up migrations that have have yet to run
	migrate down 		- reverses the most receive migration
	migrate reset 		- runs all down migrations in reverse order, and then all up migrations
	make migration <name>	- creates two new up and down migrations in the migrations folder
	`)
}