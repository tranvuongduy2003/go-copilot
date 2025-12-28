package migrations

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Runner struct {
	migrationsPath string
	databaseURL    string
}

func NewRunner(migrationsPath string, databaseURL string) *Runner {
	return &Runner{
		migrationsPath: migrationsPath,
		databaseURL:    databaseURL,
	}
}

func (runner *Runner) Up() error {
	migrator, err := runner.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (runner *Runner) Down() error {
	migrator, err := runner.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

func (runner *Runner) Steps(steps int) error {
	migrator, err := runner.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migration steps: %w", err)
	}

	return nil
}

func (runner *Runner) Version() (uint, bool, error) {
	migrator, err := runner.createMigrator()
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

func (runner *Runner) Force(version int) error {
	migrator, err := runner.createMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	return nil
}

func (runner *Runner) createMigrator() (*migrate.Migrate, error) {
	sourceURL := "file://" + runner.migrationsPath
	return migrate.New(sourceURL, runner.databaseURL)
}
