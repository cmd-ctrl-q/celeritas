package main

import (
	"os"

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
