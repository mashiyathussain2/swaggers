package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateCategoryOpts serialize fields for create a new category
type CreateCategoryOpts struct {
	Name          string             `json:"name" validate:"required"`
	ParentID      primitive.ObjectID `json:"parent_id"`
	Thumbnail     img                `json:"thumbnail" validate:"required"`
	FeaturedImage img                `json:"featured_image" validate:"required"`
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
