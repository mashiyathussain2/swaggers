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
	CategoryID    []primitive.ObjectID `json:"category_id" validate:"required,gt=0"`
	// HashTags      []string             `json:"hashtags" validate:"required,gt=0"`
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
	CategoryID    []primitive.ObjectID `json:"category_id"`
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
	Path          []model.Path         `json:"category_path,omitempty"`
	HashTags      []string             `json:"hashtags,omitempty"`
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
	Path          []model.Path         `json:"category_path" bson:"category_path,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// GetContentFilter contains list of supported filter to be applied while fetching content from DB
type GetContentFilter struct {
	IsActive    *bool                `json:"is_active"`
	IsProcessed *bool                `json:"is_processed"`
	MediaType   string               `json:"media_type" validate:"isdefault|oneof=image video"`
	Type        string               `json:"type" validate:"isdefault|oneof=pebble catalog_content review_story"`
	BrandIDs    []primitive.ObjectID `json:"brand_ids"`
	CatalogIDs  []primitive.ObjectID `json:"catalog_ids"`
	Hashtags    []string             `json:"hashtags"`
	IDs         []primitive.ObjectID `json:"ids"`

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
}

// CreateVideoReviewContentOpts contains and validates args required to create a catalog story review
type CreateVideoReviewContentOpts struct {
	FileName  string             `json:"file_name" validate:"required"`
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	BrandID   primitive.ObjectID `json:"brand_id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
}

