package app

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

// Content contains methods to implement and operation pebble(video-only) content
type Content interface{}

// PebbleImpl implements `Pebble` functionality
type PebbleImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// PebbleOpts contains args required to create a new instance of `PebbleImpl`
type PebbleOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitPebble returns a new instance of `Pebble` Implementation
func InitPebble(opts *PebbleOpts) Content {
	p := PebbleImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &p
}

// CreatePebble creates new create a new pebble document in the document, and generates and returns a token to upload video
func (pi *PebbleImpl) CreatePebble() {}
