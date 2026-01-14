package main

import (
	"errors"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
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
	url, err := su.repo.Get(req.Context(), id)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrIDNotFound) {
			logger.Log.Debug("id does not exist",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if errors.Is(err, errorsInternal.ErrIDDeleted) {
			logger.Log.Debug("id is deleted",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusGone), http.StatusGone)
			return
		} else {
			logger.Log.Debug("could not get id",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	logger.Log.Info("got page",
		zap.String("id", id),
		zap.String("url", url),
	)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
