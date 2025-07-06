package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type MiddlewareHandler struct{}

type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error

func NewMiddleware() *MiddlewareHandler {
	return &MiddlewareHandler{}
}

func (h *MiddlewareHandler) Middleware(handle Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		start := time.Now()
		query := r.URL.Query()

		elapsed := time.Since(start).Seconds() * 1000
		elapsedStr := strconv.FormatFloat(elapsed, 'f', -1, 64)
		err := handle(w, r, params)
		if err != nil {
			log.WithFields(log.Fields{
				"time":   elapsedStr,
				"method": r.Method,
				"path":   r.URL.Path,
				"query":  query.Encode(),
				"tags":   []string{"request"},
				// "body":   string(body),
			}).Warning(err.Error())
		} else {
			log.WithFields(log.Fields{
				"time":   elapsedStr,
				"method": r.Method,
				"path":   r.URL.Path,
				"query":  query.Encode(),
				"tags":   []string{"request"},
				// "body":   string(body),
			}).Info("success")
		}
	}
}

func SetupLogger() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
}
