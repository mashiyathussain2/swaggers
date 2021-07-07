package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSizeProfileOpts struct {
	Name  string              `json:"name,omitempty" validate:"required"`
	Specs []map[string]string `json:"specs,omitempty" validate:"required"`
	Image *Img                `json:"image" validate:"required"`
}

type GetAllSizeProfilesResp struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Image *model.IMG         `json:"image,omitempty" bson:"image,omitempty"`
}

type GetSizeProfileResp struct {
	ID    primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string              `json:"name,omitempty" bson:"name,omitempty"`
	Specs []map[string]string `json:"specs,omitempty" bson:"specs,omitempty"`
	Image *model.IMG          `json:"image,omitempty" bson:"image,omitempty"`
}

type GetSizeProfileForBrandResp struct {
	ID    primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string              `json:"name,omitempty" bson:"name,omitempty"`
	Specs []map[string]string `json:"specs,omitempty" bson:"specs,omitempty"`
	Image *model.IMG          `json:"image,omitempty" bson:"image,omitempty"`
}

type AddBrandToSizeProfileOpts struct {
	IDs     []primitive.ObjectID `json:"ids,omitempty" validate:"required"`
	BrandID primitive.ObjectID   `json:"brand_id,omitempty" validate:"required"`
}

type EditSizeProfileOpts struct {
	ID    primitive.ObjectID `json:"id,omitempty" validate:"required"`
	Image *Img               `json:"image" validate:"required"`
}
