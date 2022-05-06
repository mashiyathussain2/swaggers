package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetCollectionVariantResp struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Attribute string             `json:"attribute,omitempty"`
	IsDeleted bool               `json:"is_deleted,omitempty"`
}

type GetCollectionDiscountInfoResp struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Type     model.DiscountType `json:"type,omitempty"`
	Value    uint               `json:"value,omitempty"`
	MaxValue uint               `json:"max_value,omitempty"`
}

type GetCollectionCatalogInfoResp struct {
	ID            primitive.ObjectID             `json:"id,omitempty"`
	BrandID       primitive.ObjectID             `json:"brand_id,omitempty"`
	BrandInfo     *BrandInfoResp                 `json:"brand_info,omitempty"`
	Name          string                         `json:"name,omitempty"`
	FeaturedImage *model.IMG                     `json:"featured_image,omitempty"`
	Slug          string                         `json:"slug,omitempty"`
	VariantType   string                         `json:"variant_type,omitempty"`
	Variants      []GetCollectionVariantResp     `json:"variants,omitempty"`
	BasePrice     *model.Price                   `json:"base_price,omitempty"`
	RetailPrice   *model.Price                   `json:"retail_price,omitempty"`
	DiscountID    primitive.ObjectID             `json:"discount_id,omitempty"`
	DiscountInfo  *GetCollectionDiscountInfoResp `json:"discount_info,omitempty"`
}

type GetSubCollectionESResp struct {
	ID                 primitive.ObjectID             `json:"id,omitempty"`
	Name               string                         `json:"name,omitempty"`
	Image              *model.IMG                     `json:"image,omitempty"`
	CatalogIDs         []primitive.ObjectID           `json:"catalog_ids,omitempty"`
	CatalogInfo        []GetCollectionCatalogInfoResp `json:"catalog_info,omitempty"`
	FeaturedCatalogIDs []primitive.ObjectID           `json:"featured_catalog_ids,omitempty"`
}

