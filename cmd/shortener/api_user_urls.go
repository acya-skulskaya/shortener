package main

import (
	"encoding/json"
	"net/http"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"go.uber.org/zap"
)

// apiUserURLs handles the HTTP request to get a list of short URLs, that were added by an authenticated user
// Endpoint: GET /api/user/urls
// Returns:
//   - 200 OK returns a list of short URLs
//   - 204 No Content if an authenticated user has no added short URLs
//   - 401 Unauthorized if user is not authorized
//   - 500 Internal Server Error on failure
func (su *ShortUrlsService) apiUserURLs(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID, ok := ctx.Value(middleware.AuthContextKey(middleware.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could nt get userID from context")
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	list, err := su.Repo.GetUserUrls(req.Context(), userID)
	if err != nil {
		logger.Log.Debug("error getting user urls", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(res)
	if err := enc.Encode(list); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
