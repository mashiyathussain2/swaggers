package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//SubCollectionOpts specifies the data for SubCollectionOpts to be inputted
type SubCollectionOpts struct {
	Name       string               `json:"name" validate:"required"`
	Image      *Img                 `json:"image"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"required"`
}

// CreateCollectionOpts serialize the create collection api arguments
type CreateCollectionOpts struct {
	Type          string              `json:"type" validate:"required,oneof=bourbon dial product editorial"`
	Genders       []string            `json:"genders" validate:"required,dive,oneof=M F O"`
	Title         string              `json:"title" validate:"required"`
	SubCollection []SubCollectionOpts `json:"sub_collections" validate:"required"`
}

// CreateCollectionResp serialize the create collection api response
type CreateCollectionResp struct {
	ID             primitive.ObjectID    `json:"id" bson:"_id"`
	Type           string                `json:"type" bson:"type"`
	Name           string                `json:"name" bson:"name"`
	Genders        []string              `json:"genders" bson:"genders"`
	Title          string                `json:"title" bson:"title"`
	SubCollections []model.SubCollection `json:"sub_collections" bson:"sub_collections"`
	Order          int                   `json:"order" bson:"order"`
}

// CollectionResp serialize the get collections api response
type CollectionResp struct {
	ID             primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string                `json:"name" bson:"name,omitempty"`
	Type           string                `json:"type,omitempty" bson:"type,omitempty"`
	Genders        []string              `json:"genders,omitempty" bson:"genders,omitempty"`
	Title          string                `json:"title,omitempty" bson:"title,omitempty"`
	SubCollections []model.SubCollection `json:"sub_collections,omitempty" bson:"sub_collections,omitempty"`
	Status         string                `json:"status,omitempty" bson:"status,omitempty"`
	Order          int                   `json:"order,omitempty" bson:"order,omitempty"`
}

// EditCollectionOpts serialize the edit collection api arguments
type EditCollectionOpts struct {
	ID      primitive.ObjectID `json:"id" validate:"required"`
	Genders []string           `json:"genders"`
	Title   string             `json:"title"`
	Order   int                `json:"order"`
}

// AddSubCollectionOpts serialize the add sub collection api arguments
type AddSubCollectionOpts struct {
	ID             primitive.ObjectID  `json:"id" validate:"required"`
	SubCollections []SubCollectionOpts `json:"sub_collections" validate:"required"`
}

// UpdateSubCollectionImageOpts serialize the edit collection api arguments
type UpdateSubCollectionImageOpts struct {
	ColID primitive.ObjectID `json:"col_id" validate:"required"`
	SubID primitive.ObjectID `json:"sub_id" validate:"required"`
	Image string             `json:"image" validate:"required"`
}

//UpdateCatalogsInSubCollectionOpts serialize the add or remove catalogs in the sub collection api
type UpdateCatalogsInSubCollectionOpts struct {
	ColID      primitive.ObjectID   `json:"col_id" validate:"required"`
	SubID      primitive.ObjectID   `json:"sub_id" validate:"required"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"required"`
}

type SubCollectionCatalogInfoKafkaMessageResp struct {
	ID            primitive.ObjectID `json:"_id,omitempty"`
	BrandID       primitive.ObjectID `json:"brand_id,omitempty"`
	BrandInfo     *BrandInfoResp     `json:"brand_info,omitempty"`
	Name          string             `json:"name,omitempty"`
	FeaturedImage *model.IMG         `json:"featured_image,omitempty"`
	Slug          string             `json:"slug,omitempty"`
	VariantType   string             `json:"variant_type,omitempty"`
	Variants      []struct {
		ID        primitive.ObjectID `json:"_id,omitempty"`
		Attribute string             `json:"attribute,omitempty"`
		IsDeleted bool               `json:"is_deleted,omitempty"`
	} `json:"variants,omitempty"`
	BasePrice    *model.Price       `json:"base_price,omitempty"`
	RetailPrice  *model.Price       `json:"retail_price,omitempty"`
	DiscountID   primitive.ObjectID `json:"discount_id,omitempty"`
	DiscountInfo *struct {
		ID       primitive.ObjectID `json:"_id,omitempty"`
		Type     model.DiscountType `json:"type,omitempty"`
		Value    uint               `json:"value,omitempty"`
		MaxValue uint               `json:"max_value,omitempty"`
	} `json:"discount_info,omitempty"`
}

