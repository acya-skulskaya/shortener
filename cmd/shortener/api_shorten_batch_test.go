package main

import (
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_apiShortenBatch(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		body                string // добавляем тело запроса в табличные тесты
		expectedCode        int
		expectedContentType string
	}{
		{
			name:   "short urls stored",
			method: http.MethodPost,
			body: fmt.Sprintf(`[
    {
        "correlation_id": "%s",
        "original_url": "http://test.com"
    },
    {
        "correlation_id": "%s",
        "original_url": "http://test2.com"
    }
]`, helpers.RandStringRunes(10), helpers.RandStringRunes(10)),
			expectedCode:        http.StatusCreated,
			expectedContentType: "application/json",
		},
	}

	shortURLService := NewShortUrlsService(&shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: "./urls.json"})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(test.body)
			request := httptest.NewRequest(test.method, "/", bodyReader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			shortURLService.apiShortenBatch(w, request)

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
