package main

import "C"
import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
	"io"
	"math"
	"math/rand"
	"net/http"
)

/*
Эндпоинт с методом POST и путём /. Сервер принимает в теле запроса строку URL как text/plain и возвращает ответ с кодом 201 и сокращённым URL как text/plain.

Пример запроса к серверу:
POST / HTTP/1.1
Host: localhost:8080
Content-Type: text/plain

https://practicum.yandex.ru/

Пример ответа от сервера:
HTTP/1.1 201 Created
Content-Type: text/plain
Content-Length: 30

http://localhost:8080/EwHXdJfB
*/
func apiPageMain(res http.ResponseWriter, req *http.Request) {
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

	// TODO save url and id in db
	id := RandStringRunes(10)
	url := string(body)
	Cont.add(id, url)
	//ShortUrls[id] = string(body)

	logger.Log.Info("short url was created",
		zap.String("id", id),
		zap.String("url", url),
	)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.Values.URLAddress + "/" + id))
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	if n < 0 {
		n = int(math.Abs(float64(n)))
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
