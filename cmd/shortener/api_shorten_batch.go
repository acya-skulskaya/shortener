package main

import (
	"encoding/json"
	"errors"
	"net/http"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

// apiShorten handles the HTTP request to shorten a list of URLs in the request body in JSON format with IDs and return a list of shortened URLs with ID in JSON format
// Endpoint: POST /api/shorten/batch
// Expected request body: [{"correlation_id":"example123","original_url":"http://example.test/1"},{"correlation_id":"example456","original_url":"http://example.test/2"}]
// Returns:
//   - 201 Created when URLs were successfully shortened
//   - 401 Unauthorized if user is not authorized
//   - 409 Conflict when one of URLs in the request was already shortened
//   - 500 Internal Server Error on failure
func (su *ShortUrlsService) apiShortenBatch(res http.ResponseWriter, req *http.Request) {
	var list []jsonModel.BatchURLList
	if err := json.NewDecoder(req.Body).Decode(&list); err != nil {
		logger.Log.Debug("could not parse request body",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	ctx := req.Context()
	userID, ok := ctx.Value(middleware.AuthContextKey(middleware.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could nt get userID from context")
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	listShortened, err := su.Repo.StoreBatch(req.Context(), list, userID)
	if err != nil && (errors.Is(err, errorsInternal.ErrConflictOriginalURL) || errors.Is(err, errorsInternal.ErrConflictID)) {
		res.WriteHeader(http.StatusConflict)
	} else {
		logger.Log.Info("short urls were created",
			zap.Any("list", listShortened),
		)

		res.WriteHeader(http.StatusCreated)
	}

	// сериализуем ответ сервера
	enc := json.NewEncoder(res)
	if err := enc.Encode(listShortened); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
