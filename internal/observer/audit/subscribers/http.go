package subscribers

import (
	"context"
	"sync"
	"time"

	"github.com/acya-skulskaya/shortener/internal/logger"
	model "github.com/acya-skulskaya/shortener/internal/model/json"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type HTTPAuditSubscriber struct {
	URL       string
	name      string
	eventChan chan model.AuditEvent
	client    *resty.Client
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closed    bool
}

func NewHTTPAuditSubscriber(ctx context.Context, url string) *HTTPAuditSubscriber {
	ctx, cancel := context.WithCancel(ctx)
	client := resty.
		New().
		SetTimeout(3*time.Second).
		SetHeader("Content-Type", "application/json")

	subscriber := &HTTPAuditSubscriber{
		URL:       url,
		name:      "HTTPAuditSubscriber",
		eventChan: make(chan model.AuditEvent, eventChanSize),
		client:    client,
		ctx:       ctx,
		cancel:    cancel,
	}

	subscriber.wg.Add(1)
	go subscriber.worker()

	return subscriber
}

func (s *HTTPAuditSubscriber) ReceiveNewEvent(event model.AuditEvent) {
	if s.closed {
		logger.Log.Debug("HTTPAuditSubscriber::ReceiveNewEvent can not receive new event because channel is closed", zap.Any("event", event))
		return
	}

	select {
	case s.eventChan <- event:
		logger.Log.Debug("HTTPAuditSubscriber::ReceiveNewEvent new event was sent to file channel", zap.Any("event", event))
	case <-s.ctx.Done():
		logger.Log.Debug("HTTPAuditSubscriber::ReceiveNewEvent ctx.Done()")
		if !s.closed {
			close(s.eventChan)
			s.closed = true
		}
	}
}

func (s *HTTPAuditSubscriber) GetName() string {
	return s.name
}

func (s *HTTPAuditSubscriber) worker() {
	logger.Log.Debug("FileAuditSubscriber::worker")
	defer s.wg.Done()

	for {
		select {
		case event, ok := <-s.eventChan:
			if !ok {
				logger.Log.Debug("HTTPAuditSubscriber::worker eventChan was closed")
				return
			}
			go func(event model.AuditEvent) {
				resp, err := s.client.R().SetBody(event).Post(s.URL)
				if err != nil {
					logger.Log.Error("HTTPAuditSubscriber::worker could not post audit event", zap.Any("event", event), zap.Error(err))
					return
				}
				if resp.StatusCode() >= 400 {
					logger.Log.Warn("HTTPAuditSubscriber::worker received error status trying to post audit event", zap.Any("event", event), zap.Any("Status", resp.Status()))
					return
				}
				logger.Log.Debug("HTTPAuditSubscriber::worker audit event was sent", zap.Any("event", event))
			}(event)

		case <-s.ctx.Done():
			logger.Log.Debug("HTTPAuditSubscriber::worker ctx.Done()")
		}
	}
}

func (s *HTTPAuditSubscriber) Stop() {
	logger.Log.Debug("HTTPAuditSubscriber::stop")
	s.cancel()
	s.wg.Wait()
}
