package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortUrlsService_apiInternalStats(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name          string
		trustedSubnet string
		headerIP      string
		want          want
	}{
		{
			name:          "endpoint returns stats",
			trustedSubnet: "192.168.1.0/24",
			headerIP:      "192.168.1.1",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          "user IPis not in trusted subnet",
			trustedSubnet: "192.168.1.0/24",
			headerIP:      "92.68.1.1",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:          "trusted subnet is not set",
			trustedSubnet: "",
			headerIP:      "92.68.1.1",
			want: want{
				code: http.StatusUnauthorized,
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
			config.Values.TrustedSubnet = test.trustedSubnet

			request, err := http.NewRequest(http.MethodGet, testServer.URL+"/api/internal/stats", nil)
			require.NoError(t, err)
			request.Header.Add("X-Real-IP", test.headerIP)

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
		})
	}
}
