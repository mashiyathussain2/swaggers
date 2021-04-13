package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name in DB
const (
	CustomerColl string = "customer"
)

// Gender is a single digit representation of gender
/*
	Male -> M
	Female -> F
	Others -> O
*/
type Gender = string

// list of supported gender type
const (
	Invalid Gender = ""
	Male    Gender = "M"
	Female  Gender = "F"
	Others  Gender = "O"
)

// GetGender returns gender object
func GetGender(g string) Gender {
	switch g {
	case "M":
		return Male
	case "F":
		return Female
	case "O":
		return Others
	default:
		return Invalid
	}
}

// Customer represents instance of customer (app user)
type Customer struct {
	ID                    primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	UserID                primitive.ObjectID   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CartID                primitive.ObjectID   `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	FullName              string               `json:"full_name,omitempty" bson:"full_name,omitempty"`
	DOB                   time.Time            `json:"dob,omitempty" bson:"dob,omitempty"`
	Gender                *Gender              `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfileImage          *IMG                 `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	BrandFollowCount      uint                 `json:"brand_follow_count,omitempty" bson:"brand_follow_count,omitempty"`
	InfluencerFollowCount uint                 `json:"influencer_follow_count,omitempty" bson:"influencer_follow_count,omitempty"`
	BrandFollowing        []primitive.ObjectID `json:"brand_following,omitempty" bson:"brand_following,omitempty"`
	InfluencerFollowing   []primitive.ObjectID `json:"influencer_following,omitempty" bson:"influencer_following,omitempty"`
	Addresses             []Address            `json:"addresses,omitempty" bson:"addresses,omitempty"`
	CreatedAt             time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt             time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// State represents a state containing its iso representation and name of the state
type State struct {
	ISOCode string `json:"iso_code,omitempty" bson:"iso_code,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
}

// Country represents country contains its iso representation and name of the country
type Country struct {
	ISOCode string `json:"iso_code,omitempty" bson:"iso_code,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
}

// PhoneNumber represents a contact number contains prefix (country code) and phone number
type PhoneNumber struct {
	Prefix string `json:"prefix,omitempty" bson:"prefix,omitempty"`
	Number string `json:"number,omitempty" bson:"number,omitempty"`
}

// Address contains address related info a user
type Address struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DisplayName       string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
	Line1             string             `json:"line1,omitempty" bson:"line1,omitempty"`
	Line2             string             `json:"line2,omitempty" bson:"line2,omitempty"`
	District          string             `json:"district,omitempty" bson:"district,omitempty"`
	City              string             `json:"city,omitempty" bson:"city,omitempty"`
	State             *State             `json:"state,omitempty" bson:"state,omitempty"`
	PostalCode        string             `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	Country           *Country           `json:"country,omitempty" bson:"country,omitempty"`
	PlainAddress      string             `json:"plain_address,omitempty" bson:"plain_address,omitempty"`
	IsBillingAddress  bool               `json:"is_billing_address,omitempty" bson:"is_billing_address,omitempty"`
	IsShippingAddress bool               `json:"is_shipping_address,omitempty" bson:"is_shipping_address,omitempty"`
	IsDefaultAddress  bool               `json:"is_default_address,omitempty" bson:"is_default_address,omitempty"`
	ContactNumber     *PhoneNumber       `json:"contact_number,omitempty" bson:"contact_number,omitempty"`
}
