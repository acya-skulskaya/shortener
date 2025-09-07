package shortUrlJsonFile

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/model"
	"go.uber.org/zap"
)

type JSONFileShortURLRepository struct {
}

func (repo *JSONFileShortURLRepository) Get(id string) (originalURL string) {
	reader, _ := NewFileReader(config.Values.FileStoragePath)
	defer reader.Close()
	list, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("id", id),
			zap.String("file", config.Values.FileStoragePath),
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

// func (repo *InJSONFileShortURLRepository) Store(originalURL string) (id string, err error) {
func (repo *JSONFileShortURLRepository) Store(originalURL string) (id string) {
	id = helpers.RandStringRunes(10)

	row := model.URLList{
		ID:          id,
		ShortURL:    config.Values.URLAddress + "/" + id,
		OriginalURL: originalURL,
	}

	writer, _ := NewFileWriter(config.Values.FileStoragePath)
	err := writer.WriteFile(row)
	if err != nil {
		logger.Log.Debug("could not write short url to file",
			zap.Error(err),
			zap.String("id", id),
			zap.String("url", originalURL),
			zap.String("file", config.Values.FileStoragePath),
		)
		return ""
	}

	return id
}
