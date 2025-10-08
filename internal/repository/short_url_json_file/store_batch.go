package shorturljsonfile

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

func (repo *JSONFileShortURLRepository) StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList, userID string) (listShorten []jsonModel.BatchURLList, err error) {
	writer, err := NewFileWriter(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create file writer",
			zap.Error(err),
		)
		return nil, err
	}

	reader, err := NewFileReader(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create reader",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return nil, err
	}
	defer reader.Close()
	existingRows, err := reader.ReadFile(repo)

	var rows []jsonModel.URLList

	var errs []error

	for _, item := range listOriginal {
		err = nil
		for _, existingRow := range existingRows {
			if existingRow.ID == item.CorrelationID {
				err = errorsInternal.ErrConflictID
			}
			if existingRow.OriginalURL == item.OriginalURL {
				err = errorsInternal.ErrConflictOriginalURL
			}
		}

		listShortenItem := jsonModel.BatchURLList{
			CorrelationID: item.CorrelationID,
			ShortURL:      config.Values.URLAddress + "/" + item.CorrelationID,
		}

		if err != nil {
			errs = append(errs, err)
			listShortenItem.Err = fmt.Sprint(err)
		} else {
			rows = append(rows, jsonModel.URLList{
				ID:          item.CorrelationID,
				OriginalURL: item.OriginalURL,
				ShortURL:    config.Values.URLAddress + "/" + item.CorrelationID,
				UserID:      userID,
			})
		}

		listShorten = append(listShorten, listShortenItem)
	}

	err = writer.WriteFileRows(repo, rows)
	if err != nil {
		logger.Log.Debug("could not write short urls to file",
			zap.Error(err),
			zap.Any("rows", rows),
		)
		return nil, err
	}

	return listShorten, errors.Join(errs...)
}
