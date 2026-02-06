package publisher

import (
	"github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/subscribers"
)

type Publisher interface {
	Subscribe(subscriber subscribers.Subscriber)
	Unsubscribe(subscriber subscribers.Subscriber)
	Notify(event json.AuditEvent)
	Shutdown()
}

type AuditPublisher struct {
	Subscribers map[string]subscribers.Subscriber
	eventChan   chan json.AuditEvent
}

func (e *AuditPublisher) Shutdown() {
	for _, subscriber := range e.Subscribers {
		subscriber.Stop()
	}
}

func NewAuditPublisher() *AuditPublisher {
	publisher := &AuditPublisher{}

	return publisher
}

func (e *AuditPublisher) Subscribe(subscriber subscribers.Subscriber) {
	if e.Subscribers == nil {
		e.Subscribers = make(map[string]subscribers.Subscriber)
	}
	e.Subscribers[subscriber.GetName()] = subscriber
}

func (e *AuditPublisher) Unsubscribe(subscriber subscribers.Subscriber) {
	delete(e.Subscribers, subscriber.GetName())
}

func (e *AuditPublisher) Notify(event json.AuditEvent) {
	for _, subscriber := range e.Subscribers {
		subscriber.ReceiveNewEvent(event)
	}
}
