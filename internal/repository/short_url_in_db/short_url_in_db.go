package shorturlindb

import (
	"context"
	"database/sql"
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
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

	return db, nil
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
