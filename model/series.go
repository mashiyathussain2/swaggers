package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name in mongodb
const (
	PebbleSeriesColl string = "pebble_series"
)

type PebbleSeries struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty" bson:"name,omitempty" `
	Thumbnail  *IMG                 `json:"thumbnail,omitempty" bson:"thumbnail,omitempty" `
	PebbleIds  []primitive.ObjectID `json:"pebble_ids,omitempty" bson:"pebble_ids,omitempty" `
	PebbleInfo interface{}          `json:"pebble_info,omitempty" bson:"pebble_info,omitempty"`
	Label      *SeriesLabel         `json:"label,omitempty" bson:"label,omitempty"`
	IsActive   bool                 `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

//SeriesLabel will contain meta datapoint for content
type SeriesLabel struct {
	// Interests []string `json:"interests,omitempty" bson:"interests,omitempty"`
	// AgeGroups []string `json:"age_groups,omitempty" bson:"age_groups,omitempty"`
	Genders []string `json:"genders,omitempty" bson:"genders,omitempty"`
}
