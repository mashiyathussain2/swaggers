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
	SNS             SNS
	SES             SES
	Elasticsearch   Elasticsearch
	User            User
	Customer        Customer
	Brand           Brand
	Influencer      Influencer
	Cart            Cart
	KeeperUser      KeeperUser
	ExpressCheckout ExpressCheckout
	Wishlist        Wishlist

	// Consumer
	CustomerChanges   kafka.Consumer
	BrandChanges      kafka.Consumer
	InfluencerChanges kafka.Consumer
	DiscountChanges   kafka.Consumer

	// Producer
	BrandFullProducer      kafka.Producer
	InfluencerFullProducer kafka.Producer

	// Processor
	UserProcessor       *UserProcessor
	BrandProcessor      *BrandProcessor
	InfluencerProcessor *InfluencerProcessor
	CartProcessor       *CartProcessor
}

// NewApp returns new app instance
func NewApp(opts *Options) *App {
	return &App{
		MongoDB:       opts.MongoDB,
		Logger:        opts.Logger,
		Config:        opts.Config,
		SNS:           NewSNSImpl(&SNSOpts{Config: &opts.Config.SNSConfig}),
		SES:           NewSESImpl(&SESImplOpts{Config: &opts.Config.SESConfig}),
		Elasticsearch: InitElasticsearch(&ElasticsearchOpts{Config: &opts.Config.ElasticsearchConfig, Logger: opts.Logger}),
	}
}

// Close closes all the resources linked with the app
func (a *App) Close() {
	// terminating connections to all consumes
	CloseConsumer(a)
	CloseProducer(a)
}
