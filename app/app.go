package app

import (
	"go-app/server/config"
	"go-app/server/kafka"

	mongostorage "go-app/server/storage/mongodb"

	"github.com/rs/zerolog"
)

// Options contains arguments required to create a new app instance
type Options struct {
	MongoDB *mongostorage.MongoStorage
	Logger  *zerolog.Logger
	Config  *config.APPConfig
}

// App := contains resources to implement business logic
type App struct {
	MongoDB *mongostorage.MongoStorage
	Logger  *zerolog.Logger
	Config  *config.APPConfig

	// List of services this app is implementing
	SNS        SNS
	SES        SES
	User       User
	Customer   Customer
	Brand      Brand
	Influencer Influencer
	Cart       Cart

	// Consumer
	BrandChanges      kafka.Consumer
	InfluencerChanges kafka.Consumer

	// Producer
	BrandFullProducer      kafka.Producer
	InfluencerFullProducer kafka.Producer

	// Processor
	BrandProcessor      *BrandProcessor
	InfluencerProcessor *InfluencerProcessor
}

// NewApp returns new app instance
func NewApp(opts *Options) *App {
	return &App{
		MongoDB: opts.MongoDB,
		Logger:  opts.Logger,
		Config:  opts.Config,
		SNS:     NewSNSImpl(&SNSOpts{Config: &opts.Config.SNSConfig}),
		SES:     NewSESImpl(&SESImplOpts{Config: &opts.Config.SESConfig}),
	}
}

// Close closes all the resources linked with the app
func (a *App) Close() {
	// terminating connections to all consumes
	CloseConsumer(a)
}
