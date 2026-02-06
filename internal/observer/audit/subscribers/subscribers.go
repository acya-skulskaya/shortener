package subscribers

import "github.com/acya-skulskaya/shortener/internal/model/json"

const eventChanSize = 50

type Subscriber interface {
	ReceiveNewEvent(event json.AuditEvent)
	GetName() string
	Stop()
}
