package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// list of collection name in mongodb
const (
	ContentColl string = "content"
)

// Content contains linked media (image/video) with influencer, catalog or customer
type Content struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type    string             `json:"type,omitempty" bson:"type,omitempty"`
	MediaID primitive.ObjectID `json:"media_id,omitempty" bson:"media_id,omitempty"`
}
