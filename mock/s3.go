package mock

import (
	"io"

	progimage "github.com/Tinee/prog-image"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Storage mock.
type Storage struct {
	SaveImageFn      func(r io.Reader) (string, error)
	SaveImageInvoked bool

	GetFn      func(id string) (*progimage.Image, error)
	GetInvoked bool
}

// SaveImage mock.
func (s *Storage) SaveImage(r io.Reader) (string, error) {
	s.SaveImageInvoked = true
	return s.SaveImageFn(r)
}

// Get mock.
func (s *Storage) Get(id string) (*progimage.Image, error) {
	s.GetInvoked = true
	return s.GetFn(id)
}

// S3 mock
type S3 struct {
	PutObjectFn      func(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	PutObjectInvoked bool

	GetObjectFn      func(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	GetObjectInvoked bool
}

// PutObject mock.
func (s *S3) PutObject(obj *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	s.PutObjectInvoked = true
	return s.PutObjectFn(obj)
}

// GetObject mock.
func (s *S3) GetObject(obj *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	s.GetObjectInvoked = true
	return s.GetObjectFn(obj)
}
