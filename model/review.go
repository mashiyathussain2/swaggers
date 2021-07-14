package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ReviewColl = "review"
)

type ReviewStory struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID   primitive.ObjectID `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	BrandID     primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	Rating      *uint              `json:"rating,omitempty" bson:"rating,omitempty"`
	StoryID     primitive.ObjectID `json:"story_id,omitempty" bson:"story_id,omitempty"`
	IsProcessed bool               `json:"is_processed,omitempty" bson:"is_processed,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
