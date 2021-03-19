package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//LabelOpts will hold the keywords related to pebbles.
type LabelOpts struct {
	Interests []string `json:"interests" validate:"required,min=1"`
	AgeGroup  []string `json:"age_group"`
	Gender    []string `json:"gender" validate:"required,min=1,dive,oneof=M F O"`
}

// CreatePebbleOpts contains and validates args required to create a pebble
type CreatePebbleOpts struct {
	FileName      string               `json:"file_name" validate:"required"`
	Caption       string               `json:"caption" validate:"required"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" validate:"required,min=1"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids" validate:"required,min=1"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids"`
	Label         *LabelOpts           `json:"label" validate:"required"`
}

//CreatePebbleResp returns token required for uploading the content to S3 in the background
type CreatePebbleResp struct {
	ID    primitive.ObjectID `json:"id"`
	Token string             `json:"token"`
}

//EditLabelOpts contains and validates fields to update Label of a content
type EditLabelOpts struct {
	Interests []string `json:"interests"`
	AgeGroup  []string `json:"age_group"`
	Gender    []string `json:"gender" validate:"dive,oneof=M F O"`
}

// EditPebbleOpts contains and validates args required to update an existing pebble content
type EditPebbleOpts struct {
	ID            primitive.ObjectID   `json:"id,omitempty" validate:"required"`
	Caption       string               `json:"caption"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids"`
	Label         *EditLabelOpts       `json:"label"`
	IsActive      *bool                `json:"is_active"`
}

// EditPebbleResp contains fields to be returned in EditPebble operation
type EditPebbleResp struct {
	ID            primitive.ObjectID   `json:"id"  validate:"required"`
	Caption       string               `json:"caption,omitempty"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids,omitempty"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids,omitempty"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids,omitempty"`
	Label         *EditLabelOpts       `json:"label,omitempty"`
	IsActive      *bool                `json:"is_active,omitempty"`
}

// ProcessVideoContentOpts contains fields to mark content as processed and link the media object
type ProcessVideoContentOpts = CreateVideoOpts

// GetContentResp contains fields to be returned while querying for content
type GetContentResp struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	MediaType string             `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID   primitive.ObjectID `json:"media_id,omitempty" bson:"media_id,omitempty"`

	// MediaInfo stores video document when lookup aggregation is applied
	MediaInfo *GetMediaResp `json:"media_info,omitempty" bson:"media_info,omitempty"`

	InfluencerIDs []primitive.ObjectID `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	UserID        primitive.ObjectID   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	Label         *model.Label         `json:"label,omitempty" bson:"label,omitempty"`
	IsActive      bool                 `json:"is_active,omitempty" bson:"is_active,omitempty"`
	Caption       string               `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags      []string             `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	CreatedAt     time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// GetContentFilter contains list of supported filter to be applied while fetching content from DB
type GetContentFilter struct {
	IsActive    *bool                `json:"is_active"`
	IsProcessed *bool                `json:"is_processed"`
	MediaType   string               `json:"media_type" validate:"oneof=image video"`
	Type        string               `json:"type" validate:"oneof=pebble catalog_content"`
	BrandIDs    []primitive.ObjectID `json:"brand_ids"`
	CatalogIDs  []primitive.ObjectID `json:"catalog_ids"`
	Hashtags    []string             `json:"hashtags"`

	// Date range filter applied on CreatedAt field
	From time.Time `json:"from"`
	To   time.Time `json:"to" validate:"gtefield=From"`

	Page uint `json:"page"`
}

// CreateVideoCatalogContentOpts contains and validates args required to create a catalog video-content
type CreateVideoCatalogContentOpts struct {
	FileName  string             `json:"file_name" validate:"required"`
	BrandID   primitive.ObjectID `json:"brand_id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	Label     *LabelOpts         `json:"label" validate:"required"`
}

// CreateImageCatalogContentOpts contains and validates args required to create an image content
type CreateImageCatalogContentOpts struct {
	MediaID   primitive.ObjectID `json:"media_id" validate:"required"`
	BrandID   primitive.ObjectID `json:"brand_id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	Label     *LabelOpts         `json:"label" validate:"required"`
}

