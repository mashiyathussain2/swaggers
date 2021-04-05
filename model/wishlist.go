package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//model for wishlist
const (
	WishlistColl string = "wishlist"
)

//Wishlist defines structure for wishlist
type Wishlist struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogIDS []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
