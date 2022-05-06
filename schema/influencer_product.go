package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model AddInfluencerProductsOpts
type AddInfluencerProductsOpts struct {
	// swagger:strfmt ObjectID
	InfluencerID primitive.ObjectID `json:"influencer_id" validate:"required"`
	// swagger:strfmt ObjectID
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"required,gt=0"`
}

// swagger:model RemoveInfluencerProductsOpts
type RemoveInfluencerProductsOpts struct {
	// swagger:strfmt ObjectID
	InfluencerID primitive.ObjectID `json:"influencer_id" validate:"required"`
	// swagger:strfmt ObjectID
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"required,gt=0"`
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

// swagger:model GetInfluencerProducts
type GetInfluencerProducts struct {
	InfluencerID string `qs:"id" json:"id" validate:"required"`
	Type         string `qs:"type" json:"type" validate:"required"`
	Page         int    `qs:"page" json:"page"`
}
