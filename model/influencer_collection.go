package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InfluencerCollectionColl defines the name of the collection for Influencer Collections
const (
	InfluencerCollectionColl = "influencer_collection"
)

// InfluncerCollectionCollection contains Collection specific data such as Name, Image and CatalogIDs for Influencer
type InfluncerCollection struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Slug         string               `json:"slug,omitempty" bson:"slug,omitempty"`
	Image        *IMG                 `json:"image,omitempty" bson:"image,omitempty"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	Order        int                  `json:"order,omitempty" bson:"order,omitempty"`
	Status       string               `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	// FeaturedCatalogIDs []primitive.ObjectID `json:"featured_catalog_ids,omitempty" bson:"featured_catalog_ids,omitempty"`

}
