package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	InfluencerColl string = "influencer"
)

type Influencer struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	CoverImg      *IMG               `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage  *IMG               `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	SocialAccount *SocialAccount     `json:"social_account,omitempty" bson:"social_account,omitempty"`
	ExternalLinks []string           `json:"external_links,omitempty" bson:"external_links,omitempty"`
	Bio           string             `json:"bio,omitempty" bson:"bio,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
