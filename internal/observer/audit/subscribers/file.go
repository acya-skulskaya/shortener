package subscribers

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/acya-skulskaya/shortener/internal/logger"
	model "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

type FileAuditSubscriber struct {
	name      string
	file      *os.File
	encoder   *json.Encoder
	eventChan chan model.AuditEvent
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closed    bool
}

func NewFileAuditSubscriber(ctx context.Context, filePath string) (*FileAuditSubscriber, error) {
	ctx, cancel := context.WithCancel(ctx)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Error("could not open file", zap.Error(err))
		return nil, err
	}
	encoder := json.NewEncoder(file)

	subscriber := &FileAuditSubscriber{
		file:      file,
		encoder:   encoder,
		name:      "FileAuditSubscriber",
		eventChan: make(chan model.AuditEvent, eventChanSize),
		ctx:       ctx,
		cancel:    cancel,
	}

	subscriber.wg.Add(1)
	go subscriber.worker()

	return subscriber, nil
}

func (s *FileAuditSubscriber) worker() {
	logger.Log.Debug("FileAuditSubscriber::worker")
	defer s.wg.Done()

	for {
		select {
		case event, ok := <-s.eventChan:
			if !ok {
				logger.Log.Debug("FileAuditSubscriber::worker eventChan was closed")
				s.closeFile()
				return
			}
			logger.Log.Debug("FileAuditSubscriber::worker writing new event to file", zap.Any("event", event))
			if err := s.encoder.Encode(event); err != nil {
				logger.Log.Error("FileAuditSubscriber::worker could not write event to audit file", zap.Any("event", event), zap.Error(err))
			}

		case <-s.ctx.Done():
			logger.Log.Debug("FileAuditSubscriber::worker ctx.Done()")
			s.mu.Lock()
			if !s.closed {
				close(s.eventChan)
				s.closed = true
			}
			s.mu.Unlock()
			s.closeFile()
			return
		}
	}
}

func (s *FileAuditSubscriber) ReceiveNewEvent(event model.AuditEvent) {
	select {
	case s.eventChan <- event:
		logger.Log.Debug("FileAuditSubscriber::ReceiveNewEvent new event was sent to file channel", zap.Any("event", event))
	case <-s.ctx.Done():
		logger.Log.Debug("FileAuditSubscriber::ReceiveNewEvent ctx.Done()")
	}
}

func (s *FileAuditSubscriber) GetName() string {
	return s.name
}

func (s *FileAuditSubscriber) closeFile() {
	logger.Log.Debug("FileAuditSubscriber::closeFile")
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.file != nil {
		s.file.Close()
		s.file = nil
	}
}

func (s *FileAuditSubscriber) Stop() {
	logger.Log.Debug("FileAuditSubscriber::stop")
	s.cancel()
	s.wg.Wait()
}
