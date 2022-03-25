package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddInfluencerProductsOpts struct {
	InfluencerID primitive.ObjectID   `json:"influencer_id" validate:"required"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" validate:"required,gt=0"`
}

type RemoveInfluencerProductsOpts struct {
	InfluencerID primitive.ObjectID   `json:"influencer_id" validate:"required"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" validate:"required,gt=0"`
}

type InfluencerProductKafkaConsumerMessage struct {
	ID           primitive.ObjectID   `json:"_id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" bson:"influencer_id"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty"`
}

type InfluencerProductKafkaProducerMessage struct {
	ID           primitive.ObjectID   `json:"id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty"`
}

type GetInfluencerProducts struct {
	InfluencerID string `qs:"id" json:"id" validate:"required"`
	Type         string `qs:"type" json:"type" validate:"required"`
	Page         int    `qs:"page" json:"page"`
}
