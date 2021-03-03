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
	S3      S3
}

// App := contains resources to implement business logic
type App struct {
	MongoDB *mongostorage.MongoStorage
	S3      S3
	Logger  *zerolog.Logger
	Config  *config.APPConfig

	// List of services this app is implementing
	Media   Media
	Content Content
	// Kafka Consumer

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
	}
}

// Close closes all the resources linked with the app
func (a *App) Close() {
	// terminating connections to all consumes
	CloseConsumer(a)
}
