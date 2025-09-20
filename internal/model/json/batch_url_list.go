package json

type BatchURLList struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}
