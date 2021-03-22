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
	S3      S3
}

// App := contains resources to implement business logic
type App struct {
	MongoDB       *mongostorage.MongoStorage
	Elasticsearch Elasticsearch
	S3            S3
	Logger        *zerolog.Logger
	Config        *config.APPConfig

	// List of services this app is implementing
	Media   Media
	Content Content
	Live    Live

	// Update Processor
	ContentUpdateProcessor *ContentUpdateProcessor

	// Kafka Consumer
	LiveComments      kafka.Consumer
	BrandChanges      kafka.Consumer
	InfluencerChanges kafka.Consumer
	CatalogChanges    kafka.Consumer
	ContentChanges    kafka.Consumer

	// Kafka Producer
	LiveCommentProducer kafka.Producer
	ContentFullProducer kafka.Producer
}

// NewApp returns new app instance
func NewApp(opts *Options) *App {
	return &App{
		MongoDB: opts.MongoDB,
		Logger:  opts.Logger,
		Config:  opts.Config,
		S3: InitS3(&S3Opts{
			Config: &opts.Config.S3Config,
		}),
		Elasticsearch: InitElasticsearch(&ElasticsearchOpts{Config: &opts.Config.ElasticsearchConfig, Logger: opts.Logger}),
	}
}

// Close closes all the resources linked with the app
func (a *App) Close() {
	// terminating connections to all consumes
	CloseConsumer(a)
	CloseProducer(a)
}
