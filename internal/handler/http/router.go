package http

import (
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// NewRouter initiates a new router with API's endpoints:
//   - POST / - creates a shortened URL
//   - GET /{id} — redirects to the original URL
//   - GET /ping — tests connection to DB
//   - POST /api/shorten - creates a shortened URL
//   - POST /api/shorten/batch - creates a  list of shortened URLs
//   - GET /api/user/urls - returns a list of URLs that were added by an authenticated user
//   - DELETE /api/user/urls - deletes a list of shortened URLs
func NewRouter(su *ShortUrlsService) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestLogger)
	router.Use(middleware.RequestCompressor)
	router.Use(middleware.CookieAuth)

	// pprof
	//router.Route("/debug/pprof", func(r chi.Router) {
	//	r.Handle("/", http.HandlerFunc(pprof.Index))
	//	r.Handle("/profile", http.HandlerFunc(pprof.Profile))
	//	r.Handle("/symbol", http.HandlerFunc(pprof.Symbol))
	//	r.Handle("/cmdline", http.HandlerFunc(pprof.Cmdline))
	//	r.Handle("/heap", pprof.Handler("heap"))
	//})

	router.Post("/", su.apiPageMain)
	router.Get("/{id}", su.apiPageByID)
	router.Get("/ping", su.apiPingDB)
	router.Post("/api/shorten", su.apiShorten)
	router.Post("/api/shorten/batch", su.apiShortenBatch)
	router.Get("/api/user/urls", su.apiUserURLs)
	router.Delete("/api/user/urls", su.apiDeleteUserURLs)
	router.Get("/api/internal/stats", su.apiInternalStats)

	return router
}
