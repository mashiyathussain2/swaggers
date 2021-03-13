//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_sns.go -package=mock go-app/app SNS

package app

import (
	"go-app/server/config"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNS contains methods to implement AWS SNS service
type SNS interface {
	Publish(*sns.PublishInput) (*sns.PublishOutput, error)
}

// SNSImpl implements SNS interface methods
type SNSImpl struct {
	SNS    *sns.SNS
	Config *config.SNSConfig
}

// SNSOpts contains args passed to initialize new instance of SNSImpl
type SNSOpts struct {
	Config *config.SNSConfig
}

// NewSNSImpl returns new instance of SNSImpl
func NewSNSImpl(opts *SNSOpts) SNS {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(opts.Config.Region),
		Credentials: credentials.NewStaticCredentials(
			opts.Config.AccessKeyID,
			opts.Config.SecretAccessKey,
			"",
		),
	})

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Create a IVS client with additional configuration
	svc := sns.New(sess)
	awssns := SNSImpl{SNS: svc, Config: opts.Config}
	return &awssns
}

// Publish methods calls sns-sdk's Publish method
func (snsi *SNSImpl) Publish(opts *sns.PublishInput) (*sns.PublishOutput, error) {
	return snsi.SNS.Publish(opts)
}
