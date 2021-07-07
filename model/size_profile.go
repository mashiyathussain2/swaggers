package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	SizeProfileColl string = "size profile"
)

type SizeProfile struct {
	ID       primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string               `json:"name,omitempty" bson:"name,omitempty"`
	Specs    []map[string]string  `json:"specs,omitempty" bson:"specs,omitempty"`
	BrandIDs []primitive.ObjectID `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	Image    *IMG                 `json:"image,omitempty" bson:"image,omitempty"`
}
