//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_ivs.go -package=mock go-app/app IVS

package app

import (
	"go-app/server/config"
	"log"
	"os"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ivs"
)

// IVS contains methods to implement AWS IVS service
type IVS interface {
	CreateChannel(string) (*ivs.CreateChannelOutput, error)
	// PutMetadata(*ivs.PutMetadataInput) (*ivs.PutMetadataOutput, error)
	StopStream(string) (*ivs.StopStreamOutput, error)
}

// IVSImpl implements IVS interface methods
type IVSImpl struct {
	IVS    *ivs.IVS
	Config *config.IVSConfig
}

// IVSOpts contains args passed to initialize new instance of IVSImpl
type IVSOpts struct {
	Config *config.IVSConfig
}

// NewIVSImpl returns new instance of IVSImpl
func NewIVSImpl(opts *IVSOpts) IVS {
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
	// Create a IVS client from just a session.
	// svc := ivs.New(mySession)

	// Create a IVS client with additional configuration
	svc := ivs.New(sess)
	awsivs := IVSImpl{IVS: svc, Config: opts.Config}
	return &awsivs
}

// CreateChannel create a new channel
func (ivsi *IVSImpl) CreateChannel(name string) (*ivs.CreateChannelOutput, error) {
	// Creating a new aws IVS-Channel
	opts := &ivs.CreateChannelInput{
		Name:        &name,
		Authorized:  &ivsi.Config.AuthorizeChannel,
		LatencyMode: &ivsi.Config.LatencyMode,
		Type:        &ivsi.Config.ChannelType,
	}
	return ivsi.IVS.CreateChannel(opts)
}

// GetChannel returns the channel info based on arn
func (ivsi *IVSImpl) GetChannel() (*ivs.GetChannelOutput, error) {
	opts := &ivs.GetChannelInput{
		Arn: &ivsi.Config.ARN,
	}
	return ivsi.IVS.GetChannel(opts)
}

// GetStreamKey returns the streaming key
func (ivsi *IVSImpl) GetStreamKey() (*ivs.GetStreamKeyOutput, error) {
	opts := &ivs.GetStreamKeyInput{
		Arn: &ivsi.Config.ARN,
	}
	return ivsi.IVS.GetStreamKey(opts)
}

// PutMetadata sends metadata into stream such as comment, purchase info, live viewers count etc
func (ivsi *IVSImpl) PutMetadata(opts *ivs.PutMetadataInput) (*ivs.PutMetadataOutput, error) {
	return ivsi.IVS.PutMetadata(opts)
}

// StopStream ends the active stream on
func (ivsi *IVSImpl) StopStream(arn string) (*ivs.StopStreamOutput, error) {
	opts := &ivs.StopStreamInput{
		ChannelArn: &ivsi.Config.ARN,
	}
	return ivsi.IVS.StopStream(opts)
}
