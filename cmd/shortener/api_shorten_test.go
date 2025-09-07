package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
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

	shortURLService := NewShortUrlsService(&shortUrlJsonFile.JSONFileShortURLRepository{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.body)
			request := httptest.NewRequest(test.method, "/", bodyReader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			shortURLService.apiShorten(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.expectedCode, res.StatusCode)
			// проверяем Content-Type
			assert.Equal(t, test.expectedContentType, res.Header.Get("Content-Type"))
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			//resBody, err := io.ReadAll(res.Body)
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			// TODO check id in db

		})
	}
}
