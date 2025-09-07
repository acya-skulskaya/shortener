package model

type URLList struct {
	ID          string `json:"id"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
