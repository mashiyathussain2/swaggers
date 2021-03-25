//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_s3.go -package=mock go-app/app S3

package app

import (
	"go-app/server/config"
	"log"
	"os"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 contains methods which interact with AWS S3 Specific SDK
type S3 interface {
	GetPutObjectRequestURL(*s3.PutObjectInput) (string, error)
}

// S3Opts contains args required to create S3Impl instance
type S3Opts struct {
	Config *config.S3Config
}

// S3Impl implements S3 interface
type S3Impl struct {
	S3     *s3.S3
	Config *config.S3Config
}

// InitS3 returns a new instance of S3 service
func InitS3(opts *S3Opts) S3 {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(opts.Config.Region),
		Credentials: credentials.NewStaticCredentials(
			opts.Config.S3UploadAccessKeyID,
			opts.Config.S3UploadSecretAccessKey,
			"", // a token will be created when the session it's used.
		),
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// Create S3 service client
	svc := s3.New(sess)
	awss3 := S3Impl{S3: svc, Config: opts.Config}
	return &awss3
}

// GetPutObjectRequestURL calls aws sdk PutObjectRequest method and returns a presigned url.
func (s3i *S3Impl) GetPutObjectRequestURL(input *s3.PutObjectInput) (string, error) {
	req, _ := s3i.S3.PutObjectRequest(input)
	return req.Presign(s3i.Config.PresignedURLValidity * time.Minute)
}
