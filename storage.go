package progimage

import (
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage interface {
	SaveImage(io.Reader) (string, error)
	Get(id string) (*Image, error)
}

type S3 interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}
