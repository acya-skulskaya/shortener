package shorturljsonfile

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

type JSONFileShortURLRepository struct {
	FileStoragePath string
}

func (repo *JSONFileShortURLRepository) Get(ctx context.Context, id string) (originalURL string) {
	reader, err := NewFileReader(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create reader",
			zap.Error(err),
			zap.String("id", id),
			zap.String("file", repo.FileStoragePath),
		)
		return ""
	}
	defer reader.Close()

	list, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("id", id),
			zap.String("file", repo.FileStoragePath),
		)
		return ""
	}

	for _, l := range list {
		if l.ID == id {
			return l.OriginalURL
		}
	}

	return ""
}

func (repo *JSONFileShortURLRepository) Store(ctx context.Context, originalURL string) (id string, err error) {
	reader, err := NewFileReader(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create reader",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}
	defer reader.Close()
	existingRows, err := reader.ReadFile()

	for _, existingRow := range existingRows {
		if existingRow.OriginalURL == originalURL {
			return existingRow.ID, errorsInternal.ErrConflictOriginalURL
		}
	}

	id = helpers.RandStringRunes(10)

	row := jsonModel.URLList{
		ID:          id,
		ShortURL:    config.Values.URLAddress + "/" + id,
		OriginalURL: originalURL,
	}

	writer, err := NewFileWriter(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create file writer",
			zap.Error(err),
			zap.String("id", id),
			zap.String("url", originalURL),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}
	err = writer.WriteFile(row)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrConflictOriginalURL) || errors.Is(err, errorsInternal.ErrConflictID) {
			return id, err
		}

		logger.Log.Debug("could not write short url to file",
			zap.Error(err),
			zap.String("id", id),
			zap.String("url", originalURL),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}

	return id, nil
}

func (repo *JSONFileShortURLRepository) StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList, err error) {
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
	existingRows, err := reader.ReadFile()

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
			})
		}

		listShorten = append(listShorten, listShortenItem)
	}

	err = writer.WriteFileRows(rows)
	if err != nil {
		logger.Log.Debug("could not write short urls to file",
			zap.Error(err),
			zap.Any("rows", rows),
		)
		return nil, err
	}

	return listShorten, errors.Join(errs...)
}
