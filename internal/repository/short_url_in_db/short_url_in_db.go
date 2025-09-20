package shorturlindb

import (
	"context"
	"database/sql"
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
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

	if err := runMigrations(); err != nil {
		logger.Log.Error("Failed to run migrations",
			zap.Error(err),
		)
		panic(err)
	}

	return db, nil
}

func runMigrations() error {
	dbDSN := config.Values.DatabaseDSN

	m, err := migrate.New(
		"file://./migrations",
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

func (repo *InDBShortURLRepository) Get(id string) (originalURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx,
		"SELECT original_url FROM short_urls where id = $1", id)

	err := row.Scan(&originalURL) // разбираем результат
	if err != nil {
		logger.Log.Debug("could not get row from db",
			zap.Error(err),
		)
		return ""
	}

	return originalURL
}

func (repo *InDBShortURLRepository) Store(originalURL string) (id string) {
	id = helpers.RandStringRunes(10)

	_, err := repo.DB.ExecContext(context.Background(),
		"INSERT INTO short_urls (id, short_url, original_url) VALUES ($1, $2, $3)",
		id, config.Values.URLAddress+"/"+id, originalURL)
	if err != nil {
		logger.Log.Debug("could not insert into db",
			zap.Error(err),
		)
		return ""
	}

	return id
}

func (repo *InDBShortURLRepository) StoreBatch(listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList) {
	tx, err := repo.DB.Begin()
	if err != nil {
		logger.Log.Debug("could not start transaction",
			zap.Error(err),
		)
		return nil
	}

	ctx := context.Background()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO short_urls (id, short_url, original_url) VALUES ($1, $2, $3)")
	if err != nil {
		logger.Log.Debug("could not prepare statement",
			zap.Error(err),
		)
		return nil
	}
	defer stmt.Close()

	for _, item := range listOriginal {
		_, err := stmt.ExecContext(ctx, item.CorrelationId, item.ShortURL, item.OriginalURL)
		if err != nil {
			// если ошибка, то откатываем изменения
			tx.Rollback()
			logger.Log.Debug("could not insert row",
				zap.Error(err),
				zap.Any("item", item),
			)
			return nil
		}

		listShorten = append(listShorten, jsonModel.BatchURLList{
			CorrelationId: item.CorrelationId,
			ShortURL:      config.Values.URLAddress + "/" + item.CorrelationId,
		})
	}

	// завершаем транзакцию
	tx.Commit()

	return listShorten
}
