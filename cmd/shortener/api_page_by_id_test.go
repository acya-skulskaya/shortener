package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_apiPageByID(t *testing.T) {
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
		{
			name:   "400 error when method (post) is wrong",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "400 error when method (put) is wrong",
			method: http.MethodPut,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "400 error when method (delete) is wrong",
			method: http.MethodDelete,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Bad Request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/ZYzivdwTSw", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			apiPageByID(w, request)

			res := w.Result()
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// проверяем Content-Type
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			// TODO проверяем Location по id
			//assert.Equal(t, test.want.headerLocation, res.Header.Get("Location"))
			// получаем и проверяем тело запроса

			//resBody, err := io.ReadAll(res.Body)
			//
			//require.NoError(t, err)
			//assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
