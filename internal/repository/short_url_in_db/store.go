package shorturlindb

import (
	"context"
	"errors"

	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (repo *InDBShortURLRepository) Store(ctx context.Context, originalURL string, userID string) (id string, err error) {
	id = helpers.RandStringRunes(10)

	_, err = repo.DB.ExecContext(ctx,
		"INSERT INTO short_urls (id, short_url, original_url, user_id) VALUES ($1, $2, $3, $4)",
		id, config.Values.URLAddress+"/"+id, originalURL, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = errorsInternal.ErrConflictOriginalURL

			row := repo.DB.QueryRowContext(ctx,
				"SELECT id FROM short_urls where original_url = $1", originalURL)

			errScan := row.Scan(&id) // разбираем результат
			if errScan != nil {
				logger.Log.Debug("could not get row from db",
					zap.Error(errScan),
				)
				return "", errScan
			}
			return id, err

		}
		logger.Log.Debug("could not insert into db",
			zap.Error(err),
		)
		return "", err
	}

	return id, nil
}
