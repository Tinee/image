package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Handler is a collection of all the service handlers.
type Handler struct {
	ImageHandler *ImageHandler
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	if head == "image" {
		h.ImageHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	logger.Printf("http error: %s (code=%d)", err, code)

	if code == http.StatusInternalServerError {
		err = errors.New("Status Internal Server Error")
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// NotFound writes an API error message to the response.
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}` + "\n"))
}

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func getQueryParam(key string, url *url.URL) (string, error) {
	m := url.Query()
	values := m[key]
	if len(values) == 0 {
		return "", errors.New("no query params in URL")
	}
	if len(values) > 1 {
		return "", errors.New("more then one value in URL")
	}

	return values[0], nil
}

func encodeJSON(w http.ResponseWriter, logger *log.Logger, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}
