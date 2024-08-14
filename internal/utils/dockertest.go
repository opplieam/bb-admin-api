package utils

import (
	"database/sql"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	postgresUser     = "test"
	postgresPassword = "test"
	postgresDb       = "buy-better-test"
	hostPort         = "5435"
)

func CreateDockerTestContainer() (*sql.DB, *dockertest.Pool, *dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, err
	}
	err = pool.Client.Ping()
	if err != nil {
		return nil, nil, nil, err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2-alpine3.19",
		Env: []string{
			"POSTGRES_USER=" + postgresUser,
			"POSTGRES_PASSWORD=" + postgresPassword,
			"POSTGRES_DB=" + postgresDb,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {{HostIP: "0.0.0.0", HostPort: hostPort}},
		},
	})
	if err != nil {
		_ = pool.Purge(resource)
		return nil, nil, nil, err
	}
	var testDB *sql.DB
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("postgres", "postgresql://test:test@localhost:5435/buy-better-test?sslmode=disable")
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		return nil, nil, nil, err
	}

	// Tell docker to hard kill the container in 60 seconds
	err = resource.Expire(60)
	if err != nil {
		return nil, nil, nil, err
	}

	return testDB, pool, resource, nil
}

func MigrateDB(db *sql.DB, sourceURL string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(sourceURL, postgresDb, driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func SeedDataFromSQL(db *sql.DB, sqlLocation string) error {
	tableSQL, err := os.ReadFile(sqlLocation)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(tableSQL))
	if err != nil {
		return err
	}
	return nil
}
