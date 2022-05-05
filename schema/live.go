package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateLiveStreamOpts contains fields and validations to create a new live stream.

// swagger:model CreateLiveStreamOpts
type CreateLiveStreamOpts struct {
	Name           string    `json:"name" validate:"required"`
	FeaturedImage  *Img      `json:"featured_image" validate:"required"`
	StreamEndImage *Img      `json:"stream_end_image,omitempty" validate:"required"`
	ScheduledAt    time.Time `json:"scheduled_at" validate:"required"`
	// swagger:strfmt ObjectID
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" validate:"required,min=1"`
	// swagger:strfmt ObjectID
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"required,min=1"`
}

// CreateLiveStreamResp contains field to be returned in response to create live stream

// swagger:model CreateLiveStreamResp
type CreateLiveStreamResp struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`
	// swagger:strfmt ObjectID
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" bson:"influencer_ids,omitempty"`
	ScheduledAt   time.Time            `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	// swagger:strfmt ObjectID
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
	Status         *model.StreamStatus  `json:"status,omitempty" bson:"status,omitempty"`
	StatusHistory  []model.StreamStatus `json:"status_history,omitempty" bson:"status_history,omitempty"`
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

// GetLiveStreamsFilter contains and validates supported filter to get live streams

// swagger:model GetAppLiveStreamsFilter
type GetAppLiveStreamsFilter struct {
	Page int `queryparam:"page"`
}

// StartLiveStreamResp contains fields to be returned in response to start live

// swagger:model StartLiveStreamResp
type StartLiveStreamResp struct {
	StreamKey string `json:"stream_key"`
	IngestURL string `json:"ingest_url"`
}

// CreateLiveCommentOpts contains fields and validations to push a comment in kafka topic and ivs meta data

// swagger:model CreateLiveCommentOpts
type CreateLiveCommentOpts struct {
	Type string `json:"type"`
	// swagger:strfmt ObjectID
	LiveID primitive.ObjectID `json:"live_id" validate:"required"`
	// swagger:strfmt ObjectID
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	ARN          string             `json:"arn" validate:"required"`
	Name         string             `json:"name" validate:"required"`
	ProfileImage *Img               `json:"profile_image" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	CreatedAt    time.Time          `json:"created_at"`
}

// CreateIVSCommentMetaData contains fields to be returned to aws putmeta data api
type CreateIVSCommentMetaData struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	ProfileImage *Img               `json:"profile_image"`
	Description  string             `json:"description"`
}

type CreateIVSCatalogMetaData struct {
	ID primitive.ObjectID `json:"id"`
}

type CreateIVSOrderMetaData struct {
	Name         string   `json:"name"`
	ProfileImage *ImgResp `json:"profile_image"`
}

type CreateIVSNewJoinMetaData struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

type IVSMetaData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// swagger:model PushCatalogOpts
type PushCatalogOpts struct {
	ARN string `json:"arn" validate:"required"`
	// swagger:strfmt ObjectID
	ID primitive.ObjectID `json:"id" validate:"required"`
}

type PushNewOrderOpts struct {
	ARN          string   `json:"arn" validate:"required"`
	Name         string   `json:"name" validate:"required"`
	ProfileImage *ImgResp `json:"profile_image" validate:"required"`
}

// swagger:model GetAppLiveStreamResp
type GetAppLiveStreamResp struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`
	// swagger:strfmt ArrayObjectID
	InfluencerIDs  []primitive.ObjectID `json:"influencer_ids" bson:"influencer_ids,omitempty"`
	InfluencerName string               `json:"influencer_name,omitempty" bson:"influencer_name,omitempty"`
	// swagger:strfmt ArrayObjectID
	CatalogIDs     []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	ScheduledAt    time.Time            `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	FeaturedImage  *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	StreamEndImage *model.IMG           `json:"stream_end_image,omitempty" bson:"stream_end_image,omitempty"`
	Status         *model.StreamStatus  `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// swagger:model JoinLiveStreamResp
type JoinLiveStreamResp struct {
	ARN         string `json:"arn"`
	PlaybackURL string `json:"playbackURL"`
}

// swagger:model PushJoinOpts
type PushJoinOpts struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id"`
	ARN  string             `json:"arn"`
	Name string             `json:"name"`
}

type LiveOrderKafkaMessage struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	ProfileImage *ImgResp           `json:"profile_image"`
}

type PushViewerCount struct {
	ARN string `json:"arn"`
}

type ViewerCountMetadata struct {
	Count uint `json:"count"`
}

// swagger:model GetAppLiveStreamInfluencerResp
type GetAppLiveStreamInfluencerResp struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`
	// swagger:strfmt ObjectID
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" bson:"influencer_ids,omitempty"`
	// swagger:strfmt ObjectID
	CatalogIDs     []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	IVS            *model.IVS           `json:"ivs,omitempty" bson:"ivs,omitempty"`
	ScheduledAt    time.Time            `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	FeaturedImage  *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	StreamEndImage *model.IMG           `json:"stream_end_image,omitempty" bson:"stream_end_image,omitempty"`
	Status         *model.StreamStatus  `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// //Influencer info resp from entity service
// type GetInfluencerBasicESEesp struct {
// 	ID               primitive.ObjectID `json:"id,omitempty"`
// 	Name             string             `json:"name,omitempty"`
// 	Username         string             `json:"username,omitempty"`
// 	ProfileImage     *model.IMG         `json:"profile_image,omitempty"`
// 	IsFollowedByUser bool               `json:"is_followed_by_user,omitempty"`
// }

// type GetInfluencerInfoResp struct {
// 	Success bool                      `json:"success"`
// 	Data    *GetInfluencerBasicESEesp `json:"data"`
// }

// swagger:model GetLiveByInfluencerID
type GetLiveByInfluencerID struct {
	Upcoming  []GetAppLiveStreamInfluencerResp `json:"upcoming" bson:"upcoming"`
	Completed []GetAppLiveStreamInfluencerResp `json:"completed" bson:"completed"`
}
