package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
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
		{
			name:   "400 error when method (get) is wrong",
			method: http.MethodGet,
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
			link := "https://practicum.yandex.ru/learn/go-advanced/courses/7154aca2-2665-440e-99ef-9dec1dfa1cd1/sprints/634244/topics/75da540c-e78d-4fdb-be66-c94ca0f88f58/lessons/6f432b47-f47c-4544-a686-7e2a94105cd6/"
			bodyReader := strings.NewReader(link)
			request := httptest.NewRequest(test.method, "/", bodyReader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			apiPageMain(w, request)

			res := w.Result()
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
				// TODO check id in db
			} else {
				assert.Equal(t, test.want.response, string(resBody))
			}
		})
	}
}

func TestRandStringRunes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "request length 1",
			args: args{
				n: 1,
			},
			want: 1,
		},
		{
			name: "request length 10",
			args: args{
				n: 10,
			},
			want: 10,
		},
		{
			name: "request length -1",
			args: args{
				n: -1,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := RandStringRunes(tt.args.n)
			assert.Equal(t, tt.want, len(str))
		})
	}
}
