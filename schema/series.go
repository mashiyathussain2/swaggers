package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSeriesOpts struct {
	Name      string               `json:"name" validate:"required"`
	Thumbnail Img                  `json:"thumbnail" validate:"required"`
	PebbleIds []primitive.ObjectID `json:"pebble_ids" validate:"required"`
	Label     *SeriesLabelOpts     `json:"label" validate:"required"`
}

type UpdateSeriesOpts struct {
	ID        primitive.ObjectID   `json:"id" validate:"required"`
	Name      string               `json:"name,omitempty"`
	Thumbnail *Img                 `json:"thumbnail,omitempty"`
	PebbleIds []primitive.ObjectID `json:"pebble_ids,omitempty"`
	IsActive  *bool                `json:"is_active,omitempty"`
	Label     *SeriesEditLabelOpts `json:"label"`
}

type PebbleSeriesKafkaUpdateOpts struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty" bson:"name,omitempty" `
	Thumbnail  *model.IMG           `json:"thumbnail,omitempty" bson:"thumbnail,omitempty" `
	PebbleIds  []primitive.ObjectID `json:"pebble_ids,omitempty" bson:"pebble_ids,omitempty" `
	PebbleInfo []ContentForSeries   `json:"pebble_info,omitempty" bson:"pebble_info,omitempty"`
	Label      *model.SeriesLabel   `json:"label,omitempty" bson:"label,omitempty"`
	IsActive   bool                 `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type GetPebbleSeriesFilter struct {
	UserID  string   `json:"user_id,omitempty" queryparam:"user_id"`
	Genders []string `json:"genders,omitempty" queryparam:"genders"`
	// Interests []string `json:"interests,omitempty" queryparam:"interests"`
	Page uint `json:"page,omitempty" queryparam:"page"`
}

// swagger:model GetPebbleSeriesESResp
type GetPebbleSeriesESResp struct {
	// swagger:strfmt ObjectID
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty" `
	Thumbnail  model.IMG          `json:"thumbnail,omitempty" bson:"thumbnail,omitempty" `
	PebbleIds  []interface{}      `json:"pebble_ids,omitempty" bson:"pebble_ids,omitempty" `
	PebbleInfo []*GetPebbleESResp `json:"pebble_info,omitempty" bson:"pebble_info,omitempty"`
	IsActive   bool               `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

//SeriesLabelOpts will hold the keywords related to pebbles series.
type SeriesLabelOpts struct {
	// Interests []string `json:"interests" validate:"required,min=1"`
	// AgeGroup  []string `json:"age_group"`
	Gender []string `json:"gender" validate:"required,min=1,dive,oneof=M F O"`
}

//SeriesEditLabelOpts contains and validates fields to update Label of a series
type SeriesEditLabelOpts struct {
	// Interests []string `json:"interests"`
	// AgeGroup  []string `json:"age_group"`
	Gender []string `json:"gender" validate:"dive,oneof=M F O"`
}

// GetContentRespForSeries contains fields to be returned while querying for content
type GetContentRespForSeries struct {
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

// ContentForSeries contains linked media (image/video) with influencer, catalog or customer
type ContentForSeries struct {
	//fields required for Linking
	ID             primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Type           string                 `json:"type,omitempty" bson:"type,omitempty"`
	MediaType      string                 `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID        primitive.ObjectID     `json:"media_id,omitempty" bson:"media_id,omitempty"`
	MediaInfo      *GetMediaResp          `json:"media_info,omitempty" bson:"media_info,omitempty"`
	InfluencerIDs  []primitive.ObjectID   `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo []model.InfluencerInfo `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	BrandIDs       []primitive.ObjectID   `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo      []model.BrandInfo      `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	UserID         primitive.ObjectID     `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Label          *model.Label           `json:"label,omitempty" bson:"label,omitempty"`

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
	CatalogInfo []model.CatalogInfo  `json:"catalog_info" bson:"catalog_info"`

	SeriesIDs []primitive.ObjectID `json:"series_ids,omitempty" bson:"series_ids,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// swagger:model GetSeriesByIDs
type GetSeriesByIDs struct {
	UserID string   `json:"user_id,omitempty" queryparam:"user_id"`
	ID     []string `json:"id,omitempty" queryparam:"id"`
	Page   int      `json:"page,omitempty" queryparam:"page"`
}

// swagger:model GetPebbleByCategoryIDOpts
type GetPebbleByCategoryIDOpts struct {
	UserID     string `json:"user_id,omitempty" queryparam:"user_id"`
	Page       uint   `qs:"page"`
	CategoryID string `qs:"categoryID"`
	Sort       int    `qs:"sort"`
}

type SearchPebbleByCaption struct {
	Caption string `qs:"caption"`
	Page    uint   `qs:"page"`
}

type GetSeriesKeeperFilter struct {
	Page     int    `qs:"page"`
	IsActive bool   `qs:"is_active"`
	Name     string `qs:"name"`
}

type PebbleSeriesResp struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty" bson:"name,omitempty" `
	Thumbnail  *model.IMG           `json:"thumbnail,omitempty" bson:"thumbnail,omitempty" `
	PebbleIds  []primitive.ObjectID `json:"pebble_ids,omitempty" bson:"pebble_ids,omitempty" `
	PebbleInfo interface{}          `json:"pebble_info,omitempty" bson:"pebble_info,omitempty"`
	Label      *model.SeriesLabel   `json:"label,omitempty" bson:"label,omitempty"`
	IsActive   bool                 `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type KeeperPebbleSeriesResp struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty" bson:"name,omitempty" `
	Thumbnail  *model.IMG           `json:"thumbnail,omitempty" bson:"thumbnail,omitempty" `
	PebbleIds  []primitive.ObjectID `json:"pebble_ids,omitempty" bson:"pebble_ids,omitempty" `
	PebbleInfo []GetContentResp     `json:"pebble_info,omitempty" bson:"pebble_info,omitempty"`
	Label      *model.SeriesLabel   `json:"label,omitempty" bson:"label,omitempty"`
	IsActive   bool                 `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type KeeperGetSeriesByID struct {
	ID string `json:"id,omitempty" queryparam:"id"`
}
type KeeperGetSeriesBasic struct {
	IDs []primitive.ObjectID `json:"ids" validate:"required"`
}
type KeeperGetSeriesBasicResp struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" `
	Thumbnail *model.IMG         `json:"thumbnail,omitempty" bson:"thumbnail,omitempty" `
}

// type SearchSeriesByName struct {
// 	Name string `qs:"name"`
// 	Page uint   `qs:"page"`
// }
