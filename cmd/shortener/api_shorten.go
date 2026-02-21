package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	models "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

// apiShorten handles the HTTP request to shorten the URL in the request body in JSON format and return a URL with ID in JSON format
// Endpoint: POST /api/shorten
// Expected request body: {"url":"http://example.test"}
// Returns:
//   - 201 Created when ID was successfully created
//   - 401 Unauthorized if user is not authorized
//   - 409 Conflict when the URL in the request was already shortened
//   - 500 Internal Server Error on failure
func (su *ShortUrlsService) apiShorten(res http.ResponseWriter, req *http.Request) {
	type RequestData struct {
		URL string `json:"url"`
	}

	type ResponseData struct {
		Result string `json:"result"`
	}

	var requestData RequestData
	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		logger.Log.Debug("could not parse request body",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url := requestData.URL
	ctx := req.Context()
	userID, ok := ctx.Value(middleware.AuthContextKey(middleware.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could nt get userID from context")
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, err := su.Repo.Store(req.Context(), url, userID)

	if len(id) == 0 {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err != nil && errors.Is(err, errorsInternal.ErrConflictOriginalURL) {
		res.WriteHeader(http.StatusConflict)
	} else {
		logger.Log.Info("short url was created",
			zap.String("id", id),
			zap.String("url", url),
		)

		res.WriteHeader(http.StatusCreated)
	}

	resp := ResponseData{Result: config.Values.URLAddress + "/" + id}

	su.auditPublisher.Notify(models.AuditEvent{
		Timestamp:   time.Now().Unix(),
		Action:      models.AuditEventActionTypeShorten,
		UserID:      userID,
		OriginalURL: url,
	})

	// сериализуем ответ сервера
	enc := json.NewEncoder(res)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
