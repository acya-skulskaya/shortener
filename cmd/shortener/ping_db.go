package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// apiPingDB handles the HTTP request to test connection to DB
// Endpoint: GET /ping
// Returns:
//   - 200 OK if a connection was successful
//   - 500 Internal Server Error on failure
func (su *ShortUrlsService) apiPingDB(res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("pgx", config.Values.DatabaseDSN)
	if err != nil {
		logger.Log.Debug("could not open db connection",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		logger.Log.Debug("could not ping db",
			zap.Error(err),
		)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