// CreateVideoCatalogContentResp returns content id and video upload token
type CreateVideoCatalogContentResp = CreatePebbleResp

// CreateImageCatalogContentResp contains fields to be returned for image catalog content
type CreateImageCatalogContentResp struct {
	ID primitive.ObjectID `json:"id,omitempty"`
}

// EditCatalogContentOpts contains fields and validations required to edit existing catalog content
type EditCatalogContentOpts struct {
	ID       primitive.ObjectID `json:"id" validate:"required"`
	IsActive *bool              `json:"is_active,omitempty"`
	Label    *EditLabelOpts     `json:"label,omitempty"`
}

// EditCatalogContentResp contains fields to be returned in respose of edit catalog content
type EditCatalogContentResp struct {
	ID       primitive.ObjectID `json:"id" validate:"required"`
	IsActive *bool              `json:"is_active,omitempty"`
	Label    *EditLabelOpts     `json:"label,omitempty"`
}

type CreateCommentOpts struct {
	ResourceType string             `json:"resource_type" validate:"required,oneof=live pebble"`
	ResourceID   primitive.ObjectID `json:"resource_id" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	CreatedAt    time.Time          `json:"created_at"`
}

type CreateCommentResp struct {
	ID           primitive.ObjectID `json:"id"`
	ResourceType string             `json:"resource_type" validate:"required,oneof=live pebble"`
	ResourceID   primitive.ObjectID `json:"resource_id" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	CreatedAt    time.Time          `json:"created_at"`
}

type CreateViewOpts struct {
	ResourceType string             `json:"resource_type" validate:"required,oneof=live pebble"`
	ResourceID   primitive.ObjectID `json:"resource_id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	Duration     time.Duration      `json:"duration" validate:"required"`
	// Timestamp of instance when user started watching video
	CreatedAt time.Time `json:"created_at"`
}

type CreateLikeOpts struct {
	ResourceType string             `json:"resource_type" validate:"required,oneof=live pebble"`
	ResourceID   primitive.ObjectID `json:"resource_id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
}

type ContentAddOpts struct {
	ID             string                 `json:"_id,omitempty" bson:"_id,omitempty"`
	Type           string                 `json:"type,omitempty" bson:"type,omitempty"`
	MediaType      string                 `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID        string                 `json:"media_id,omitempty" bson:"media_id,omitempty"`
	MediaInfo      *GetMediaResp          `json:"media_info,omitempty" bson:"media_info,omitempty"`
	InfluencerIDs  []string               `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo []model.InfluencerInfo `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	BrandIDs       []string               `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo      []model.BrandInfo      `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	UserID         string                 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Label          *model.Label           `json:"label,omitempty" bson:"label,omitempty"`
	IsProcessed    bool                   `json:"is_processed" bson:"is_processed"`
	IsActive       bool                   `json:"is_active" bson:"is_active"`
	Caption        string                 `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags       []string               `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	CatalogIDs     []string               `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	CatalogInfo    []model.CatalogInfo    `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt    time.Time              `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type KafkaMeta struct {
	ID        interface{}         `bson:"_id,omitempty" json:"_id,omitempty"`
	Timestamp primitive.Timestamp `bson:"ts" json:"ts"`
	Namespace string              `bson:"ns" json:"ns"`
	Operation string              `bson:"op,omitempty" json:"op,omitempty"`
}

type KafkaMessage struct {
	Meta KafkaMeta              `bson:"meta" json:"meta"`
	Data map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`
}

type UpdateContentBrandInfoOpts struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Logo *model.IMG         `json:"logo,omitempty" bson:"logo,omitempty"`
}

type UpdateContentInfluencerInfoOpts struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	ProfileImage *model.IMG         `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
}

type UpdateContentCatalogInfoOpts struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	FeaturedImage *model.IMG         `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	BasePrice     *model.Price       `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   *model.Price       `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
}
