package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name in mongodb
const (
	ContentColl string = "content"
	CommentColl string = "comment"
	ViewColl    string = "view"
	LikeColl    string = "like"
)

// list of supported type of content
const (
	PebbleType         string = "pebble"
	CatalogContentType string = "catalog_content"
	LiveType           string = "live"
)

type BrandInfo struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Logo *IMG               `json:"logo,omitempty" bson:"logo,omitempty"`
}

type InfluencerInfo struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	ProfileImage *IMG               `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
}

type CatalogInfo struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	FeaturedImage *IMG               `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	BasePrice     *Price             `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   *Price             `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
}

// Content contains linked media (image/video) with influencer, catalog or customer
type Content struct {
	//fields required for Linking
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Type           string               `json:"type,omitempty" bson:"type,omitempty"`
	MediaType      string               `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID        primitive.ObjectID   `json:"media_id,omitempty" bson:"media_id,omitempty"`
	InfluencerIDs  []primitive.ObjectID `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo []InfluencerInfo     `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	BrandIDs       []primitive.ObjectID `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo      []BrandInfo          `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	UserID         primitive.ObjectID   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Label          *Label               `json:"label,omitempty" bson:"label,omitempty"`

	// Flag to enable content availability when processing is done
	IsProcessed bool `json:"is_processed" bson:"is_processed"`
	// Flag to toggle content visibility
	IsActive bool `json:"is_active" bson:"is_active"`

	Caption  string   `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags []string `json:"hashtags,omitempty" bson:"hashtags,omitempty"`

	//Catalog Linking
	CatalogIDs  []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	CatalogInfo []CatalogInfo        `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

//Label will contain meta datapoint for content
type Label struct {
	Interests []string `json:"interests,omitempty" bson:"interests,omitempty"`
	AgeGroups []string `json:"age_groups,omitempty" bson:"age_groups,omitempty"`
	Genders   []string `json:"genders,omitempty" bson:"genders,omitempty"`
}

// Comment stores comment with linked content type

type Comment struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// Type of content ie pebble, catalog, live
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//Like has user's liking reference wrt a particular content
type Like struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

//View stores the amount of time for which the user has watched a particular content
type View struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// Type of content ie pebble, catalog, live
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Duration     time.Duration      `json:"duration,omitempty" bson:"duration,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
