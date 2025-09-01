package main

import (
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

/*
Эндпоинт с методом GET и путём /{id}, где id — идентификатор сокращённого URL (например, /EwHXdJfB). В случае успешной обработки запроса сервер возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.

Пример запроса к серверу:
GET /EwHXdJfB HTTP/1.1
Host: localhost:8080
Content-Type: text/plain

Пример ответа от сервера:
HTTP/1.1 307 Temporary Redirect
Location: https://practicum.yandex.ru/
*/
func apiPageByID(res http.ResponseWriter, req *http.Request) {
	// На любой некорректный запрос сервер должен возвращать ответ с кодом 400.
	if req.Method != http.MethodGet {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(req, "id")
	// TODO get url from db
	url := Cont.getUrl(id)
	//url := ShortUrls[id]
	//res.WriteHeader(http.StatusTemporaryRedirect)
	//res.Header().Set("Location", url)
	//	res.Header().Add("Location", url)

	logger.Log.Info("got page",
		zap.String("id", id),
		zap.String("url", url),
	)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
