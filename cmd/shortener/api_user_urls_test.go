package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_apiUserURLs(t *testing.T) {
	type want struct {
		code        int
		numUrls     int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		method  string
		numUrls int
		want    want
	}{
		{
			name:    "endpoint returns urls",
			method:  http.MethodGet,
			numUrls: 2,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				numUrls:     2,
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
			var res *http.Response
			var cookies []*http.Cookie
			for i := 0; i < test.numUrls; i++ {
				link := "https://practicum.yandex.ru/learn/go-advanced/courses/7154aca2-2665-440e-99ef-9dec1dfa1cd1/sprints/634244/topics/75da540c-e78d-4fdb-be66-c94ca0f88f58/lessons/" + uuid.New().String()
				bodyReader := strings.NewReader(link)
				request, err := http.NewRequest(http.MethodPost, testServer.URL+"/", bodyReader)
				require.NoError(t, err)
				if i > 0 {
					for _, cookie := range cookies {
						request.AddCookie(cookie)
					}
				}

				res, err = testServer.Client().Do(request)
				require.NoError(t, err)
				res.Body.Close()

				if i == 0 {
					cookies = res.Cookies()
				}
			}

			request, err := http.NewRequest(test.method, testServer.URL+"/api/user/urls", nil)
			require.NoError(t, err)
			for _, cookie := range cookies {
				request.AddCookie(cookie)
			}

			response, err := testServer.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.want.code, response.StatusCode)
			// проверяем Content-Type
			assert.Equal(t, test.want.contentType, response.Header.Get("Content-Type"))

			var list []jsonModel.BatchURLList

			err = json.NewDecoder(response.Body).Decode(&list)
			require.NoError(t, err)
			assert.Equal(t, test.want.numUrls, len(list))
		})
	}
}
