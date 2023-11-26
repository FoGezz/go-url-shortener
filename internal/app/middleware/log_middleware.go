package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// берём структуру для хранения сведений об ответе
type responseData struct {
	status int
	size   int
}

// добавляем реализацию http.ResponseWriter
type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData        *responseData
}

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

func ZapLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestLogger := zap.Must(zap.NewProduction()).Sugar()
		defer func() {
			_ = requestLogger.Sync()
		}()
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, req)
		requestLogger.Info(
			zap.String("Method", req.Method),
			zap.String("URI", req.RequestURI),
			zap.String("ExecTime", time.Since(start).String()),
			zap.Int("RespStatus", responseData.status),
			zap.Int("RespBodySize", responseData.size),
		)

	})
}
