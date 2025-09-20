package json

type BatchURLList struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}
