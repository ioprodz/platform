package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorOrange = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return colorGreen
	case "POST", "PUT", "PATCH", "DELETE":
		return colorOrange
	default:
		return colorCyan
	}
}

func colorForStatus(code int) string {
	switch {
	case code < 300:
		return colorGreen
	case code < 400:
		return colorCyan
	case code < 500:
		return colorOrange
	default:
		return colorRed
	}
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(rec, r)
		log.Print(fmt.Sprintf("%s%-6s%s %s %s%d%s %s",
			colorForMethod(r.Method), r.Method, colorReset,
			r.URL.Path,
			colorForStatus(rec.status), rec.status, colorReset,
			time.Since(start),
		))
	})
}
