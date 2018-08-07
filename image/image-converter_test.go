package image

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Tinee/prog-image"
)

func TestConvert(t *testing.T) {
	jpeg := getDomainImage("../testdata/github.jpg", "image/jpeg")
	png := getDomainImage("../testdata/github.png", "image/png")

	type args struct {
		pi progimage.Image
		c  ConvertFunc
	}
	tests := []struct {
		name            string
		args            args
		wantContentType string
		wantErr         bool
	}{
		{
			name: "If I give a jpeg I expect a png back",
			args: args{
				pi: jpeg,
				c:  PNGConverter,
			},
			wantContentType: "image/png",
			wantErr:         false,
		},
		{
			name: "If I give a png I expect a jpeg back",
			args: args{
				pi: png,
				c:  JPEGConverter,
			},
			wantContentType: "image/jpeg",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.args.pi, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantContentType != got.ContentType {
				t.Errorf("expected (--%v--) but got (--%v--)", tt.wantContentType, got.ContentType)
			}

			if bytes.Equal(tt.args.pi.Body, got.Body) {
				t.Error("The input Body and the output body is the same. That shouldn't happen.")
			}
		})
	}
}

func getDomainImage(path, contentType string) progimage.Image {
	f, _ := os.Open(path)
	bs, _ := ioutil.ReadAll(f)

	return progimage.Image{
		ContentType: contentType,
		Body:        bs,
		ID:          "random",
	}
}
