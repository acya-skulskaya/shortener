package main

import (
	"encoding/json"
	"errors"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (su *ShortUrlsService) apiDeleteUserURLs(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	var list []string

	err = json.Unmarshal(body, &list)
	if err != nil {
		logger.Log.Debug("could not parse request body",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctx := req.Context()
	userID := ctx.Value(middleware.AuthContextKey(middleware.AuthContextKeyUserID)).(string)
	res.Header().Set("Content-Type", "application/json")

	err = su.repo.DeleteUserUrls(req.Context(), list, userID)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrIDNotFound) {
			logger.Log.Debug("one of ids does not exist",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, errorsInternal.ErrIDDeleted) {
			logger.Log.Debug("one of ids is already deleted",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusGone), http.StatusGone)
			return
		} else if errors.Is(err, errorsInternal.ErrUserIDUnauthorized) {
			logger.Log.Debug("one of ids does not belong to current user",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else {
			logger.Log.Debug("could not delete urls",
				zap.Error(err),
			)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	res.WriteHeader(http.StatusAccepted)
}
