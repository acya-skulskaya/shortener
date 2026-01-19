package shorturlindb

import (
	"context"
	"errors"
	"fmt"

	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (repo *InDBShortURLRepository) StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList, userID string) (listShorten []jsonModel.BatchURLList, err error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		logger.Log.Debug("could not start transaction",
			zap.Error(err),
		)
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO short_urls (id, short_url, original_url, user_id) VALUES ($1, $2, $3, $4)")
	if err != nil {
		logger.Log.Debug("could not prepare statement",
			zap.Error(err),
		)
		return nil, err
	}
	defer stmt.Close()

	var errs []error

	for _, item := range listOriginal {
		listShortenItem := jsonModel.BatchURLList{
			CorrelationID: item.CorrelationID,
			ShortURL:      config.Values.URLAddress + "/" + item.CorrelationID,
		}

		_, err := stmt.ExecContext(ctx, item.CorrelationID, config.Values.URLAddress+"/"+item.CorrelationID, item.OriginalURL, userID)
		if err != nil {
			// если ошибка, то откатываем изменения
			tx.Rollback()
			logger.Log.Debug("could not insert row",
				zap.Error(err),
				zap.Any("item", item),
			)

			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				listShortenItem.Err = fmt.Sprint(err)
				if pgErr.Code == pgerrcode.UniqueViolation {
					err = errorsInternal.ErrConflictOriginalURL

					row := repo.DB.QueryRowContext(ctx,
						"SELECT id FROM short_urls where original_url = $1", item.OriginalURL)
					var id string
					errScan := row.Scan(&id) // разбираем результат
					if errScan != nil {
						logger.Log.Debug("could not get row from db",
							zap.Error(errScan),
						)
						return nil, errScan
					}
					listShortenItem.CorrelationID = id
					listShortenItem.ShortURL = config.Values.URLAddress + "/" + id
				} else {
					err = errorsInternal.ErrConflictID

					row := repo.DB.QueryRowContext(ctx,
						"SELECT original_url FROM short_urls where id = $1", item.CorrelationID)
					var originalURL string
					errScan := row.Scan(&originalURL) // разбираем результат
					if errScan != nil {
						logger.Log.Debug("could not get row from db",
							zap.Error(errScan),
						)
						return nil, errScan
					}
					listShortenItem.OriginalURL = originalURL
				}

				errs = append(errs, err)
			}

			listShorten = []jsonModel.BatchURLList{listShortenItem}
			return listShorten, errors.Join(errs...)
		}

		listShorten = append(listShorten, listShortenItem)
	}

	// завершаем транзакцию
	tx.Commit()

	return listShorten, errors.Join(errs...)
}
