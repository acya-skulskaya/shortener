package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_apiPageMain(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "short url created",
			method: http.MethodPost,
			want: want{
				code:        http.StatusCreated,
				response:    "url",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	os.Remove("./urls.json")

	auditPublisher := publisher.NewAuditPublisher()
	shortURLService := NewShortUrlsService(&shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: "./urls.json"}, auditPublisher)

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			link := "https://practicum.yandex.ru/learn/go-advanced/courses/7154aca2-2665-440e-99ef-9dec1dfa1cd1/sprints/634244/topics/75da540c-e78d-4fdb-be66-c94ca0f88f58/lessons/6f432b47-f47c-4544-a686-7e2a94105cd6/"
			bodyReader := strings.NewReader(link)
			request, err := http.NewRequest(test.method, testServer.URL+"/", bodyReader)
			require.NoError(t, err)

			client := &http.Client{}

			res, err := client.Do(request)
			require.NoError(t, err)
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// проверяем Content-Type
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			if test.want.response == "url" {
				_, err := url.ParseRequestURI(string(resBody))
				assert.NoError(t, err)
			} else {
				assert.Equal(t, test.want.response, string(resBody))
			}
		})
	}
}
