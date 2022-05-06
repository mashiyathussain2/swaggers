package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name
const (
	BrandColl string = "brand"
)

type Brand struct {
	ID                 primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name               string               `json:"name,omitempty" bson:"name,omitempty"`
	LName              string               `json:"lname,omitempty" bson:"lname,omitempty"`
	Username           string               `json:"username,omitempty" bson:"username,omitempty"`
	RegisteredName     string               `json:"registered_name,omitempty" bson:"registered_name,omitempty"`
	FulfillmentEmail   string               `json:"fulfillment_email,omitempty" bson:"fulfillment_email,omitempty"`
	FulfillmentCCEmail []string             `json:"fulfillment_cc_email,omitempty" bson:"fulfillment_cc_email,omitempty"`
	Domain             string               `json:"domain,omitempty" bson:"domain,omitempty"`
	Website            string               `json:"website,omitempty" bson:"website,omitempty"`
	Logo               *IMG                 `json:"logo,omitempty" bson:"logo,omitempty"`
	FollowersCount     uint                 `json:"followers_count,omitempty" bson:"followers_count,omitempty"`
	FollowingCount     uint                 `json:"following_count,omitempty" bson:"following_count,omitempty"`
	Bio                string               `json:"bio,omitempty" bson:"bio,omitempty"`
	CoverImg           *IMG                 `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	SocialAccount      *SocialAccount       `json:"social_account,omitempty" bson:"social_account,omitempty"`
	FollowersID        []primitive.ObjectID `json:"followers_id,omitempty" bson:"followers_id,omitempty"`
	FollowingID        []primitive.ObjectID `json:"following_id,omitempty" bson:"following_id,omitempty"`
	Policies           []Policy             `json:"policies,omitempty" bson:"policies,omitempty"`
	IsCODAvailable     bool                 `json:"is_cod_available" bson:"is_cod_available"`
	CreatedAt          time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt          time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	SizeProfiles       []primitive.ObjectID `json:"size_profiles,omitempty" bson:"size_profiles,omitempty"`
}

// SocialMedia contains followers_count for a specific account
/*
	Type -> facebook
		-> twitter
		-> youtube
		-> instagram
*/
type SocialMedia struct {
	FollowersCount uint   `json:"followers_count,omitempty" bson:"followers_count,omitempty"`
	URL            string `json:"url,omitempty" bson:"url,omitempty"`
}

// SocialAccount contains info about social media pages such as facebook, instagram, etc
type SocialAccount struct {
	Facebook  *SocialMedia `json:"facebook,omitempty" bson:"facebook,omitempty"`
	Instagram *SocialMedia `json:"instagram,omitempty" bson:"instagram,omitempty"`
	Twitter   *SocialMedia `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Youtube   *SocialMedia `json:"youtube,omitempty" bson:"youtube,omitempty"`
}

type BrandClaim struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

// Policy contains brand policies in key:value format
type Policy struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}
