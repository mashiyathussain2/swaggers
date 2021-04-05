package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddToWishlistOpts struct {
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
}
type RemoveFromWishlistOpts struct {
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
}

type GetWishlistResp struct {
	CatalogInfo []model.CatalogInfo `json:"catalog_info"`
}
type CatalogWishListinfo struct {
	ID            primitive.ObjectID          `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string                      `json:"name,omitempty" bson:"name,omitempty"`
	BrandName     string                      `json:"brand_name,omitempty" bson:"brand_name,omitempty"`
	FeaturedImage *model.CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`

	BasePrice   model.Price `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice model.Price `json:"retail_price,omitempty" bson:"retail_price,omitempty"`

	Status *model.Status `json:"status,omitempty" bson:"status,omitempty"`

	DiscountInfo  *model.DiscountInfoResp `json:"discount_info,omitempty" bson:"discount_info,omitempty"`
	InventoryInfo []model.Inventory       `json:"inventory_info,omitempty" bson:"inventory_info,omitempty"`
}
