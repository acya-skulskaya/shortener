package shorturlindb

import (
	"database/sql"

	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type InDBShortURLRepository struct {
	DB *sql.DB
}

func NewInDBShortURLRepository(databaseDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.Values.DatabaseDSN)
	if err != nil {
		logger.Log.Debug("could not open db connection",
			zap.Error(err),
		)
		return nil, err
	}

	if err = runMigrations(); err != nil {
		logger.Log.Error("Failed to run migrations",
			zap.Error(err),
		)
		return nil, err
	}

	return db, nil
}

func runMigrations() error {
	dbDSN := config.Values.DatabaseDSN

	m, err := migrate.New(
		"file://./migrations",
		//"file://./../../migrations",
		dbDSN,
	)
	if err != nil {
		logger.Log.Warn("failed to initialize migrations",
			zap.Error(err),
			zap.String("dbDSN", dbDSN),
		)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Warn("failed to apply migrations",
			zap.Error(err),
		)
		return err
	}

	logger.Log.Debug("Migrations applied successfully!",
		zap.Error(err),
	)

	return nil
}