// CreateImageCatalogContentOpts contains and validates args required to create an image content
type CreateImageCatalogContentOpts struct {
	MediaID   primitive.ObjectID `json:"media_id" validate:"required"`
	BrandID   primitive.ObjectID `json:"brand_id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
}

// CreateVideoCatalogContentResp returns content id and video upload token
type CreateVideoCatalogContentResp = CreatePebbleResp

// CreateVideoReviewContentResp returns content id and video upload token
type CreateVideoReviewContentResp = CreatePebbleResp

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
	ResourceType string             `json:"resource_type"`
	ResourceID   primitive.ObjectID `json:"resource_id"`
	Description  string             `json:"description"`
	UserID       primitive.ObjectID `json:"user_id"`
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

type ContentUpdateOpts struct {
	ID             primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Type           string                 `json:"type,omitempty" bson:"type,omitempty"`
	MediaType      string                 `json:"media_type,omitempty" bson:"media_type,omitempty"`
	MediaID        primitive.ObjectID     `json:"media_id,omitempty" bson:"media_id,omitempty"`
	MediaInfo      *GetMediaResp          `json:"media_info,omitempty" bson:"media_info,omitempty"`
	CreatorID      primitive.ObjectID     `json:"creator_id,omitempty" bson:"creator_id,omitempty"`
	InfluencerIDs  []primitive.ObjectID   `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo []model.InfluencerInfo `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	BrandIDs       []primitive.ObjectID   `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo      []model.BrandInfo      `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	UserID         primitive.ObjectID     `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Label          *model.Label           `json:"label,omitempty" bson:"label,omitempty"`
	IsProcessed    bool                   `json:"is_processed" bson:"is_processed"`
	IsActive       bool                   `json:"is_active" bson:"is_active"`
	ViewCount      uint                   `json:"view_count" bson:"view_count"`
	LikeCount      uint                   `json:"like_count" bson:"like_count"`
	LikeIDs        []primitive.ObjectID   `json:"like_ids" bson:"like_ids"`
	LikedBy        []primitive.ObjectID   `json:"liked_by" bson:"liked_by"`
	CommentCount   uint                   `json:"comment_count" bson:"comment_count"`
	Caption        string                 `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags       []string               `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	Paths          []model.Path           `json:"category_path,omitempty" bson:"category_path,omitempty"`
	CatalogIDs     []primitive.ObjectID   `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	CatalogInfo    []model.CatalogInfo    `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	SeriesIDs      []primitive.ObjectID   `json:"series_ids,omitempty" bson:"series_ids,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt    time.Time              `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type KafkaMeta struct {
	ID        interface{}         `bson:"_id,omitempty" json:"_id,omitempty"`
	Timestamp primitive.Timestamp `bson:"ts" json:"ts"`
	Namespace string              `bson:"ns" json:"ns"`
	Operation string              `bson:"op,omitempty" json:"op,omitempty"`
	Updates   interface{}         `bson:"updates,omitempty" json:"updates,omitempty"`
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

type GetPebbleFilter struct {
	UserID    string   `json:"user_id,omitempty" queryparam:"user_id"`
	Genders   []string `json:"genders,omitempty" queryparam:"genders"`
	Interests []string `json:"interests,omitempty" queryparam:"interests"`
	Page      uint     `json:"page,omitempty" queryparam:"page"`
	IsSeries  bool
}

type GetPebbleByIDFilter struct {
	UserID string `queryparam:"user_id"`
	ID     string `queryparam:"id"`
}

type GetPebbleMediaESResp struct {
	ID            primitive.ObjectID `json:"id,omitempty"`
	IsPortrait    bool               `json:"is_portrait,omitempty"`
	Dimensions    *model.Dimensions  `json:"dimensions,omitempty"`
	PlaybackURL   string             `json:"hls_playback_url,omitempty"`
	Duration      float32            `json:"duration,omitempty"`
	ThumbnailURLS []string           `json:"thumbnail_url,omitempty"`
}

type GetPebbleESResp struct {
	ID             primitive.ObjectID     `json:"id,omitempty"`
	Type           string                 `json:"type,omitempty"`
	MediaType      string                 `json:"media_type,omitempty"`
	MediaID        primitive.ObjectID     `json:"media_id,omitempty"`
	MediaInfo      *GetPebbleMediaESResp  `json:"media_info,omitempty"`
	InfluencerIDs  []primitive.ObjectID   `json:"influencer_ids,omitempty"`
	InfluencerInfo []model.InfluencerInfo `json:"influencer_info,omitempty"`
	BrandIDs       []primitive.ObjectID   `json:"brand_ids,omitempty"`
	BrandInfo      []model.BrandInfo      `json:"brand_info,omitempty"`
	LikeCount      int                    `json:"like_count,omitempty"`
	CommentCount   int                    `json:"comment_count,omitempty"`
	ViewCount      int                    `json:"view_count,omitempty"`
	Paths          []model.Path           `json:"category_path,omitempty" bson:"category_path,omitempty"`
	Caption        string                 `json:"caption,omitempty"`
	Hashtags       []string               `json:"hashtags,omitempty"`
	CatalogIDs     []primitive.ObjectID   `json:"catalog_ids,omitempty"`
	CatalogInfo    []model.CatalogInfo    `json:"catalog_info,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty"`
	IsActive       bool                   `json:"is_active"`
	IsLikedByUser  bool                   `json:"is_liked_by_user,omitempty"`
}

type LikeESResp struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty"`
}

type GetBrandInfoResp struct {
	Success bool              `json:"success"`
	Payload []model.BrandInfo `json:"payload"`
}

type GetInfluencerInfoResp struct {
	Success bool                   `json:"success"`
	Payload []model.InfluencerInfo `json:"payload"`
}

type GetCatalogInfoResp struct {
	Success bool                `json:"success"`
	Payload []model.CatalogInfo `json:"payload"`
}

type ProcessLikeOpts struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type ProcessLikeESResp struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type ProcessViewOpts struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Duration     time.Duration      `json:"duration,omitempty" bson:"duration,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type ProcessViewESResp struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Duration     time.Duration      `json:"duration,omitempty" bson:"duration,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type ProcessCommentOpts struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceType string             `json:"resource_type,omitempty" bson:"resource_type,omitempty"`
	ResourceID   primitive.ObjectID `json:"resource_id,omitempty" bson:"resource_id,omitempty"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type GetPebblesKeeperFilter struct {
	Type          string               `json:"type" validate:"required,oneof=pebble catalog_content"`
	Status        []string             `json:"status"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids"`
	Page          uint                 `json:"page"`
	Caption       string               `json:"caption"`
}

type ChangeContentStatusOpts struct {
	ID       primitive.ObjectID `json:"id,omitempty" validate:"required"`
	IsActive *bool              `json:"is_active" validate:"required"`
}

type GetPebbleByBrandID struct {
	UserID  string `json:"user_id,omitempty" queryparam:"user_id"`
	BrandID string `queryparam:"brand_id"`
	Page    int    `queryparam:"page"`
}

type GetPebbleByInfluencerID struct {
	UserID       string `json:"user_id,omitempty" queryparam:"user_id"`
	InfluencerID string `queryparam:"influencer_id"`
	Page         int    `queryparam:"page"`
	IsActive     bool
}

