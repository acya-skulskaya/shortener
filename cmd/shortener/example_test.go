package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
)

// This example demonstrates a successful request to shorten a URL via endpoint POST /
// Endpoint POST / returns a shortened URL
func Example_apiPageMain_Test() {
	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, publisher.NewAuditPublisher())

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Making a request to POST / endpoint with a URL that needs to be shortened in the request body
	link := "https://practicum.yandex.ru/learn/go-advanced/courses/7154aca2-2665-440e-99ef-9dec1dfa1cd1/sprints/634244/topics/75da540c-e78d-4fdb-be66-c94ca0f88f58/lessons/6f432b47-f47c-4544-a686-7e2a94105cd6/"
	bodyReader := strings.NewReader(link)
	request, _ := http.NewRequest(http.MethodPost, testServer.URL+"/", bodyReader)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	// Result:
	// Status: 201 Created
	// Response Body:
	// http://localhost:8080/CqCvChqHPd
}

// This example demonstrates a successful request to delete a shortened URL via endpoint DELETE /api/user/urls
// Endpoint DELETE /api/user/urls returns status 202 if the request is accepted
func Example_apiDeleteUserURLs_Test() {
	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, publisher.NewAuditPublisher())

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Making a request to DELETE /api/user/urls endpoint with a list of shortened URL IDs
	bodyReader := strings.NewReader(`["tgkPCfkqXF"]`)
	request, _ := http.NewRequest(http.MethodDelete, testServer.URL+"/api/user/urls", bodyReader)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	// Result:
	// Status: 202 Accepted
	// Response Body: empty
}

// This example demonstrates a successful request to follow a shortened URL via endpoint GET /{id}
// Endpoint GET /{id} returns status 307 and redirects to the original URL
func Example_apiPageByID_Test() {
	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, publisher.NewAuditPublisher())

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Making a request to GET /{id} endpoint with a list of shortened URL IDs
	request, _ := http.NewRequest(http.MethodGet, testServer.URL+"/tgkPCfkqXF", http.NoBody)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	// Result:
	// Status: 307 Temporary Redirect
	// Location: https://original-url-domain.test/some/path
}

// This example demonstrates a successful request to shorten a URL via endpoint POST /api/shorten
// Endpoint POST /api/shorten returns status 201 if a shortened URL was created
func Example_apiShorten_Test() {
	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, publisher.NewAuditPublisher())

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Making a request to POST /api/shorten endpoint with a URL that needs to be shortened in the request body
	bodyReader := strings.NewReader(`{"url": "http://test.test"}`)
	request, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/shorten", bodyReader)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	// Result:
	// Status: 201 Created
	// Response Body:
	// {
	// 	"result": "http://localhost:8080/jOIIqnBYWi"
	// }
}

// This example demonstrates a successful request to shorten a list of URLs via endpoint POST /api/shorten/batch
// Endpoint POST /api/shorten/batch returns status 201 if a shortened URL was created
func Example_apiShortenBatch_Test() {
	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, publisher.NewAuditPublisher())

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Making a request to POST /api/shorten/batch endpoint with a list of URLs to be shortened in the request body
	bodyReader := strings.NewReader(`[{"correlation_id":"example123","original_url":"http://example.test/1"},{"correlation_id":"example456","original_url":"http://example.test/2"}]`)
	request, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/shorten/batch", bodyReader)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	// Result:
	// Status: 201 Created
	// Response Body:
	// [
	//	{
	//		"correlation_id": "example123",
	//		"short_url": "http://localhost:8080/tgkPCfkqXF"
	//	},
	//	{
	//		"correlation_id": "example456",
	//		"short_url": "http://localhost:8080/AgkPCfkRxf"
	//	}
	// ]
}

// This example demonstrates a successful request to get a list of URLs that were added by an authenticated user via endpoint GET /api/user/urls
// Endpoint GET /api/user/urls returns status 200 and a list of URLs
func Example_apiUserURLs_Test() {
	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, publisher.NewAuditPublisher())

	router := NewRouter(shortURLService)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Making a request to GET /api/user/urls endpoint with a list of URLs to be shortened in the request body
	request, _ := http.NewRequest(http.MethodGet, testServer.URL+"/api/user/urls", http.NoBody)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	// Result:
	// Status: 20o OK
	// Response Body:
	// [
	//	{
	//		"short_url": "http://localhost:8080/CqCvChqHPd",
	//		"original_url": "http://example.test/1"
	//	},
	//	{
	//		"short_url": "http://localhost:8080/CACvChqHPf",
	//		"original_url": "http://example.test/2"
	//	},
	//	{
	//		"short_url": "http://localhost:8080/CqAvChqgPd",
	//		"original_url": "http://example.test/3"
	//	},
	//	{
	//		"short_url": "http://localhost:8080/dqsvChqHPd",
	//		"original_url": "http://example.test/4"
	//	}
	//]
}
