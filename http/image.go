package http

import (
	"strings"

	"github.com/Tinee/prog-image/image"

	// Todo: Read up on why I have to import these to get image.Decode(io.Reader) to work.
	// Issue -> https://github.com/golang/go/issues/9184
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Tinee/prog-image"
)

// ImageHandler represents an HTTP API handler for dials.
type ImageHandler struct {
	Logger  *log.Logger
	Storage progimage.Storage
}

// ServeHTTP takes request and direct them to their correct handler
func (h *ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodGet:
		id, err := getQueryParam("id", r.URL)
		if id == "" || err != nil {
			NotFound(w)
			return
		}
		split := strings.Split(id, ".")

		h.handleGetWithID(split[0], split[1])(w, r)
	default:
		NotFound(w)
	}
}

func (h *ImageHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rd := io.LimitReader(r.Body, 10*1024*1024)

	id, err := h.Storage.SaveImage(rd)
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	encodeJSON(w, h.Logger, map[string]string{
		"data": id,
	})
}

func (h *ImageHandler) handleGetWithID(id string, extension string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		i, err := h.Storage.Get(id)
		if err != nil {
			Error(w, errors.Wrapf(err, "could not get by that id: %v", id), http.StatusBadRequest, h.Logger)
			return
		}

		switch extension {
		case "png":
			i, err = image.Convert(*i, image.PNGConverter)
		case "jpeg":
			i, err = image.Convert(*i, image.JPEGConverter)
		}
		if err != nil {
			Error(w, errors.Wrapf(err, "could not encode image", id), http.StatusBadRequest, h.Logger)
			return
		}

		encodeJSON(w, h.Logger, map[string]interface{}{
			"data": i,
		})
	}
}
