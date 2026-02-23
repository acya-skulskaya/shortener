package http

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	models "github.com/acya-skulskaya/shortener/internal/model/json"
	authService "github.com/acya-skulskaya/shortener/internal/service/auth"
	"go.uber.org/zap"
)

// apiPageMain handles the HTTP request to shorten the URL in the request body and return a URL with ID
// Endpoint: POST /
// Expected request body: plain text with the original URL
// Returns:
//   - 201 Created when ID was successfully created
//   - 401 Unauthorized if user is not authorized
//   - 409 Conflict when the URL in the request was already shortened
//   - 500 Internal Server Error on failure
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
	ctx := req.Context()
	userID, ok := ctx.Value(authService.AuthContextKey(authService.AuthContextKeyUserID)).(string)
	if !ok {
		logger.Log.Debug("could nt get userID from context")
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, err := su.Repo.Store(ctx, url, userID)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err != nil {
		if errors.Is(err, errorsInternal.ErrConflictOriginalURL) {
			res.WriteHeader(http.StatusConflict)
		} else {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		logger.Log.Info("short url was created",
			zap.String("id", id),
			zap.String("url", url),
		)
		res.WriteHeader(http.StatusCreated)
	}

	su.auditPublisher.Notify(models.AuditEvent{
		Timestamp:   time.Now().Unix(),
		Action:      models.AuditEventActionTypeShorten,
		UserID:      userID,
		OriginalURL: url,
	})

	res.Write([]byte(config.Values.URLAddress + "/" + id))
}
