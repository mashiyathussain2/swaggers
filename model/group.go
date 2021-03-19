package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GroupColl defines collection name in mongo db
const (
	GroupColl string = "group"
)

//Group contain catalog group specific data
type Group struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Basis      string               `json:"basis,omitempty" bson:"basis,omitempty"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	Status     *GroupStatus         `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

//GroupStatus stores group status such as unlist (default), publish, archive
type GroupStatus struct {
	Value     string    `json:"value,omitempty" bson:"value,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
