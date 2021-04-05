package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateLiveStreamOpts contains fields and validations to create a new live stream.
type CreateLiveStreamOpts struct {
	Name           string               `json:"name" validate:"required"`
	FeaturedImage  *Img                 `json:"featured_image" validate:"required"`
	StreamEndImage *Img                 `json:"stream_end_image,omitempty" validate:"required"`
	ScheduledAt    time.Time            `json:"scheduled_at" validate:"required"`
	InfluencerIDs  []primitive.ObjectID `json:"influencer_ids" validate:"required,min=1"`
	CatalogIDs     []primitive.ObjectID `json:"catalog_ids" validate:"required,min=1"`
}

// CreateLiveStreamResp contains field to be returned in response to create live stream
type CreateLiveStreamResp struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string               `json:"name,omitempty" bson:"name,omitempty"`
	Slug           string               `json:"slug,omitempty" bson:"slug,omitempty"`
	InfluencerIDs  []primitive.ObjectID `json:"influencer_ids" bson:"influencer_ids,omitempty"`
	ScheduledAt    time.Time            `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	CatalogIDs     []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	FeaturedImage  *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	StreamEndImage *model.IMG           `json:"stream_end_image,omitempty" bson:"stream_end_image,omitempty"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// GetLiveStreamResp contains fields to be returned in response to get live stream api
type GetLiveStreamResp struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string               `json:"name,omitempty" bson:"name,omitempty"`
	Slug           string               `json:"slug,omitempty" bson:"slug,omitempty"`
	InfluencerIDs  []primitive.ObjectID `json:"influencer_ids" bson:"influencer_ids,omitempty"`
	CatalogIDs     []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	ScheduledAt    time.Time            `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	FeaturedImage  *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	StreamEndImage *model.IMG           `json:"stream_end_image,omitempty" bson:"stream_end_image,omitempty"`
	IVS            *model.IVS           `json:"ivs,omitempty" bson:"ivs,omitempty"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// GetLiveStreamsFilter contains and validates supported filter to get live streams
type GetLiveStreamsFilter struct {
	Page            int       `queryparam:"page"`
	Status          []string  `queryparam:"status"`
	ScheduledAtFrom time.Time `queryparam:"scheduled_at_from"`
	ScheduledAtTo   time.Time `queryparam:"scheduled_at_to"`
	CreatedAtFrom   time.Time `queryparam:"created_at_from"`
	CreatedAtTo     time.Time `queryparam:"created_at_to"`
}

// StartLiveStreamResp contains fields to be returned in response to start live
type StartLiveStreamResp struct {
	StreamKey string `json:"stream_key"`
	IngestURL string `json:"ingest_url"`
}

// CreateLiveCommentOpts contains fields and validations to push a comment in kafka topic and ivs meta data
type CreateLiveCommentOpts struct {
	Type         string             `json:"type"`
	LiveID       primitive.ObjectID `json:"live_id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	ARN          string             `json:"arn" validate:"required"`
	Name         string             `json:"name" validate:"required"`
	ProfileImage *Img               `json:"profile_image" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	CreatedAt    time.Time          `json:"created_at"`
}

// CreateIVSCommentMetaData contains fields to be returned to aws putmeta data api
type CreateIVSCommentMetaData struct {
	Name         string `json:"name"`
	ProfileImage *Img   `json:"profile_image"`
	Description  string `json:"description"`
}

type PushCatalogInfo struct {
	ID            primitive.ObjectID `json:"id" validate:"required"`
	Name          string             `json:"name" validate:"required"`
	FeaturedImage *ImgResp           `json:"featured_image" validate:"required"`
	BasePrice     *PriceOpts         `json:"base_price" validate:"required"`
	RetailPrice   *PriceOpts         `json:"retail_price" validate:"required"`
}
