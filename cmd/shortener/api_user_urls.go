package main

import (
	"encoding/json"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

func (su *ShortUrlsService) apiUserURLs(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(middleware.AuthContextKey(middleware.AuthContextKeyUserID)).(string)

	list, err := su.repo.GetUserUrls(req.Context(), userID)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(res)
	if err := enc.Encode(list); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
