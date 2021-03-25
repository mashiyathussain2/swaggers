package app

import (
	"go-app/server/config"
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
	// Example       Example
	KeeperCatalog KeeperCatalog
	Brand         Brand
	Category      Category
	Discount      Discount
	Group         Group
	Collection    Collection
	Inventory     Inventory
}

// NewApp returns new app instance
func NewApp(opts *Options) *App {
	return &App{
		MongoDB: opts.MongoDB,
		Logger:  opts.Logger,
		Config:  opts.Config,
	}
}

// Close closes all the resources linked with the app
func (a *App) Close() {
	// terminating connections to all consumes
	CloseConsumer(a)
}
