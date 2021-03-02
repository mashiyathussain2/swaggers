package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name in mongodb
const (
	ContentColl string = "content"
)

// list of supported type of content
const (
	PebbleType         string = "pebble"
	CatalogContentType string = "catalog_content"
)

// Content contains linked media (image/video) with influencer, catalog or customer
type Content struct {
	//fields required for Linking
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Type          string               `json:"type,omitempty" bson:"type,omitempty"`
	MediaType     string               `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID       primitive.ObjectID   `json:"media_id,omitempty" bson:"media_id,omitempty"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	CustomerID    primitive.ObjectID   `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	Label         *Label               `json:"label,omitempty" bson:"label,omitempty"`

	// Flag to enable content availability when processing is done
	IsProcessed bool `json:"is_processed" bson:"is_processed"`
	// Flag to toggle content visibility
	IsActive bool `json:"is_active" bson:"is_active"`

	Caption  string   `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags []string `json:"hashtags,omitempty" bson:"hashtags,omitempty"`

	//Catalog Linking
	CatalogIDs []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

//Like has user's liking reference wrt a particular content
type Like struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ContentID  primitive.ObjectID `json:"content_id,omitempty" bson:"content_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//View stores the amount of time for which the user has watched a particular content
type View struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ContentID  primitive.ObjectID `json:"content_id,omitempty" bson:"content_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	Duration   time.Duration      `json:"duration,omitempty" bson:"duration,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//Label will contain meta datapoint for content
type Label struct {
	Interests []string `json:"interests,omitempty" bson:"interests,omitempty"`
	AgeGroups []string `json:"age_groups,omitempty" bson:"age_groups,omitempty"`
	Genders   []string `json:"genders,omitempty" bson:"genders,omitempty"`
}
