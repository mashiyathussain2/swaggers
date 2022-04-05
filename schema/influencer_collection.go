package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateInfluencerCollectionOpts specifies the data for InfluncerCollection to be inputted
type CreateInfluencerCollectionOpts struct {
	Name         string               `json:"name" validate:"required"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" validate:"required"`
	Image        *Img                 `json:"image" validate:"required"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" validate:"required"`
	Order        uint                 `json:"order"`
	IsDraft      bool                 `json:"is_draft"`
	// FeaturedCatalogIDs []primitive.ObjectID `json:"feat_cat_ids" validate:"required"`
}

// InfluencerCollectionResp serialize the create collection api response
type InfluencerCollectionResp struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" bson:"influencer_id"`
	Name         string               `json:"name" bson:"name"`
	Slug         string               `json:"slug" bson:"slug"`
	Image        *model.IMG           `json:"image" bson:"image"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	Status       string               `json:"status" bson:"status"`
	Order        int                  `json:"order" bson:"order"`
	CreatedAt    time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// EditInfluencerCollectionOpts serialize the edit collection api arguments
type EditInfluencerCollectionOpts struct {
	ID         primitive.ObjectID   `json:"id" validate:"required"`
	Name       string               `json:"name"`
	Image      *Img                 `json:"image"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids"`
	Order      int                  `json:"order"`
	Status     string               `json:"status"`
}

type GetInfluencerCollectionsOpts struct {
	InfluencerID string `qs:"id" json:"id"`
	Status       string `qs:"status" json:"status"`
	Page         int    `qs:"page" json:"page"`
}

type InfluencerCollectionKafkaMessage struct {
	ID           primitive.ObjectID   `json:"_id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" bson:"influencer_id"`
	Name         string               `json:"name"`
	Slug         string               `json:"slug" bson:"slug"`
	Image        *Img                 `json:"image" bson:"image"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	Status       string               `json:"status"`
	Order        int                  `json:"order" bson:"order"`
	CreatedAt    time.Time            `json:"created_at,omitempty"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty"`
}

type InfluencerCollectionFullKafkaMessage struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	InfluencerID primitive.ObjectID `json:"influencer_id"`
	// InfluencerInfo *InfluencerInfo       `json:"influencer_info"`
	Name       string               `json:"name"`
	Slug       string               `json:"slug" bson:"slug"`
	Image      *Img                 `json:"image" bson:"image"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	// CatalogInfo    []GetCatalogBasicResp `json:"catalog_info"`
	Status    string    `json:"status"`
	Order     int       `json:"order" bson:"order"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type InfluencerInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
}

// EditInfluencerCollectionAppOpts serialize the edit collection api arguments
type EditInfluencerCollectionAppOpts struct {
	ID           primitive.ObjectID   `json:"id" validate:"required"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" validate:"required"`
	Name         string               `json:"name"`
	Image        *Img                 `json:"image"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids"`
	Order        int                  `json:"order"`
	Status       string               `json:"status"`
}
