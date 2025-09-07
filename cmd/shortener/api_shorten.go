package main

import (
	"encoding/json"
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type RequestData struct {
	URL string `json:"url"`
}

type ResponseData struct {
	Result string `json:"result"`
}

func (su *ShortUrlsService) apiShorten(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	var requestData RequestData

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		logger.Log.Debug("could not parse request body", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url := requestData.URL

	id := su.repo.Store(url)

	logger.Log.Info("short url was created",
		zap.String("id", id),
		zap.String("url", url),
	)

	resp := ResponseData{Result: config.Values.URLAddress + "/" + id}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	// сериализуем ответ сервера
	enc := json.NewEncoder(res)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