type SubCollectionKafkaMessageResp struct {
	ID          primitive.ObjectID                         `json:"_id,omitempty"`
	Name        string                                     `json:"name,omitempty"`
	Image       *model.IMG                                 `json:"image,omitempty"`
	CatalogIDs  []primitive.ObjectID                       `json:"catalog_ids,omitempty"`
	CatalogInfo []SubCollectionCatalogInfoKafkaMessageResp `json:"catalog_info,omitempty"`
	CreatedAt   time.Time                                  `json:"created_at,omitempty"`
	UpdatedAt   time.Time                                  `json:"updated_at,omitempty"`
}

type CollectionKafkaMessageResp struct {
	ID             primitive.ObjectID              `json:"_id,omitempty"`
	Name           string                          `json:"name"`
	Type           string                          `json:"type,omitempty"`
	Genders        []string                        `json:"genders,omitempty"`
	Title          string                          `json:"title,omitempty"`
	SubCollections []SubCollectionKafkaMessageResp `json:"sub_collections,omitempty"`
	CreatedAt      time.Time                       `json:"created_at,omitempty"`
	UpdatedAt      time.Time                       `json:"updated_at,omitempty"`
	Status         string                          `json:"status,omitempty"`
	Order          int                             `json:"order,omitempty"`
}

type SubCollectionInfoResp struct {
	ID          primitive.ObjectID               `json:"id,omitempty"`
	Name        string                           `json:"name,omitempty"`
	Image       *model.IMG                       `json:"image,omitempty"`
	CatalogIDs  []primitive.ObjectID             `json:"catalog_ids,omitempty"`
	CatalogInfo []SubCollectionCatalogInfoSchema `json:"catalog_info,omitempty"`
	CreatedAt   time.Time                        `json:"created_at,omitempty"`
	UpdatedAt   time.Time                        `json:"updated_at,omitempty"`
}

type CollectionInfoResp struct {
	ID             primitive.ObjectID      `json:"id,omitempty"`
	Name           string                  `json:"name"`
	Type           string                  `json:"type,omitempty"`
	Genders        []string                `json:"genders,omitempty"`
	Title          string                  `json:"title,omitempty"`
	SubCollections []SubCollectionInfoResp `json:"sub_collections,omitempty"`
	CreatedAt      time.Time               `json:"created_at,omitempty"`
	UpdatedAt      time.Time               `json:"updated_at,omitempty"`
	Status         string                  `json:"status,omitempty"`
	Order          int                     `json:"order,omitempty"`
}

type SubCollectionCatalogInfoDiscountInfoResp struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Type     model.DiscountType `json:"type,omitempty"`
	Value    uint               `json:"value,omitempty"`
	MaxValue uint               `json:"max_value,omitempty"`
}

type SubCollectionCatalogInfoVariantsResp struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Attribute string             `json:"attribute,omitempty"`
	IsDeleted bool               `json:"is_deleted,omitempty"`
}

type SubCollectionCatalogInfoSchema struct {
	ID            primitive.ObjectID                        `json:"id,omitempty"`
	BrandID       primitive.ObjectID                        `json:"brand_id,omitempty"`
	BrandInfo     *BrandInfoResp                            `json:"brand_info,omitempty"`
	Name          string                                    `json:"name,omitempty"`
	FeaturedImage *model.IMG                                `json:"featured_image,omitempty"`
	Slug          string                                    `json:"slug,omitempty"`
	VariantType   string                                    `json:"variant_type,omitempty"`
	Variants      []SubCollectionCatalogInfoVariantsResp    `json:"variants,omitempty"`
	BasePrice     *model.Price                              `json:"base_price,omitempty"`
	RetailPrice   *model.Price                              `json:"retail_price,omitempty"`
	DiscountID    primitive.ObjectID                        `json:"discount_id,omitempty"`
	DiscountInfo  *SubCollectionCatalogInfoDiscountInfoResp `json:"discount_info,omitempty"`
}

type UpdateCollectionStatus struct {
	ID     primitive.ObjectID `json:"id"`
	Status string             `json:"status"`
}
