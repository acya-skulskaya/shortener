package http

import (
	"errors"
	"net/http"
	"time"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	models "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// apiPageByID handles the HTTP request to follow the original URL from the ID
// Endpoint: GET /{id}
// Returns:
//   - 307 Temporary Redirect when ID is found
//   - 401 Unauthorized if user is not authorized
//   - 404 Not Found when ID is not found
//   - 410 Gone  when ID is deleted
//   - 500 Internal Server Error on failure
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

	su.auditPublisher.Notify(models.AuditEvent{
		Timestamp:   time.Now().Unix(),
		Action:      models.AuditEventActionTypeFollow,
		OriginalURL: url,
	})

	logger.Log.Info("got page",
		zap.String("id", id),
		zap.String("url", url),
	)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
