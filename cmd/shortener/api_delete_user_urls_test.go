package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_apiDeleteUserURLs(t *testing.T) {
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
			name:   "endpoint returns accepted response",
			method: http.MethodDelete,
			want: want{
				code: http.StatusAccepted,
			},
		},
	}

	os.Remove("./urls.json")

	auditPublisher := publisher.NewAuditPublisher()
	repo := &shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: "./urls.json"}
	shortURLService := NewShortUrlsService(repo, auditPublisher)

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyReader := strings.NewReader(`[{"correlation_id": "tgkPCfkqXF","original_url": "http://test.com"},{"correlation_id": "tgkPCfkqXg","original_url": "http://test2.com"}]`)
			request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/shorten/batch", bodyReader)
			require.NoError(t, err)
			res, err := testServer.Client().Do(request)
			require.NoError(t, err)
			res.Body.Close()
			cookies := res.Cookies()

			bodyReader = strings.NewReader(`["tgkPCfkqXF"]`)
			request, err = http.NewRequest(test.method, testServer.URL+"/api/user/urls", bodyReader)
			require.NoError(t, err)
			for _, cookie := range cookies {
				request.AddCookie(cookie)
			}

			response, err := testServer.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.want.code, response.StatusCode)
		})
	}
}
