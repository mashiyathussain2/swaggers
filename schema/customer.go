package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCustomerInfoResp contains fields to be returned in response to get customer
type GetCustomerInfoResp struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CartID       primitive.ObjectID `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	FullName     string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	DOB          time.Time          `json:"dob,omitempty" bson:"dob,omitempty"`
	Gender       *model.Gender      `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfileImage *model.IMG         `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
}

type GetCustomerProfileInfoResp struct {
	ID                    primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserID                primitive.ObjectID   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CartID                primitive.ObjectID   `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	FullName              string               `json:"full_name,omitempty" bson:"full_name,omitempty"`
	DOB                   time.Time            `json:"dob,omitempty" bson:"dob,omitempty"`
	Gender                *model.Gender        `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfileImage          *model.IMG           `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	UserInfo              *GetUserInfoResp     `json:"user_info,omitempty" bson:"user_info,omitempty"`
	BrandFollowing        []primitive.ObjectID `json:"brand_following,omitempty" bson:"brand_following,omitempty"`
	InfluencerFollowing   []primitive.ObjectID `json:"influencer_following,omitempty" bson:"influencer_following,omitempty"`
	BrandFollowCount      uint                 `json:"brand_follow_count,omitempty" bson:"brand_follow_count,omitempty"`
	InfluencerFollowCount uint                 `json:"influencer_follow_count,omitempty" bson:"influencer_follow_count,omitempty"`
}

// EmailLoginCustomerOpts contains fields and validations required to allow customer to login via email
type EmailLoginCustomerOpts struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// EmailLoginCustomerResp contains fields to be returned in respose to customer email login
type EmailLoginCustomerResp struct {
	Token string `json:"email" validate:"required,email"`
}

// UpdateCustomerOpts contains fields and validations to update existing customer
type UpdateCustomerOpts struct {
	ID           primitive.ObjectID `json:"id" validate:"required"`
	FullName     string             `json:"full_name"`
	DOB          time.Time          `json:"dob"`
	Gender       string             `json:"gender"`
	ProfileImage *Img               `json:"profile_image"`
}

//AddAddressOpts contains field required to add new address
type AddAddressOpts struct {
	UserID            primitive.ObjectID `json:"user_id" validate:"required"`
	DisplayName       string             `json:"display_name"`
	Line1             string             `json:"line1" validate:"required"`
	Line2             string             `json:"line2"`
	District          string             `json:"district"`
	City              string             `json:"city" validate:"required"`
	State             *model.State       `json:"state" validate:"required"`
	PostalCode        string             `json:"postal_code" validate:"required"`
	Country           *model.Country     `json:"country" validate:"required"`
	PlainAddress      string             `json:"plain_address" validate:"required"`
	IsBillingAddress  bool               `json:"is_billing_address" validate:"required"`
	IsShippingAddress bool               `json:"is_shipping_address" validate:"required"`
	IsDefaultAddress  bool               `json:"is_default_address" validate:"required"`
	ContactNumber     *model.PhoneNumber `json:"contact_number" validate:"required"`
}

//AddAddressResp contains field required to add new address
type AddAddressResp struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DisplayName   string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
	ContactNumber *model.PhoneNumber `json:"phone_no,omitempty" bson:"contact_number,omitempty"`
	Line1         string             `json:"line1,omitempty" bson:"line1,omitempty"`
	Line2         string             `json:"line2,omitempty" bson:"line2,omitempty"`
	District      string             `json:"district,omitempty" bson:"district,omitempty"`
	City          string             `json:"city,omitempty" bson:"city,omitempty"`
	State         *model.State       `json:"state,omitempty" bson:"state,omitempty"`
	PostalCode    string             `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	PlainAddress  string             `json:"plain_address,omitempty" bson:"plain_address,omitempty"`
}
