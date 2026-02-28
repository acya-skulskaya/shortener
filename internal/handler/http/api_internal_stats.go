package http

import (
	"encoding/json"
	"net/http"

	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	ip2 "github.com/vikram1565/request-ip"
	"go.uber.org/zap"
)

// apiInternalStats handles the HTTP request to show statistics of shortened URLs and users stored in service
// Endpoint: GET /api/internal/stats
// Returns:
//   - 200 OK returns statistics of shortened URLs and users stored in service
//   - 401 Unauthorized if user's ip is not in trusted subnet or subnet is not configured
//   - 500 Internal Server Error on failure
func (su *ShortUrlsService) apiInternalStats(res http.ResponseWriter, req *http.Request) {
	type ResponseData struct {
		URLs  int `json:"urls"`
		Users int `json:"users"`
	}

	res.Header().Set("Content-Type", "application/json")

	if config.Values.TrustedSubnet == "" {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	ip := ip2.GetClientIP(req)
	checkIPResult, err := helpers.CheckIPSubnet(ip, config.Values.TrustedSubnet)
	if err != nil {
		logger.Log.Debug("could not check ip subnet",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !checkIPResult {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	urls, users, err := su.repo.GetInternalStats(req.Context())
	if err != nil {
		logger.Log.Debug("could not get internal stats", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := ResponseData{
		URLs:  urls,
		Users: users,
	}

	res.WriteHeader(http.StatusOK)

	// сериализуем ответ сервера
	enc := json.NewEncoder(res)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
