package main

import (
	"encoding/json"
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

/*
Инкремент 7
Задание по треку «Сервис сокращения URL»
Добавьте в код сервера новый эндпоинт POST /api/shorten, который будет принимать в теле запроса JSON-объект {"url":"<some_url>"} и возвращать в ответ объект {"result":"<short_url>"}.
Запрос может иметь такой вид:

POST http://localhost:8080/api/shorten HTTP/1.1
Host: localhost:8080
Content-Type: application/json
{
  "url": "https://practicum.yandex.ru"
}

Ответ может быть таким:

HTTP/1.1 201 OK
Content-Type: application/json
Content-Length: 30
{
 "result": "http://localhost:8080/EwHXdJfB"
}

Удостоверьтесь, что в ответе от сервера присутствует HTTP-заголовок Content-Type со значением application/json. Он указывает клиенту, в каком формате передано тело ответа.
Также не забудьте добавить тесты на новый эндпоинт, как и на предыдущие.
При реализации задействуйте одну из распространённых библиотек:

    encoding/json,
    github.com/mailru/easyjson,
    github.com/pquerna/ffjson,
    github.com/labstack/echo,
    github.com/goccy/go-json.
*/

type RequestData struct {
	URL string `json:"url"`
}

type ResponseData struct {
	Result string `json:"result"`
}

func apiShorten(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	var requestData RequestData

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		logger.Log.Debug("could not parse request body", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url := requestData.URL

	// TODO save url and id in db
	// TODO make repository
	id := helpers.RandStringRunes(10)
	Cont.add(id, url)

	logger.Log.Info("short url was created",
		zap.String("id", id),
		zap.String("url", url),
	)

	resp := ResponseData{Result: config.Values.URLAddress + "/" + id}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	// сериализуем ответ сервера
	enc := json.NewEncoder(res)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
