package shorturljsonfile

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

type JSONFileShortURLRepository struct {
	FileStoragePath string
}

func (repo *JSONFileShortURLRepository) Get(id string) (originalURL string) {
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

	list, _, err := reader.ReadFile()
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

func (repo *JSONFileShortURLRepository) Store(originalURL string) (id string, err error) {
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

func (repo *JSONFileShortURLRepository) StoreBatch(listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList) {
	writer, err := NewFileWriter(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create file writer",
			zap.Error(err),
		)
		return nil
	}

	var rows []jsonModel.URLList

	for _, item := range listOriginal {
		// TODO check if id already exists
		rows = append(rows, jsonModel.URLList{
			ID:          item.CorrelationID,
			OriginalURL: item.OriginalURL,
			ShortURL:    config.Values.URLAddress + "/" + item.CorrelationID,
		})

		listShorten = append(listShorten, jsonModel.BatchURLList{
			CorrelationID: item.CorrelationID,
			ShortURL:      config.Values.URLAddress + "/" + item.CorrelationID,
		})
	}

	err = writer.WriteFileRows(rows)
	if err != nil {
		logger.Log.Debug("could not write short urls to file",
			zap.Error(err),
			zap.Any("rows", rows),
		)
		return nil
	}

	return listShorten
}
