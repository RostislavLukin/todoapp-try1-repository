package response

import "net/http"

var (
	statusCodeUninitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
	}
}

func (h *ResponseWriter) WriteHeader(statusCode int) {
	h.ResponseWriter.WriteHeader(statusCode)
	h.statusCode = statusCode
}

func (h *ResponseWriter) GetStatusCodeOrPanic() int {
	if h.statusCode == statusCodeUninitialized {
		panic("status code not set")
	}

	return h.statusCode
}
