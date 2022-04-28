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

// swagger:model getCustomerInfo
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

// swagger:model EmailLoginCustomerOpts
type EmailLoginCustomerOpts struct {
	//  required: true
	Email string `json:"email" validate:"required,email"`
	//  required: true
	Password string `json:"password" validate:"required,min=6"`
}

// EmailLoginCustomerResp contains fields to be returned in respose to customer email login

// swagger:model SuccessfulLogin
type EmailLoginCustomerResp struct {
	// Token after successful login
	Token string `json:"token"`
}

// UpdateCustomerOpts contains fields and validations to update existing customer

// swagger:model UpdateCustomerOpts
type UpdateCustomerOpts struct {
	ID           primitive.ObjectID `json:"id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	FullName     string             `json:"full_name"`
	DOB          time.Time          `json:"dob"`
	Gender       string             `json:"gender"`
	ProfileImage *Img               `json:"profile_image"`
	Email        string             `json:"email"`
	PhoneNo      *PhoneNoOpts       `json:"phone_no"`
}

//AddAddressOpts contains field required to add new address

// swagger:model AddAddressOpts
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
	PlainAddress      string             `json:"plain_address"`
	IsBillingAddress  bool               `json:"is_billing_address" validate:"required"`
	IsShippingAddress bool               `json:"is_shipping_address" validate:"required"`
	IsDefaultAddress  bool               `json:"is_default_address" validate:"required"`
	ContactNumber     *model.PhoneNumber `json:"contact_number" validate:"required"`
}

//AddAddressResp contains field required to add new address

// swagger:model AddAddressResp
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

//EditAddressOpts contains field required to add new address

// swagger:model EditAddressOpts
type EditAddressOpts struct {
	// required: true
	UserID primitive.ObjectID `json:"user_id" validate:"required"`
	// required: true
	AddressID primitive.ObjectID `json:"address_id" validate:"required"`
	// required: true
	DisplayName string `json:"display_name"`
	// required: true
	Line1 string `json:"line1" validate:"required"`
	// required: true
	Line2 string `json:"line2"`
	// required: true
	District string `json:"district"`
	// required: true
	City string `json:"city" validate:"required"`
	// required: true
	State *model.State `json:"state" validate:"required"`
	// required: true
	PostalCode string `json:"postal_code" validate:"required"`
	// required: true
	Country *model.Country `json:"country" validate:"required"`
	// required: true
	PlainAddress string `json:"plain_address" `
	// required: true
	IsBillingAddress bool `json:"is_billing_address" validate:"required"`
	// required: true
	IsShippingAddress bool `json:"is_shipping_address" validate:"required"`
	// required: true
	IsDefaultAddress bool `json:"is_default_address" validate:"required"`
	// required: true
	ContactNumber *model.PhoneNumber `json:"contact_number" validate:"required"`
}

// UpdateCustomerOpts contains fields and validations to update existing customer
type UpdateCustomerOptsV2 struct {
	ID           primitive.ObjectID `json:"id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
	FullName     string             `json:"full_name"`
	DOB          time.Time          `json:"dob"`
	Gender       string             `json:"gender"`
	ProfileImage *Img               `json:"profile_image"`
	Email        string             `json:"email"`
	PhoneNo      *PhoneNoOpts       `json:"phone_no"`

	InfluencerID      primitive.ObjectID     `json:"influencer_id"`
	Username          string                 `json:"username,omitempty"`
	Bio               string                 `json:"bio,omitempty"`
	CoverImg          *Img                   `json:"cover_img,omitempty"`
	ExternalLinks     []string               `json:"external_links,omitempty"`
	SocialAccount     *SocialAccountOpts     `json:"social_account,omitempty"`
	PayoutInformation *PayoutInformationOpts `json:"payout_information,omitempty"`
}
