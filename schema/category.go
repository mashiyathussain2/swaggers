package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCategoryOpts serialize fields for create a new category
type CreateCategoryOpts struct {
	Name          string             `json:"name" validate:"required"`
	ParentID      primitive.ObjectID `json:"parent_id"`
	Thumbnail     *Img               `json:"thumbnail"`
	FeaturedImage *Img               `json:"featured_image"`
	IsMain        bool               `json:"is_main"`
}

// CreateCategoryResp contains fields to be returned when a new category is created
type CreateCategoryResp struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Slug          string               `json:"slug,omitempty" bson:"slug,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	ParentID      primitive.ObjectID   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	AncestorID    []primitive.ObjectID `json:"ancestors_id,omitempty" bson:"ancestors_id,omitempty"`
	Thumbnail     *model.IMG           `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	FeaturedImage *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	IsMain        bool                 `json:"is_main,omitempty" bson:"is_main,omitempty"`
}

// EditCategoryOpts contains fields to be updated in db
type EditCategoryOpts struct {
	ID            primitive.ObjectID `json:"id" validate:"required"`
	Name          string             `json:"name"`
	Thumbnail     *Img               `json:"thumbnail"`
	FeaturedImage *Img               `json:"featured_image"`
	IsMain        *bool              `json:"is_main"`
}

// EditCategoryResp contains updated category fields
type EditCategoryResp = CreateCategoryResp

// GetCategoriesResp contains fields to be returned for category doc
type GetCategoriesResp struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	ParentID      primitive.ObjectID   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	AncestorID    []primitive.ObjectID `json:"ancestors_id,omitempty" bson:"ancestors_id,omitempty"`
	Thumbnail     *model.IMG           `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	FeaturedImage *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	IsMain        bool                 `json:"is_main,omitempty" bson:"is_main,omitempty"`
}

// GetCategoriesBasicResp contains id, name and IsMain of category to be returned for category doc
type GetCategoriesBasicResp struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	IsMain bool               `json:"is_main,omitempty" bson:"is_main,omitempty"`
}

// GetMainCategoriesMapResp contains fields to be returned for category map with key as id and this schema as value

// swagger:model GetMainCategoriesMapResp
type GetMainCategoriesMapResp struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	// swagger:strfmt ObjectID
	ParentID primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	// swagger:strfmt ObjectID
	AncestorID    []primitive.ObjectID `json:"ancestors_id,omitempty" bson:"ancestors_id,omitempty"`
	Thumbnail     *model.IMG           `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	FeaturedImage *model.IMG           `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
}

// GetParentCategoriesResp contains fields to be returned for all parent categories

// swagger:model GetParentCategoriesResp
type GetParentCategoriesResp struct {
	// swagger:strfmt ObjectID
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Thumbnail *model.IMG         `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
}

// GetMainCategoriesByParentIDResp contains fields to be returned for all children categories matching parent id

// swagger:model GetMainCategoriesByParentIDResp
type GetMainCategoriesByParentIDResp struct {
	// swagger:strfmt ObjectID
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	FeaturedImage *model.IMG         `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
}

// GetSubCategoriesByParentIDResp contains fields to be returned for all children categories matching parent id

// swagger:model GetSubCategoriesByParentIDResp
type GetSubCategoriesByParentIDResp struct {
	// swagger:strfmt ObjectID
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}
