package main

import (
	"context"
	"encoding/json"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

func (su *ShortUrlsService) apiDeleteUserURLs(res http.ResponseWriter, req *http.Request) {
	var list []string

	if err := json.NewDecoder(req.Body).Decode(&list); err != nil {
		logger.Log.Debug("could not parse request body",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctx := req.Context()
	userID, ok := ctx.Value(middleware.AuthContextKey(middleware.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could not get userID from context")
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	go su.repo.DeleteUserUrls(context.Background(), list, userID)

	res.WriteHeader(http.StatusAccepted)
}
