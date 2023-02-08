// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

// Package database provides persistent data storage using
// a postgres or sqlite3 database.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/harness/gitness/internal/store/database/migrate"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	// sqlForUpdate is the sql statement used for locking rows returned by select queries.
	sqlForUpdate = "FOR UPDATE"
)

// build is a global instance of the sql builder. we are able to
// hardcode to postgres since sqlite3 is compatible with postgres.
var builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// Connect to a database and verify with a ping.
func Connect(ctx context.Context, driver string, datasource string) (*sqlx.DB, error) {
	datasource, err := prepareDatasourceForDriver(driver, datasource)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare datasource: %w", err)
	}

	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, fmt.Errorf("failed to open the db: %w", err)
	}

	dbx := sqlx.NewDb(db, driver)
	if err = pingDatabase(ctx, dbx); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}

	return dbx, nil
}

// ConnectAndMigrate creates the database handle and migrates the database.
func ConnectAndMigrate(ctx context.Context, driver string, datasource string) (*sqlx.DB, error) {
	dbx, err := Connect(ctx, driver, datasource)
	if err != nil {
		return nil, err
	}

	if err = migrateDatabase(ctx, dbx); err != nil {
		return nil, fmt.Errorf("failed to setup the db: %w", err)
	}

	return dbx, nil
}

// Must is a helper function that wraps a call to Connect
// and panics if the error is non-nil.
func Must(db *sqlx.DB, err error) *sqlx.DB {
	if err != nil {
		panic(err)
	}
	return db
}

// prepareDatasourceForDriver ensures that required features are enabled on the
// datasource connection string based on the driver.
func prepareDatasourceForDriver(driver string, datasource string) (string, error) {
	switch driver {
	case "sqlite3":
		url, err := url.Parse(datasource)
		if err != nil {
			return "", fmt.Errorf("datasource is of invalid format for driver sqlite3")
		}

		// get original query and update it with required settings
		query := url.Query()

		// ensure foreign keys are always enabled (disabled by default)
		// See https://github.com/mattn/go-sqlite3#connection-string
		query.Set("_foreign_keys", "on")

		// update url with updated query
		url.RawQuery = query.Encode()

		return url.String(), nil
	default:
		return datasource, nil
	}
}

// helper function to ping the database with backoff to ensure
// a connection can be established before we proceed with the
// database setup and migration.
func pingDatabase(ctx context.Context, db *sqlx.DB) error {
	var err error
	for i := 1; i <= 30; i++ {
		err = db.PingContext(ctx)

		// No point in continuing if context was cancelled
		if errors.Is(err, context.Canceled) {
			return err
		}

		// We can complete on first successful ping
		if err == nil {
			return nil
		}

		log.Debug().Err(err).Msgf("Ping attempt #%d failed", i)

		time.Sleep(time.Second)
	}

	return fmt.Errorf("all 30 tries failed, last failure: %w", err)
}

// helper function to setup the database by performing automated
// database migration steps.
func migrateDatabase(ctx context.Context, db *sqlx.DB) error {
	return migrate.Migrate(ctx, db)
}
