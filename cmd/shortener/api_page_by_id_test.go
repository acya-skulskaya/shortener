package main

import (
	"context"
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_apiPageByID(t *testing.T) {
	config.Init()

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
			name:   "short url redirects",
			method: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "Temporary Redirect\n",
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	os.Remove("./urls.json")

	auditPublisher := publisher.NewAuditPublisher()
	repo := &shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: "./urls.json"}
	shortURLService := NewShortUrlsService(repo, auditPublisher)
	id, _ := repo.Store(context.Background(), "https://test.com", "userID123")

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(test.method, testServer.URL+"/"+id, nil)
			require.NoError(t, err)

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			response, err := client.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.want.code, response.StatusCode)
			// проверяем Content-Type
			assert.Equal(t, test.want.contentType, response.Header.Get("Content-Type"))
		})
	}
}
