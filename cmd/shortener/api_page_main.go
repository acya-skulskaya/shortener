package main

import (
	"errors"
	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (su *ShortUrlsService) apiPageMain(res http.ResponseWriter, req *http.Request) {
	// На любой некорректный запрос сервер должен возвращать ответ с кодом 400.
	if req.Method != http.MethodPost {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	url := string(body)
	id, err := su.repo.Store(url)

	if len(id) == 0 {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err != nil && errors.Is(err, errorsInternal.ErrConflict) {
		res.WriteHeader(http.StatusConflict)
	} else {
		logger.Log.Info("short url was created",
			zap.String("id", id),
			zap.String("url", url),
		)
		res.WriteHeader(http.StatusCreated)
	}

	res.Write([]byte(config.Values.URLAddress + "/" + id))
}
