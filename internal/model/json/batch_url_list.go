package json

type BatchURLList struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	Err           string `json:"err,omitempty"`
	IsDeleted     int    `json:"is_deleted,omitempty"`
}
