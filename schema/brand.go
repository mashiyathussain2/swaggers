package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialMediaOpts struct {
	FollowersCount int `json:"followers_count" validate:"gte=0"`
}

type SocialAccountOpts struct {
	Facebook  *SocialMediaOpts `json:"facebook"`
	Instagram *SocialMediaOpts `json:"instagram"`
	Twitter   *SocialMediaOpts `json:"twitter"`
	Youtube   *SocialMediaOpts `json:"youtube"`
}

// CreateBrandOpts contains and validations required to create a new brand
type CreateBrandOpts struct {
	Name               string             `json:"name" validate:"required"`
	RegisteredName     string             `json:"registered_name" validate:"required"`
	FulfillmentEmail   string             `json:"fulfillment_email" validate:"required,email"`
	FulfillmentCCEmail []string           `json:"fulfillment_cc_email" validate:"dive,email"`
	Domain             string             `json:"domain" validate:"required"`
	Website            string             `json:"website" validate:"required,url"`
	Logo               *Img               `json:"logo" validate:"required"`
	Bio                string             `json:"bio"`
	CoverImg           *Img               `json:"cover_img" validate:"required"`
	SocialAccount      *SocialAccountOpts `json:"social_account"`
}

// CreateBrandResp contains fields to be returned in response to create brand api
type CreateBrandResp struct {
	ID                 primitive.ObjectID   `json:"id"`
	Name               string               `json:"name"`
	RegisteredName     string               `json:"registered_name"`
	FulfillmentEmail   string               `json:"fulfillment_email"`
	FulfillmentCCEmail []string             `json:"fulfillment_cc_email"`
	Domain             string               `json:"domain"`
	Website            string               `json:"website"`
	Logo               *model.IMG           `json:"logo"`
	CoverImg           *model.IMG           `json:"cover_img"`
	Bio                string               `json:"bio,omitempty"`
	SocialAccount      *model.SocialAccount `json:"social_account"`
	CreatedAt          time.Time            `json:"created_at"`
}

// EditBrandOpts contains and validations required to update a new brand
type EditBrandOpts struct {
	ID                 primitive.ObjectID `json:"id" validate:"required"`
	Name               string             `json:"name"`
	RegisteredName     string             `json:"registered_name"`
	FulfillmentEmail   string             `json:"fulfillment_email" validate:"isdefault|email"`
	FulfillmentCCEmail []string           `json:"fulfillment_cc_email" validate:"dive,email"`
	Domain             string             `json:"domain"`
	Website            string             `json:"website" validate:"isdefault|url"`
	Logo               *Img               `json:"logo"`
	CoverImg           *Img               `json:"cover_img"`
	Bio                string             `json:"bio"`
	SocialAccount      *SocialAccountOpts `json:"social_account"`
}

// EditBrandResp contains fields to be returned in edit brand operation
type EditBrandResp struct {
	ID                 primitive.ObjectID   `json:"id"`
	Name               string               `json:"name,omitempty"`
	RegisteredName     string               `json:"registered_name,omitempty"`
	FulfillmentEmail   string               `json:"fulfillment_email,omitempty"`
	FulfillmentCCEmail []string             `json:"fulfillment_cc_email,omitempty"`
	Domain             string               `json:"domain,omitempty"`
	Website            string               `json:"website,omitempty"`
	Logo               *model.IMG           `json:"logo,omitempty"`
	CoverImg           *model.IMG           `json:"cover_img,omitempty"`
	SocialAccount      *model.SocialAccount `json:"social_account,omitempty"`
	Bio                string               `json:"bio,omitempty"`
	CreatedAt          time.Time            `json:"created_at,omitempty"`
	UpdatedAt          time.Time            `json:"updated_at,omitempty"`
}

// GetBrandsByIDOpts contains fields and validations for get multiple brands by ids
type GetBrandsByIDOpts struct {
	IDs []primitive.ObjectID `json:"id"`
}

