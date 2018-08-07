package s3

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Tinee/prog-image"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	s3      progimage.S3
	bucket  string
	session *session.Session
}

// NewClient returns an instance of a Client
func NewClient(bucket string, s *session.Session) *Client {
	// Set's a production ready seed.
	rand.Seed(time.Now().UnixNano())
	s3 := s3.New(s)

	return &Client{s3, bucket, s}
}

func (c *Client) SaveImage(r io.Reader) (string, error) {
	buff := bytes.Buffer{}
	tee := io.TeeReader(r, &buff)

	_, _, err := image.Decode(tee)
	if err != nil {
		return "", errors.Wrap(err, "couldn't validate image")
	}

	contentType := http.DetectContentType(buff.Bytes())
	id := generateID(25)
	_, err = c.s3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(id),
		Body:        bytes.NewReader(buff.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", errors.Wrap(err, "couldn't PutObject to S3")
	}

	return id, nil
}

func (c *Client) Get(id string) (*progimage.Image, error) {
	res, err := c.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(id),
	})
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to download file, %v", err)
	}

	rd := io.LimitReader(res.Body, 10*1024*1024)
	bs, err := ioutil.ReadAll(rd)

	if err != nil {
		return nil, errors.Wrap(err, "couldn't read image to memory")
	}

	i := &progimage.Image{
		Body:        bs,
		ContentType: *res.ContentType,
		ID:          id,
	}
	return i, nil
}

// SetSeed is a good to have function when I unit test the Client.
func (c *Client) SetSeed(seed int64) {
	rand.Seed(seed)
}

func generateID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
