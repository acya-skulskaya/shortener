package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
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
		panic(err)
	}

	// TODO save url and id in db
	id := RandStringRunes(10)
	ShortUrls[id] = string(body)

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
