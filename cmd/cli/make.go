package main

import (
	"errors"
	"fmt"
	"time"
)

func doMake(arg2, arg3 string) error {

	switch arg2 {
	case "migration":
		dbType := cel.DB.DataType
		if arg3 == "" {
			exitGraceFully(errors.New("you must give the migration name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		// create temlates for migrations
		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGraceFully(err)
		}

		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGraceFully(err)
		}
	}

	return nil
}
