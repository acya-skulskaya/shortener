package main

import (
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func (su *ShortUrlsService) apiPageByID(res http.ResponseWriter, req *http.Request) {
	// На любой некорректный запрос сервер должен возвращать ответ с кодом 400.
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")
	// TODO почему при запуске тестов епрестал чиатться параметр с помощью chi???
	if len(id) == 0 {
		id = req.PathValue("id")
	}

	url := su.repo.Get(id)

	if len(url) == 0 {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	logger.Log.Info("got page",
		zap.String("id", id),
		zap.String("url", url),
	)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
