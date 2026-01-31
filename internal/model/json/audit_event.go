package json

type AuditEventActionType string

const (
	AuditEventActionTypeShorten AuditEventActionType = "shorten"
	AuditEventActionTypeFollow  AuditEventActionType = "follow"
)

// generate:reset
type AuditEvent struct {
	Timestamp   int64                `json:"ts"`                // unix timestamp события
	Action      AuditEventActionType `json:"action"`            // действие: shorten (создание) или follow (прохождение по ссылке)
	UserID      string               `json:"user_id,omitempty"` // идентификатор пользователя, если есть
	OriginalURL string               `json:"url"`               // оригинальный (не сокращенный) URL
}
