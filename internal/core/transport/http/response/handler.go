package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	errorsgo "todoapp-try1/internal/core/errors.go"
	"todoapp-try1/internal/core/logger"

	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	rw  http.ResponseWriter
	log *logger.Logger
}

func NewHTTPResponseHandler(w http.ResponseWriter, l *logger.Logger) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		rw:  w,
		log: l,
	}
}

func (h *HTTPResponseHandler) JSONResponse(response any, statusCode int) {
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.ErrorResponse(
			"write http response",
			err,
		)
		return
	}
}

func (h *HTTPResponseHandler) ErrorResponse(msg string, err error) {
	var (
		statusCode int
		logFunc    func(msg string, fields ...zap.Field)
	)

	switch {
	case errors.Is(err, errorsgo.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	case errors.Is(err, errorsgo.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, errorsgo.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error(
			"write HTTP response",
			zap.String("message", msg),
			zap.Int("statusCode", http.StatusInternalServerError),
			zap.Error(err),
		)
		return
	}

	logFunc(msg, zap.Int("status_code", statusCode), zap.Error(err))
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.rw.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error(
			"write HTTP response",
			zap.String("message", msg),
			zap.Int("status_code", http.StatusInternalServerError),
			zap.Error(err))
		return
	}

	h.log.Debug(
		msg,
		zap.Int("status_code", http.StatusInternalServerError),
		zap.Error(err),
	)
}
