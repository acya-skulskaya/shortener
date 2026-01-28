package subscribers

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/acya-skulskaya/shortener/internal/logger"
	model "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

type FileAuditSubscriber struct {
	FilePath string
	name     string
	mu       *sync.RWMutex
}

func NewFileAuditSubscriber(filePath string) *FileAuditSubscriber {
	subscriber := &FileAuditSubscriber{
		FilePath: filePath,
		name:     "FileAuditSubscriber",
		mu:       &sync.RWMutex{},
	}

	return subscriber
}

func (s *FileAuditSubscriber) ReceiveNewEvent(event model.AuditEvent) {
	logger.Log.Info("file subscriber received new event", zap.Any("event", event))

	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.OpenFile(s.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Error("could not open file to write event", zap.Any("event", event), zap.Error(err))
		return
	}
	defer file.Close()

	data, err := json.Marshal(&event)
	if err != nil {
		logger.Log.Error("could not marshall json", zap.Any("event", event), zap.Error(err))
		return
	}
	data = append(data, '\n')

	_, err = file.Write(data)
	if err != nil {
		logger.Log.Error("could not write to file", zap.Any("event", event), zap.Error(err))
		return
	}
}

func (s *FileAuditSubscriber) GetName() string {
	return s.name
}
