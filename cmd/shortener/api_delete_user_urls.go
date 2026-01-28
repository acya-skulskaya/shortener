package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"go.uber.org/zap"
)

// apiDeleteUserURLs handles the HTTP request to delete a list of IDS of shortened URLs created by the logged in user.
// Endpoint: DELETE /api/user/urls
// Expected request body: ["ExampleID1", "ExampleID1"]
// Returns:
//   - 202 Accepted on success
//   - 401 Unauthorized if user is not authorized
//   - 500 Internal Server Error on failure
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go su.Repo.DeleteUserUrls(ctx, list, userID)

	res.WriteHeader(http.StatusAccepted)
}
