package schema

import (
	"go-app/model"

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
	ID          primitive.ObjectID             `json:"id,omitempty"`
	Name        string                         `json:"name,omitempty"`
	Image       *model.IMG                     `json:"image,omitempty"`
	CatalogIDs  []primitive.ObjectID           `json:"catalog_ids,omitempty"`
	CatalogInfo []GetCollectionCatalogInfoResp `json:"catalog_info,omitempty"`
}

type GetCollectionESResp struct {
	ID             primitive.ObjectID       `json:"id,omitempty"`
	Name           string                   `json:"name"`
	Type           string                   `json:"type,omitempty"`
	Genders        []string                 `json:"genders,omitempty"`
	Title          string                   `json:"title,omitempty"`
	SubCollections []GetSubCollectionESResp `json:"sub_collections,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Order          int                      `json:"order,omitempty"`
}

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

type SearchOpts struct {
	Query string `qs:"query"`
}

type BrandSearchResp struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	Logo *model.IMG         `json:"logo"`
}

type InfluencerSearchResp struct {
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

type SearchResp struct {
	Catalog    []CatalogSearchResp    `json:"catalog"`
	Brand      []BrandSearchResp      `json:"brand"`
	Influencer []InfluencerSearchResp `json:"influencer"`
	Content    []ContentSearchResp    `json:"content"`
}

type GetActiveCollectionsOpts struct {
	Gender string `qs:"gender"`
}
