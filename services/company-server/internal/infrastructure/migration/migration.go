// Файл company-server/internal/infrastructure/migration/migration.go содержит реализацию пакета migration.
package migration

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed sql/*.sql
var migrationsFS embed.FS

// Migrator handles database migrations using golang-migrate
type Migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewMigrator creates a new Migrator instance
func NewMigrator(db *sql.DB, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// Run executes all pending migrations
func (m *Migrator) Run() error {
	m.logger.Info("Starting database migrations")

	// Create source driver from embedded files
	sourceDriver, err := iofs.New(migrationsFS, "sql")
	if err != nil {
		return err
	}

	// Create database driver
	dbDriver, err := postgres.WithInstance(m.db, &postgres.Config{
		SchemaName: "companies",
	})
	if err != nil {
		return err
	}

	// Create migrate instance
	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}

	// Run migrations
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	version, dirty, _ := migrator.Version()
	m.logger.Info("Database migrations completed",
		zap.Uint("version", version),
		zap.Bool("dirty", dirty))

	return nil
}