// GetBrandResp returns fields contaning brand info
type GetBrandResp struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Name               string               `json:"name,omitempty" bson:"name,omitempty"`
	LName              string               `json:"lname,omitempty" bson:"lname,omitempty"`
	RegisteredName     string               `json:"registered_name,omitempty" bson:"registered_name,omitempty"`
	FulfillmentEmail   string               `json:"fulfillment_email,omitempty" bson:"fulfillment_email,omitempty"`
	FulfillmentCCEmail []string             `json:"fulfillment_cc_email,omitempty" bson:"fulfillment_cc_email,omitempty"`
	Domain             string               `json:"domain,omitempty" bson:"domain,omitempty"`
	Website            string               `json:"website,omitempty" bson:"website,omitempty"`
	Logo               *model.IMG           `json:"logo,omitempty" bson:"logo,omitempty"`
	CoverImg           *model.IMG           `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	Bio                string               `json:"bio,omitempty" bson:"bio,omitempty"`
	SocialAccount      *model.SocialAccount `json:"social_account,omitempty" bson:"social_account,omitempty"`
	FollowersID        []primitive.ObjectID `json:"followers_id,omitempty" bson:"followers_id,omitempty"`
	FollowingID        []primitive.ObjectID `json:"following_id,omitempty" bson:"following_id,omitempty"`
	FollowersCount     uint                 `json:"followers_count,omitempty" bson:"followers_count,omitempty"`
	FollowingCount     uint                 `json:"following_count,omitempty" bson:"following_count,omitempty"`
}

type AddBrandFollowerOpts struct {
	BrandID    primitive.ObjectID `json:"id" validate:"required"`
	CustomerID primitive.ObjectID `json:"customer_id" validate:"required"`
}

type BrandKafkaMessage struct {
	ID                 primitive.ObjectID   `json:"_id,omitempty"`
	Name               string               `json:"name,omitempty"`
	LName              string               `json:"lname,omitempty"`
	RegisteredName     string               `json:"registered_name,omitempty"`
	FulfillmentEmail   string               `json:"fulfillment_email,omitempty"`
	FulfillmentCCEmail []string             `json:"fulfillment_cc_email,omitempty"`
	Domain             string               `json:"domain,omitempty"`
	Website            string               `json:"website,omitempty"`
	Logo               *model.IMG           `json:"logo,omitempty"`
	FollowersCount     uint                 `json:"followers_count,omitempty"`
	FollowingCount     uint                 `json:"following_count,omitempty"`
	Bio                string               `json:"bio,omitempty"`
	CoverImg           *model.IMG           `json:"cover_img,omitempty"`
	SocialAccount      *model.SocialAccount `json:"social_account,omitempty"`
	FollowersID        []primitive.ObjectID `json:"followers_id"`
	FollowingID        []primitive.ObjectID `json:"following_id"`
	CreatedAt          time.Time            `json:"created_at,omitempty"`
	UpdatedAt          time.Time            `json:"updated_at,omitempty"`
}

type BrandFullKafkaMessageOpts struct {
	ID                 primitive.ObjectID   `json:"id,omitempty"`
	Name               string               `json:"name,omitempty"`
	LName              string               `json:"lname,omitempty"`
	RegisteredName     string               `json:"registered_name,omitempty"`
	FulfillmentEmail   string               `json:"fulfillment_email,omitempty"`
	FulfillmentCCEmail []string             `json:"fulfillment_cc_email,omitempty"`
	Domain             string               `json:"domain,omitempty"`
	Website            string               `json:"website,omitempty"`
	Logo               *model.IMG           `json:"logo,omitempty"`
	FollowersCount     uint                 `json:"followers_count"`
	FollowingCount     uint                 `json:"following_count"`
	Bio                string               `json:"bio,omitempty"`
	CoverImg           *model.IMG           `json:"cover_img,omitempty"`
	SocialAccount      *model.SocialAccount `json:"social_account,omitempty"`
	FollowersID        []primitive.ObjectID `json:"followers_id"`
	FollowingID        []primitive.ObjectID `json:"following_id"`
	CreatedAt          time.Time            `json:"created_at,omitempty"`
	UpdatedAt          time.Time            `json:"updated_at,omitempty"`
}
