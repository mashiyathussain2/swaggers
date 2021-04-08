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
	ID          primitive.ObjectID             `json:"id,omitempty"`
	Name        string                         `json:"name,omitempty"`
	Image       *model.IMG                     `json:"image,omitempty"`
	CatalogIDs  []primitive.ObjectID           `json:"catalog_ids,omitempty"`
	CatalogInfo []GetCollectionCatalogInfoResp `json:"catalog_info,omitempty"`
	CreatedAt   time.Time                      `json:"created_at,omitempty"`
	UpdatedAt   time.Time                      `json:"updated_at,omitempty"`
}

type GetCollectionESResp struct {
	ID             primitive.ObjectID       `json:"id,omitempty"`
	Name           string                   `json:"name"`
	Type           string                   `json:"type,omitempty"`
	Genders        []string                 `json:"genders,omitempty"`
	Title          string                   `json:"title,omitempty"`
	SubCollections []GetSubCollectionESResp `json:"sub_collections,omitempty"`
	CreatedAt      time.Time                `json:"created_at,omitempty"`
	UpdatedAt      time.Time                `json:"updated_at,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Order          int                      `json:"order,omitempty"`
}

type GetCatalogBySaleIDOpts struct {
	Page   uint   `qs:"page"`
	SaleID string `qs:"sale_id"`
}

type GetCatalogByCategoryIDOpts struct {
	Page       uint   `qs:"page"`
	CategoryID string `qs:"categoryID"`
}
