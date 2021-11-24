package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// setup populates the Celeritas type
func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "vesion" && arg1 != "help" {

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
	make auth 		- creates and runs migrations for authentication tables and creates models and middleware
	make handler <name> 	- creates a stub handler in the handlers directory
	make model <name>	- creates a new model in the data directory
	make session 		- creates a table in the database as a session store
	make mail <name> 	- creates two starter mail templates in the mail directory
	`)
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	// check for error
	if err != nil {
		return err
	}

	// check if current file is directory
	if fi.IsDir() {
		return nil
	}

	// only check go files
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	// have matching file
	if matched {
		// read file contents
		read, err := ioutil.ReadFile(path)
		if err != nil {
			exitGraceFully(err)
		}

		newContents := strings.Replace(string(read), "myapp", appURL, -1)

		// write the changed file
		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			exitGraceFully(err)
		}
	}

	return nil
}

func updateSource() {
	// walk entire project folder and subfolders
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		exitGraceFully(err)
	}
}
