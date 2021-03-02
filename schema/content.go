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
	ID            primitive.ObjectID   `json:"id"`
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
	CustomerID    primitive.ObjectID   `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
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
