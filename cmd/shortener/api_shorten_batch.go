package main

import (
	"encoding/json"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (su *ShortUrlsService) apiShortenBatch(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	var list []jsonModel.BatchURLList

	err = json.Unmarshal(body, &list)
	if err != nil {
		logger.Log.Debug("could not parse request body",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	listShortened := su.repo.StoreBatch(list)
	if len(listShortened) == 0 {
		logger.Log.Debug("no urls were inserted",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	logger.Log.Info("short urls were created",
		zap.Any("list", listShortened),
	)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	// сериализуем ответ сервера
	enc := json.NewEncoder(res)
	if err := enc.Encode(listShortened); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
