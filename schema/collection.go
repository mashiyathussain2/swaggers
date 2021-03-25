package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//SubCollectionOpts specifies the data for SubCollectionOpts to be inputted
type SubCollectionOpts struct {
	Name       string               `json:"name" validate:"required"`
	Image      string               `json:"image" `
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"required"`
}

// CreateCollectionOpts serialize the create collection api arguments
type CreateCollectionOpts struct {
	Type          string              `json:"type" validate:"required,oneof=bourbon dial product editorial"`
	Genders       []string            `json:"genders" validate:"required"`
	Title         string              `json:"title" validate:"required"`
	SubCollection []SubCollectionOpts `json:"sub_collection" validate:"required"`
}

// CreateCollectionResp serialize the create collection api response
type CreateCollectionResp struct {
	ID            primitive.ObjectID    `json:"id" validate:"required"`
	Type          string                `json:"type" validate:"required"`
	Name          string                `json:"name" validate:"required"`
	Genders       []string              `json:"genders" validate:"required"`
	Title         string                `json:"title" validate:"required"`
	SubCollection []model.SubCollection `json:"sub_collection" validate:"required"`
}

// EditCollectionOpts serialize the edit collection api arguments
type EditCollectionOpts struct {
	ID      primitive.ObjectID `json:"id" validate:"required"`
	Genders []string           `json:"genders"`
	Title   string             `json:"title"`
}

// AddSubCollectionOpts serialize the add sub collection api arguments
type AddSubCollectionOpts struct {
	ID            primitive.ObjectID `json:"id" validate:"required"`
	SubCollection *SubCollectionOpts `json:"sub_collection" validate:"required"`
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
