package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Tinee/prog-image"
	"github.com/pkg/errors"
)

// ConvertFunc that writes an image to a writer
type ConvertFunc func(io.Writer, image.Image) error

// Convert takes a domain image and tries to run it through a ConvertFunc
func Convert(pi progimage.Image, c ConvertFunc) (*progimage.Image, error) {
	rd := bytes.NewReader(pi.Body)
	i, format, err := image.Decode(rd)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to decode the bytes to an image.")
	}

	buff := bytes.Buffer{}
	c(&buff, i)

	ct := http.DetectContentType(buff.Bytes())
	ex, err := mime.ExtensionsByType(ct)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to determine the mimeType")
	}
	id := strings.Replace(pi.ID, "."+format, ex[0], 1)

	return &progimage.Image{
		Body:        buff.Bytes(),
		ContentType: ct,
		ID:          id,
	}, nil
}

// JPEGConverter write a common image to a JPEG
func JPEGConverter(w io.Writer, i image.Image) error {
	return jpeg.Encode(w, i, &jpeg.Options{
		Quality: 100,
	})
}

// PNGConverter write a common image to a PNG
func PNGConverter(w io.Writer, i image.Image) error {
	return png.Encode(w, i)
}