// swagger:model GetCollectionESResp
type GetCollectionESResp struct {
	// swagger:strfmt ObjectID
	ID             primitive.ObjectID       `json:"id,omitempty"`
	Name           string                   `json:"name"`
	Type           string                   `json:"type,omitempty"`
	Genders        []string                 `json:"genders,omitempty"`
	Title          string                   `json:"title,omitempty"`
	SubCollections []GetSubCollectionESResp `json:"sub_collections,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Order          int                      `json:"order,omitempty"`
}

// swagger:model GetCatalogBySaleIDOpts
type GetCatalogBySaleIDOpts struct {
	Page   uint   `qs:"page"`
	SaleID string `qs:"sale_id"`
}

type GetCatalogByCategoryIDOpts struct {
	Page       uint     `qs:"page"`
	CategoryID string   `qs:"categoryID"`
	BrandName  []string `qs:"brandName"`
	Sort       int      `qs:"sort"`
}

// swagger:model SearchOpts
type SearchOpts struct {
	Query   string `qs:"query"`
	Page    int    `qs:"page"`
	BrandID string `qs:"brand_id"`
}

// swagger:model BrandSearchResp
type BrandSearchResp struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Logo *model.IMG         `json:"logo"`
}

// swagger:model InfluencerSearchResp
type InfluencerSearchResp struct {
	// swagger:strfmt ObjectID
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	ProfileImage *model.IMG         `json:"profile_image"`
}

type CatalogSearchResp struct {
	ID            primitive.ObjectID `json:"id"`
	Name          string             `json:"name"`
	FeaturedImage *model.IMG         `json:"featured_image"`
	BasePrice     model.Price        `json:"base_price"`
	RetailPrice   model.Price        `json:"retail_price"`
	DiscountInfo  *DiscountBasicResp `json:"discount_info"`
	Variants      []struct {
		ID primitive.ObjectID `json:"id"`
	} `json:"variants"`
	BrandInfoResp *BrandInfoResp `json:"brand_info"`
}

type ContentSearchResp struct {
	ID        primitive.ObjectID `json:"id"`
	Caption   string             `json:"caption"`
	MediaInfo interface{}        `json:"media_info"`
}

// swagger:model SeriesSearchResp
type SeriesSearchResp struct {
	// swagger:strfmt ObjectID
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Thumbnail *Img               `json:"thumbnail"`
}

// swagger:model SearchResp
type SearchResp struct {
	Brand      []BrandSearchResp      `json:"brand"`
	Influencer []InfluencerSearchResp `json:"influencer"`
	Content    []ContentSearchResp    `json:"content"`
	Catalog    []CatalogSearchResp    `json:"catalog"`
}

// swagger:model ShopSearchResp
type ShopSearchResp struct {
	// Brand   []BrandSearchResp   `json:"brand"`
	Catalog []CatalogSearchResp `json:"catalog"`
}

// swagger:model HashtagSearchResp
type HashtagSearchResp struct {
	Text string `json:"text"`
}

// swagger:model DiscoverSearchResp
type DiscoverSearchResp struct {
	Brand      []BrandSearchResp      `json:"brand"`
	Influencer []InfluencerSearchResp `json:"influencer"`
	Series     []SeriesSearchResp     `json:"series"`
	Hashtag    []HashtagSearchResp    `json:"hashtags"`
}

// swagger:model GetActiveCollectionsOpts
type GetActiveCollectionsOpts struct {
	Gender string `qs:"gender"`
	Page   int    `qs:"page"`
	Size   int    `qs:"size"`
}

// swagger:model GetReviewsByCatalogIDFilter
type GetReviewsByCatalogIDFilter struct {
	Page      uint   `qs:"page"`
	CatalogID string `qs:"catalogId"`
}

type GetReviewMediaInfo struct {
	Dimensions  interface{} `json:"dimensions,omitempty"`
	Duration    float32     `json:"duration,omitempty"`
	PlaybackURL string      `json:"hls_playback_url,omitempty"`
}

type GetReviewStoryInfoResp struct {
	ID        primitive.ObjectID  `json:"id,omitempty"`
	MediaType string              `json:"media_type,omitempty"`
	MediaInfo *GetReviewMediaInfo `json:"media_info,omitempty"`
}

// swagger:model GetReviewsByCatalogIDResp
type GetReviewsByCatalogIDResp struct {
	// swagger:strfmt ObjectID
	ID        primitive.ObjectID      `json:"id,omitempty"`
	Rating    *uint                   `json:"rating,omitempty"`
	CreatedAt time.Time               `json:"created_at,omitempty"`
	UpdatedAt time.Time               `json:"updated_at,omitempty"`
	StoryInfo *GetReviewStoryInfoResp `json:"story_info,omitempty"`
	UserInfo  *ReviewUserInfo         `json:"user_info,omitempty"`
}

// swagger:model GetCatalogByBrandIDOpts
type GetCatalogByBrandIDOpts struct {
	Page    uint   `qs:"page"`
	BrandID string `qs:"brand_id"`
}

// swagger:model GetInfluencerCollectionESResp
type GetInfluencerCollectionESResp struct {
	// swagger:strfmt ObjectID
	ID primitive.ObjectID `json:"id,omitempty"`
	// swagger:strfmt ObjectID
	InfluencerID primitive.ObjectID `json:"influencer_id"`
	// InfluencerInfo *InfluencerInfo       `json:"influencer_info"`
	Name  string `json:"name"`
	Slug  string `json:"slug" bson:"slug"`
	Image *Img   `json:"image" bson:"image"`
	// swagger:strfmt ObjectID
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	// CatalogInfo    []GetCatalogBasicResp `json:"catalog_info"`
	Status    string    `json:"status" bson:"status"`
	Order     int       `json:"order" bson:"order"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// swagger:model GetActiveInfluencerCollectionsOpts
type GetActiveInfluencerCollectionsOpts struct {
	InfluencerID string `qs:"influencer_id"`
	Page         int    `qs:"page"`
	Size         int    `qs:"size"`
}

type GetInfluencerProductESResp struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" bson:"influencer_id,omitempty"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type GetPebbleByInfluencerID struct {
	UserID       string `json:"user_id,omitempty" queryparam:"user_id"`
	InfluencerID string `queryparam:"influencer_id"`
	Page         int    `queryparam:"page"`
	IsActive     bool
}

type GetPebbleESResp struct {
	ID            primitive.ObjectID   `json:"id,omitempty"`
	Type          string               `json:"type,omitempty"`
	MediaType     string               `json:"media_type,omitempty"`
	MediaID       primitive.ObjectID   `json:"media_id,omitempty"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids,omitempty"`
	LikeCount     int                  `json:"like_count,omitempty"`
	CommentCount  int                  `json:"comment_count,omitempty"`
	ViewCount     int                  `json:"view_count,omitempty"`
	Paths         []model.Path         `json:"category_path,omitempty" bson:"category_path,omitempty"`
	Caption       string               `json:"caption,omitempty"`
	Hashtags      []string             `json:"hashtags,omitempty"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids,omitempty"`
	// CatalogInfo   []model.CatalogInfo  `json:"catalog_info,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	IsActive      bool      `json:"is_active"`
	IsLikedByUser bool      `json:"is_liked_by_user,omitempty"`
}

type GetCatalogsByInfluencerID struct {
	UserID       string `json:"user_id,omitempty" queryparam:"user_id"`
	InfluencerID string `queryparam:"influencer_id"`
	Page         int    `queryparam:"page"`
}
