package shorturlindb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (repo *InDBShortURLRepository) Store(originalURL string) (id string, err error) {
	id = helpers.RandStringRunes(10)

	ctx := context.Background()

	_, err = repo.DB.ExecContext(ctx,
		"INSERT INTO short_urls (id, short_url, original_url) VALUES ($1, $2, $3)",
		id, config.Values.URLAddress+"/"+id, originalURL)
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

func (repo *InDBShortURLRepository) StoreBatch(listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList, err error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		logger.Log.Debug("could not start transaction",
			zap.Error(err),
		)
		return nil, err
	}

	ctx := context.Background()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO short_urls (id, short_url, original_url) VALUES ($1, $2, $3)")
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

		_, err := stmt.ExecContext(ctx, item.CorrelationID, item.ShortURL, item.OriginalURL)
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
