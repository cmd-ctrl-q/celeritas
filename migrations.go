package celeritas

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func (c *Celeritas) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}

	return nil
}

func (c *Celeritas) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Down(); err != nil {
		return err
	}

	return nil
}

// Steps runs migrations.
func (c *Celeritas) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Steps(n); err != nil {
		return err
	}

	return nil
}

// MigrateForce force removes the previous migration
func (c *Celeritas) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+c.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	// force remove the last migration
	if err = m.Force(-1); err != nil {
		return err
	}

	return nil
}
