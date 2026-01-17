package subscribers

import "github.com/acya-skulskaya/shortener/internal/model/json"

type Subscriber interface {
	ReceiveNewEvent(event json.AuditEvent)
	GetName() string
}
