package progimage

import (
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
)

// Storage is the domain storage interface
type Storage interface {
	SaveImage(io.Reader) (string, error)
	Get(id string) (*Image, error)
}

// S3 is being use to hide the aws S3 object under an interface.
type S3 interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}
