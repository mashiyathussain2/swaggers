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

type CatalogDiscountInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IsActive bool               `json:"is_active,omitempty" bson:"is_active,omitempty"`
	Type     string             `json:"type,omitempty" bson:"type,omitempty"`
	Value    uint               `json:"value,omitempty" bson:"value,omitempty"`
	// MaxValue will only be applicable in case of PercentOffType type where you want to restrict discount value to a limit.
	MaxValue uint `json:"max_value,omitempty" bson:"max_value,omitempty"`
}

type VariantInfo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Attribute   string             `json:"attribute,omitempty" bson:"attribute,omitempty"`
	InventoryID primitive.ObjectID `json:"inventory_id,omitempty" bson:"inventory_id,omitempty"`
	SKU         string             `json:"sku,omitempty" bson:"sku,omitempty"`
	IsDeleted   bool               `json:"is_deleted" bson:"is_deleted"`
}

type CatalogInfo struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	BrandID       primitive.ObjectID   `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	BrandInfo     *BrandInfo           `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	FeaturedImage *IMG                 `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	BasePrice     *Price               `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   *Price               `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	DiscountInfo  *CatalogDiscountInfo `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	VariantType   string               `json:"variant_type,omitempty" bson:"variant_type,omitempty"`
	Variants      []VariantInfo        `json:"variants,omitempty" bson:"variants,omitempty"`
}

// Content contains linked media (image/video) with influencer, catalog or customer
type Content struct {
	//fields required for Linking
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Type           string               `json:"type,omitempty" bson:"type,omitempty"`
	MediaType      string               `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID        primitive.ObjectID   `json:"media_id,omitempty" bson:"media_id,omitempty"`
	MediaInfo      interface{}          `json:"media_info,omitempty" bson:"media_info,omitempty"`
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

	ViewCount    uint                 `json:"view_count" bson:"view_count"`
	LikeCount    uint                 `json:"like_count" bson:"like_count"`
	LikeIDs      []primitive.ObjectID `json:"like_ids,omitempty" bson:"like_ids,omitempty"`
	LikedBy      []primitive.ObjectID `json:"liked_by,omitempty" bson:"liked_by,omitempty"`
	CommentCount uint                 `json:"comment_count" bson:"comment_count"`

	Caption  string   `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags []string `json:"hashtags,omitempty" bson:"hashtags,omitempty"`

	//Catalog Linking
	CatalogIDs  []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	CatalogInfo []CatalogInfo        `json:"catalog_info" bson:"catalog_info"`

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
