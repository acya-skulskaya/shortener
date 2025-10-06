package main

import (
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_apiShorten(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		body                string // добавляем тело запроса в табличные тесты
		expectedCode        int
		expectedContentType string
	}{
		{
			name:                "short url created",
			method:              http.MethodPost,
			body:                `{"url": "http://test.test"}`,
			expectedCode:        http.StatusCreated,
			expectedContentType: "application/json",
		},
		{
			name:                "500 error on invalid json",
			method:              http.MethodPost,
			body:                `{"url": "http://test.test}`,
			expectedCode:        http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
		},
	}

	os.Remove("./urls.json")

	shortURLService := NewShortUrlsService(&shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: "./urls.json"})

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.body)
			request, err := http.NewRequest(test.method, testServer.URL+"/api/shorten", bodyReader)
			require.NoError(t, err)

			client := &http.Client{}

			res, err := client.Do(request)
			require.NoError(t, err)
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.expectedCode, res.StatusCode)
			// проверяем Content-Type
			assert.Equal(t, test.expectedContentType, res.Header.Get("Content-Type"))
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			//resBody, err := io.ReadAll(res.Body)
			_, err = io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
