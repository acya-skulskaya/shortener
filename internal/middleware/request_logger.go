package middleware

import (
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// RequestLogger HTTP middleware setting a value on the request context
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()
		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start).String()

		// отправляем сведения о запросе в zap

		logger.Log.Info("REQUEST",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.String("duration", duration),
		)

		logger.Log.Info("REQUEST",
			zap.Int("status", responseData.status), // получаем перехваченный код статуса ответа
			zap.Int("size", responseData.size),     // получаем перехваченный размер ответа
		)
	})
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
