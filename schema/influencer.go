package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EditInfluencerOptsOpts create fields and validations required to create an new instance of influencer
type CreateInfluencerOpts struct {
	Name          string             `json:"name" validate:"required"`
	Username      string             `json:"username" validate:"required"`
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
	Username      string               `json:"username"`
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
	Username      string             `json:"username"`
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
	Username      string               `json:"username"`
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
	Username       string               `json:"username,omitempty" bson:"username,omitempty"`
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
	CustomerID   primitive.ObjectID `json:"customer_id" validate:"required"`
}

type InfluencerKafkaMessage struct {
	ID             primitive.ObjectID   `json:"_id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Username       string               `json:"username,omitempty"`
	CoverImg       *model.IMG           `json:"cover_img,omitempty"`
	ProfileImage   *model.IMG           `json:"profile_image,omitempty"`
	SocialAccount  *model.SocialAccount `json:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id,omitempty"`
	FollowingID    []primitive.ObjectID `json:"following_id,omitempty"`
	FollowersCount uint                 `json:"followers_count"`
	FollowingCount uint                 `json:"following_count"`
	CreatedAt      time.Time            `json:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty"`
}
type InfluencerFullKafkaMessageOpts struct {
	ID             primitive.ObjectID   `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Username       string               `json:"username,omitempty"`
	CoverImg       *model.IMG           `json:"cover_img,omitempty"`
	ProfileImage   *model.IMG           `json:"profile_image,omitempty"`
	SocialAccount  *model.SocialAccount `json:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id"`
	FollowingID    []primitive.ObjectID `json:"following_id"`
	FollowersCount uint                 `json:"followers_count"`
	FollowingCount uint                 `json:"following_count"`
	CreatedAt      time.Time            `json:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty"`
}

type LinkUserAccountOpts struct {
	RequestID    primitive.ObjectID `json:"request_id" validate:"required"`
	InfluencerID primitive.ObjectID `json:"influencer_id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
}

type InfluencerAccountRequestOpts struct {
	UserID     primitive.ObjectID `json:"user_id" validate:"required"`
	CustomerID primitive.ObjectID `json:"customer_id" validate:"required"`
	// InfluencerID  primitive.ObjectID `json:"influencer_id" validate:"required"`
	FullName      string             `json:"full_name" validate:"required"`
	Username      string             `json:"username,omitempty"`
	ProfileImage  Img                `json:"profile_image" validate:"required"`
	CoverImage    Img                `json:"cover_image" validate:"required"`
	Bio           string             `json:"bio" validate:"required"`
	Website       string             `json:"website"`
	SocialAccount *SocialAccountOpts `json:"social_account"`
}

type UpdateInfluencerAccountRequestStatusOpts struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	Grant     *bool              `json:"grant" validate:"required"`
	GranteeID primitive.ObjectID
}

type InfluencerAccountRequestInfluencerInfo struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	CoverImg     *model.IMG         `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage *model.IMG         `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
}

type InfluencerAccountRequestCustomerInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"full_name" bson:"full_name"`
	Gender   string             `json:"gender" bson:"gender"`
	DOB      time.Time          `json:"dob" bson:"dob"`
}

type InfluencerAccountRequestUserInfo struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNo *model.PhoneNumber `json:"phone_no" bson:"phone_no"`
	Email   string             `json:"email" bson:"email"`
}

type InfluencerAccountRequestResp struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// InfluencerID   primitive.ObjectID                      `json:"influencer_id" bson:"influencer_id"`
	ProfileImage  *model.IMG                            `json:"profile_image" bson:"profile_image"`
	CoverImage    *model.IMG                            `json:"cover_image" bson:"cover_image"`
	Bio           string                                `json:"bio" bson:"bio"`
	Website       string                                `json:"website" bson:"website"`
	SocialAccount *model.SocialAccount                  `json:"social_account" bson:"social_account"`
	CustomerInfo  *InfluencerAccountRequestCustomerInfo `json:"customer_info" bson:"customer_info"`
	UserInfo      *InfluencerAccountRequestUserInfo     `json:"user_info" bson:"user_info"`
	IsActive      bool                                  `json:"is_active" bson:"is_active"`
	GranteeID     primitive.ObjectID                    `json:"grantee_id" bson:"grantee_id"`
	CreatedAt     time.Time                             `json:"created_at" bson:"created_at"`
	GrantedAt     time.Time                             `json:"granted_at" bson:"granted_at"`
	Status        string                                `json:"status,omitempty" bson:"status,omitempty"`
}

// EditInfluencerAppOpts contains fields and validations required to edit existing influencer
type EditInfluencerAppOpts struct {
	ID primitive.ObjectID `json:"id" validate:"required"`
	// Name          string             `json:"name"`
	Username string `json:"username"`
	// Bio           string             `json:"bio"`
	// CoverImg      *Img               `json:"cover_img"`
	// ProfileImage  *Img               `json:"profile_image"`
	// ExternalLinks []string           `json:"external_links"`
	// SocialAccount *SocialAccountOpts `json:"social_account"`
}
