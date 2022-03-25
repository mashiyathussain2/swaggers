package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InfluencerProductColl defines the name of the collection for Influencer Collections
const (
	InfluencerProductColl = "influencer_product"
)

// InfluncerCollectionCollection contains Collection specific data such as Name, Image and CatalogIDs for Influencer
type InfluncerProduct struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	// FeaturedCatalogIDs []primitive.ObjectID `json:"featured_catalog_ids,omitempty" bson:"featured_catalog_ids,omitempty"`

}