type GetCatalogsByInfluencerID struct {
	UserID       string `json:"user_id,omitempty" queryparam:"user_id"`
	InfluencerID string `queryparam:"influencer_id"`
	Page         int    `queryparam:"page"`
}

type GetPebbleByHashtag struct {
	UserID  string `json:"user_id,omitempty" queryparam:"user_id"`
	Hashtag string `json:"hashtag,omitempty" queryparam:"hashtag"`
	Page    int    `json:"page,omitempty" queryparam:"page"`
}

type GetPebbleSearchCaptionResp struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	MediaInfo *GetMediaResp      `json:"media_info,omitempty" bson:"media_info,omitempty"`
}

//CreatorLabelOpts will hold the keywords related to pebbles.
type CreatorLabelOpts struct {
	// Interests []string `json:"interests" validate:"required,min=1"`
	AgeGroup []string `json:"age_group"`
	Gender   []string `json:"gender" validate:"required,min=1,dive,oneof=M F O"`
}

//CreatorEditLabelOpts contains and validates fields to update Label of a content
type CreatorEditLabelOpts struct {
	// Interests []string `json:"interests"`
	AgeGroup []string `json:"age_group"`
	Gender   []string `json:"gender" validate:"dive,oneof=M F O"`
}

// CreatePebbleAppOpts contains and validates args required to create a pebble
type CreatePebbleAppOpts struct {
	FileName      string               `json:"file_name" validate:"required"`
	Caption       string               `json:"caption" validate:"required"`
	CreatorID     primitive.ObjectID   `json:"creator_id" validate:"required"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" validate:"required,min=1"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids"`
	// Label         *CreatorLabelOpts    `json:"label" validate:"required"`
	CategoryID []primitive.ObjectID `json:"category_id" validate:"required,gt=0"`
}

type CreatePebbleAppV2Opts struct {
	FileName string `json:"file_name" validate:"required"`
	// Caption       string               `json:"caption" validate:"required"`
	CreatorID primitive.ObjectID `json:"creator_id" validate:"required"`
	// InfluencerIDs []primitive.ObjectID `json:"influencer_ids" validate:"required,min=1"`
	// BrandIDs   []primitive.ObjectID `json:"brand_ids"`
	// CatalogIDs []primitive.ObjectID `json:"catalog_ids"`
	// Label         *CreatorLabelOpts    `json:"label" validate:"required"`
	// CategoryID []primitive.ObjectID `json:"category_id" validate:"required,gt=0"`
}

// EditPebbleAppOpts contains and validates args required to update an existing pebble content
type EditPebbleAppOpts struct {
	ID            primitive.ObjectID    `json:"id,omitempty" validate:"required"`
	CreatorID     primitive.ObjectID    `json:"creator_id" validate:"required"`
	Caption       string                `json:"caption"`
	InfluencerIDs []primitive.ObjectID  `json:"influencer_ids"`
	BrandIDs      *[]primitive.ObjectID `json:"brand_ids"`
	CatalogIDs    *[]primitive.ObjectID `json:"catalog_ids"`
	// Label         *CreatorEditLabelOpts `json:"label"`
	IsActive   *bool                `json:"is_active"`
	CategoryID []primitive.ObjectID `json:"category_id"`
}

type GetPebblesCreatorFilter struct {
	Status        []string             `json:"status"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids"`
	Page          uint                 `json:"page"`
}

// EditPebbleAppResp contains fields to be returned in EditPebble operation
type EditPebbleAppResp struct {
	ID            primitive.ObjectID    `json:"id"  validate:"required"`
	Caption       string                `json:"caption,omitempty"`
	InfluencerIDs []primitive.ObjectID  `json:"influencer_ids,omitempty"`
	BrandIDs      *[]primitive.ObjectID `json:"brand_ids,omitempty"`
	CatalogIDs    *[]primitive.ObjectID `json:"catalog_ids,omitempty"`
	// Label         *CreatorEditLabelOpts `json:"label,omitempty"`
	IsActive *bool        `json:"is_active,omitempty"`
	Path     []model.Path `json:"category_path,omitempty"`
	HashTags []string     `json:"hashtags,omitempty"`
}

// CreatorGetContentResp contains fields to be returned while querying for content
type CreatorGetContentResp struct {
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
	IsActive      bool                 `json:"is_active" bson:"is_active"`
	Caption       string               `json:"caption,omitempty" bson:"caption,omitempty"`
	Hashtags      []string             `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	Path          []model.Path         `json:"category_path" bson:"category_path,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
