package subscribers

import (
	"time"

	"github.com/acya-skulskaya/shortener/internal/logger"
	model "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const eventChanSize = 20

type HTTPAuditSubscriber struct {
	URL       string
	name      string
	eventChan chan model.AuditEvent
	client    *resty.Client
}

func NewHTTPAuditSubscriber(url string) *HTTPAuditSubscriber {
	client := resty.New().SetTimeout(3*time.Second).SetHeader("Content-Type", "application/json")

	subscriber := &HTTPAuditSubscriber{
		URL:       url,
		name:      "HTTPAuditSubscriber",
		eventChan: make(chan model.AuditEvent, eventChanSize),
		client:    client,
	}

	go subscriber.worker()

	return subscriber
}

func (s *HTTPAuditSubscriber) ReceiveNewEvent(event model.AuditEvent) {
	logger.Log.Info("http subscriber received new event", zap.Any("event", event))
	s.eventChan <- event
}

func (s *HTTPAuditSubscriber) GetName() string {
	return s.name
}

func (s *HTTPAuditSubscriber) worker() {
	for event := range s.eventChan {
		func(event model.AuditEvent) {
			resp, err := s.client.R().SetBody(event).Post(s.URL)
			if err != nil {
				logger.Log.Error("could not post audit event", zap.Any("event", event), zap.Error(err))
			}
			if resp.StatusCode() >= 400 {
				logger.Log.Debug("received error status trying to post audit event", zap.Any("event", event), zap.Any("Status", resp.Status()))
			}
		}(event)
	}
}
