package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/Tinee/prog-image/http"
	"github.com/Tinee/prog-image/s3"
)

func main() {
	var (
		addr        = flag.String("port", "3000", "Which port you want the server to run on.")
		bucket      = flag.String("bucket", "", "What's the bucket name?")
		awsRegion   = flag.String("region", "eu-west-2", "which region on AWS is the bucket on?")
		serviceName = flag.String("service_name", "ImageService", "What's the name of the service? will be a prefix for the logger etc")

		l  = log.New(os.Stdout, *serviceName, log.LstdFlags)
		se = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(*awsRegion),
		}))
		s3 = s3.NewClient(*bucket, se)
	)
	flag.Parse()

	s, err := http.NewServer(":"+*addr, &http.Handler{
		ImageHandler: &http.ImageHandler{
			Storage: s3,
			Logger:  l,
		},
	})
	if err != nil {
		log.Fatalf("Error trying to create a server: %v", err)
	}

	err = s.Open()
	defer s.Close()
	if err != nil {
		log.Fatalf("Error trying to open a connection to the server: %v", err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	fmt.Println("Got a signal to terminate process", <-c)
}
