package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCustomerInfoResp contains fields to be returned in response to get customer
type GetCustomerInfoResp struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	// CartID       primitive.ObjectID `json:"cart_id,omitempty" bson:"cart_id,omitempty"`
	FullName     string        `json:"full_name,omitempty" bson:"full_name,omitempty"`
	DOB          time.Time     `json:"dob,omitempty" bson:"dob,omitempty"`
	Gender       *model.Gender `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfileImage *model.IMG    `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
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
