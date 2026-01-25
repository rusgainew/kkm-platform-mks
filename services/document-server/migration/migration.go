package migration

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

//go:embed *.sql
var migrationFS embed.FS

type Migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewMigrator создает новый миграционный объект
func NewMigrator(db *sql.DB, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// Run выполняет все миграции
func (m *Migrator) Run() error {
	d, err := iofs.New(migrationFS, ".")
	if err != nil {
		m.logger.Error("Failed to create migration source", zap.Error(err))
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", d, "postgres://")
	if err != nil {
		m.logger.Error("Failed to create migrator", zap.Error(err))
		return err
	}

	if err := migrator.Up(); err != nil {
		if err == migrate.ErrNoChange {
			m.logger.Info("No migrations to run")
			return nil
		}
		m.logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("migration failed: %w", err)
	}

	m.logger.Info("Migrations completed successfully")
	return nil
}
