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
	MongoDB       *mongostorage.MongoStorage
	Logger        *zerolog.Logger
	Config        *config.APPConfig
	Elasticsearch Elasticsearch

	// List of services this app is implementing
	KeeperCatalog        KeeperCatalog
	Brand                Brand
	Category             Category
	Discount             Discount
	Group                Group
	Collection           Collection
	InfluencerCollection InfluencerCollection
	InfluencerProducts   InfluencerProducts
	Inventory            Inventory
	Review               Review

	// Consumers
	CatalogChanges              kafka.Consumer
	CollectionCatalogChanges    kafka.Consumer
	CollectionChanges           kafka.Consumer
	BrandChanges                kafka.Consumer
	InventoryChanges            kafka.Consumer
	DiscountChanges             kafka.Consumer
	ContentChanges              kafka.Consumer
	GroupChanges                kafka.Consumer
	ReviewChanges               kafka.Consumer
	InfluencerCollectionChanges kafka.Consumer
	InfluencerProductChanges    kafka.Consumer

	CatalogFullProducer          kafka.Producer
	CollectionFullProducer       kafka.Producer
	InfluencerCollectionProducer kafka.Producer
	InfluencerProductProducer    kafka.Producer
	ReviewFullProducer           kafka.Producer

	// Processor
	CatalogProcessor    *CatalogProcessor
	CollectionProcessor *CollectionProcessor
	ReviewProcessor     *ReviewProcessor
}

// NewApp returns new app instance
func NewApp(opts *Options) *App {
	return &App{
		MongoDB:       opts.MongoDB,
		Logger:        opts.Logger,
		Config:        opts.Config,
		Elasticsearch: InitElasticsearch(&ElasticsearchOpts{Config: &opts.Config.ElasticsearchConfig, Logger: opts.Logger}),
	}
}

// Close closes all the resources linked with the app
func (a *App) Close() {
	// terminating connections to all consumes
	CloseConsumer(a)
	CloseProducer(a)
}
