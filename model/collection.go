package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type of Pebble Collections
const (
	HashTagCollection    = "hashtag"
	InfluencerCollection = "influencer"
	BrandCollection      = "brand"
	SeriesCollection     = "series"
)

// Defining the type of Collection Status
const (
	Draft   string = "draft"
	Unlist  string = "unlist"
	Archive string = "archive"
	Publish string = "publish"
)

// CollectionColl defines the name of the collection for Collections
const (
	CollectionColl = "pebble_collection"
)

type SeriesSubCollection struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Thumbnail IMG                `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	SeriesIDs []string           `json:"series_ids,omitempty" bson:"series_ids,omitempty"`
}

type Collection struct {
	ID                  primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Name                string                `json:"name,omitempty" bson:"name,omitempty"`
	Type                string                `json:"type,omitempty" bson:"type,omitempty"`
	Genders             []string              `json:"genders,omitempty" bson:"genders,omitempty"`
	Hashtags            []string              `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	BrandIDs            []string              `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo           []BrandInfo           `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	InfluencerIDs       []string              `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo      []InfluencerInfo      `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	SeriesSubCollection []SeriesSubCollection `json:"series_subcollection,omitempty" bson:"series_subcollection,omitempty"`
	Status              string                `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt           time.Time             `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
