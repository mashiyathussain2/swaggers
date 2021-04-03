package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EditInfluencerOptsOpts create fields and validations required to create an new instance of influencer
type CreateInfluencerOpts struct {
	Name          string             `json:"name" validate:"required"`
	Bio           string             `json:"bio"`
	CoverImg      *Img               `json:"cover_img" validate:"required"`
	ProfileImage  *Img               `json:"profile_image" validate:"required"`
	ExternalLinks []string           `json:"external_links" validate:"required,min=1,dive,min=6"`
	SocialAccount *SocialAccountOpts `json:"social_account"`
}

// CreateInfluencerResp contains fields to be returned in response to create influencer
type CreateInfluencerResp struct {
	ID            primitive.ObjectID   `json:"id"`
	Name          string               `json:"name"`
	Bio           string               `json:"bio"`
	CoverImg      *model.IMG           `json:"cover_img"`
	ProfileImage  *model.IMG           `json:"profile_image"`
	ExternalLinks []string             `json:"external_links"`
	SocialAccount *model.SocialAccount `json:"social_account"`
	CreatedAt     time.Time            `json:"created_at"`
}

// EditInfluencerOpts contains fields and validations required to edit existing influencer
type EditInfluencerOpts struct {
	ID            primitive.ObjectID `json:"id" validate:"required"`
	Name          string             `json:"name"`
	Bio           string             `json:"bio"`
	CoverImg      *Img               `json:"cover_img"`
	ProfileImage  *Img               `json:"profile_image"`
	ExternalLinks []string           `json:"external_links"`
	SocialAccount *SocialAccountOpts `json:"social_account"`
}

// EditInfluencerResp contains fields to be returned in response to edit influencer
type EditInfluencerResp struct {
	ID            primitive.ObjectID   `json:"id"`
	Name          string               `json:"name"`
	Bio           string               `json:"bio"`
	CoverImg      *model.IMG           `json:"cover_img"`
	ProfileImage  *model.IMG           `json:"profile_image"`
	ExternalLinks []string             `json:"external_links"`
	SocialAccount *model.SocialAccount `json:"social_account"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// GetInfluencersByIDOpts contains fields and validations required to get multiple influencer by matching id
type GetInfluencersByIDOpts struct {
	IDs []primitive.ObjectID `json:"id" validate:"required,min=1"`
}

type GetInfluencersByNameOpts struct {
	Name string `json:"name" validate:"required,min=3"`
}

// GetInfluencerResp contains fields to be returned for get influencer function
type GetInfluencerResp struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string               `json:"name,omitempty" bson:"name,omitempty"`
	CoverImg       *model.IMG           `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage   *model.IMG           `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	SocialAccount  *model.SocialAccount `json:"social_account,omitempty" bson:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty" bson:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty" bson:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id,omitempty" bson:"followers_id,omitempty"`
	FollowingID    []primitive.ObjectID `json:"following_id,omitempty" bson:"following_id,omitempty"`
	FollowersCount uint                 `json:"followers_count" bson:"followers_count"`
	FollowingCount uint                 `json:"following_count" bson:"following_count"`
}

type AddInfluencerFollowerOpts struct {
	InfluencerID primitive.ObjectID `json:"id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
}
