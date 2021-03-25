package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name
const (
	BrandColl = "brand"
)

// Logo contains brand logo image information
type Logo struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IMG
}

// BrandFeaturedImage contains brand featured image information
type BrandFeaturedImage struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IMG
}

// Brand struct contains brand specific data
type Brand struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty" bson:"name,omitempty"`
	RegisteredName string             `json:"registered_name,omitempty" bson:"registered_name,omitempty"`
	Slug           string             `json:"slug,omitempty" bson:"slug,omitempty"`

	Description string `json:"description,omitempty" bson:"description,omitempty"`
	WebsiteLink string `json:"website_link,omitempty" bson:"website_link,omitempty"`

	FeaturedImage *BrandFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	Logo          *Logo               `json:"logo,omitempty" bson:"logo,omitempty"`

	FollowerCount  uint `json:"follower_count,omitempty" bson:"follower_count,omitempty"`
	FollowingCount uint `json:"following_count,omitempty" bson:"following_count,omitempty"`
	PostCount      uint `json:"post_count,omitempty" bson:"post_count,omitempty"`

	Fulfillment *Fulfillment `json:"fulfillment,omitempty" bson:"fulfillment,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// Fulfillment contains fulfillment detail such as email, contact info etc.
type Fulfillment struct {
	Email string `json:"email,omitempty" bson:"email,omitempty"`
}
