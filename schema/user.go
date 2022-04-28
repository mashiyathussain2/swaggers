package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PhoneNoOpts contains fields and validations for mobile no

// swagger:model PhoneNoOpts
type PhoneNoOpts struct {
	// Prefix of the number
	// in: string
	// required: true
	Prefix string `json:"prefix" validate:"required,oneof=+91"`
	// Number of the user
	// in: string
	// required: true
	Number string `json:"number" validate:"required"`
}

// CreateUserOpts contains fields and validations required to create a new user.

// swagger:model CreateUserOpts
type CreateUserOpts struct {
	//  description: type of user
	//  required: true
	Type     string       `json:"type" validate:"required,oneof=customer influencer brand"`
	MobileNo *PhoneNoOpts `json:"phone_no"`

	//  description: email of user
	//  required: true
	Email string `json:"email" validate:"required_without=MobileNo|email"`

	//  description: password of user
	//  required: true
	//  Min Length: 6
	Password string `json:"password" validate:"required,min=6"`

	//  description: confirm password
	//  required: true
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// CreateUserResp contains fields to be required in response to create user

// swagger:model CreateUserResp
type CreateUserResp struct {
	// label: primitive
	ID      primitive.ObjectID `json:"id"`
	Type    string             `json:"type"`
	Email   string             `json:"email,omitempty"`
	PhoneNo *model.PhoneNumber `json:"phone_no,omitempty"`
}

// VerifyEmailOpts contains fields and validations required to verify an email

// swagger:model VerifyEmailOpts
type VerifyEmailOpts struct {
	Email            string `json:"email" validate:"required,email"`
	VerificationCode string `json:"verification_code" validate:"required"`
}

// swagger:model CheckEmailOpts
type CheckEmailOpts struct {
	// required: true
	Email string `json:"email" validate:"required,email"`
}

// swagger:model CheckPhoneNoOpts
type CheckPhoneNoOpts struct {
	// Phone no
	// required: true
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
}

// VerifyEmailOpts contains fields and validations required to verify an email

// swagger:model VerifyPhoneNoOpts
type VerifyPhoneNoOpts struct {
	// Phone no
	// required: true
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
	// required: true
	VerificationCode string `json:"verification_code" validate:"required"`
}

// GetUserResp returns fields in response to get user
type GetUserResp struct {
	ID            primitive.ObjectID `json:"id,omitempty"`
	Type          string             `json:"type,omitempty"`
	Role          string             `json:"role,omitempty"`
	Email         string             `json:"email,omitempty"`
	PhoneNo       *model.PhoneNumber `json:"phone_no,omitempty"`
	Username      string             `json:"username,omitempty"`
	EmailVerified bool               `json:"email_verified,omitempty"`
	PhoneVerified bool               `json:"phone_verified,omitempty"`
	CreatedVia    string             `json:"created_via,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty"`
}

type GetUserInfoResp struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	CustomerID      primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Type            string             `json:"type" bson:"type"`
	Role            string             `json:"role,omitempty" bson:"role,omitempty"`
	FullName        string             `json:"full_name" bson:"full_name"`
	ProfileImage    *model.IMG         `json:"profile_image" bson:"profile_image"`
	Email           string             `json:"email" bson:"email"`
	PhoneNo         *model.PhoneNumber `json:"phone_no" bson:"phone_no"`
	Username        string             `json:"username" bson:"username"`
	EmailVerifiedAt time.Time          `json:"email_verified_at" bson:"email_verified_at"`
	PhoneVerifiedAt time.Time          `json:"phone_verified_at" bson:"phone_verified_at"`
	EmailVerified   bool               `json:"email_verified" bson:"email_verified"`
	PhoneVerified   bool               `json:"phone_verified" bson:"phone_verified"`
}

//ForgotPasswordOpts contains fields and validations required to send otp to email to reset password

// swagger:model ForgotPasswordOpts
type ForgotPasswordOpts struct {
	Email string `json:"email" validate:"required,email"`
}

//ResendVerificationEmailOpts contains fields and validations required to send otp to email to reset password

// swagger:model ResendVerificationEmailOpts
type ResendVerificationEmailOpts struct {
	// required: true
	Email string `json:"email" validate:"required,email"`
}

//ResetPasswordOpts contains fields and validations required to change existing user password

// swagger:model ResetPasswordOpts
type ResetPasswordOpts struct {
	// required:true
	Email string `json:"email" validate:"required,email"`
	// required:true
	Password string `json:"password" validate:"required,min=6"`
	// required:true
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	// required:true
	OTP string `json:"otp" validate:"required"`
}

// MobileLoginCustomerUserOpts contains field and validations required to allow mobile login for customer

// swagger:model MobileLoginCustomerUserOpts
type MobileLoginCustomerUserOpts struct {
	// Phone no
	// required: true
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
	// required: true
	OTP string `json:"otp" validate:"required"`
}

// GenerateMobileLoginOTPOpts contains fields and validations to generate mobile login otp

// swagger:model GenerateMobileLoginOTPOpts
type GenerateMobileLoginOTPOpts struct {
	// required:true
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
}

// LoginWithSocial contains fields and validations required to allow customer login for social apps

// swagger:model LoginWithSocial
type LoginWithSocial struct {
	// required:true
	Type string `json:"type" validate:"required,oneof=google facebook"`
	// required:true
	Email string `json:"email" validate:"required,email"`
	// required:true
	FullName string `json:"full_name" validate:"required"`
	// required:true
	ProfileImage *Img `json:"profile_image" validate:"required"`
}

// swagger:model LoginWithApple
type LoginWithApple struct {
	// required:true
	Type     string `json:"type" validate:"required,oneof=apple"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	// required:true
	AppleID string `json:"apple_id" validate:"required"`
}

// swagger:model GetUserInfoByIDOpts
type GetUserInfoByIDOpts struct {
	ID primitive.ObjectID `json:"id" validate:"required"`
}

type KeeperUserLoginOpts struct {
	GoogleID      string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Domain        string `json:"hd"`
}

type CreateOrUpdateKeeperUser struct {
	Email        string
	FullName     string
	ProfileImage *Img
}

type KeeperUserInfoResp struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	UserInfo     *GetUserResp       `json:"user_info,omitempty" bson:"user_info,omitempty"`
	FullName     string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	ProfileImage *model.IMG         `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	UserGroups   []model.UserGroup  `json:"user_groups,omitempty" bson:"user_groups,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UpdateUserEmailOpts struct {
	ID    primitive.ObjectID `json:"id" validate:"required"`
	Email string             `json:"email" validate:"email"`
}

type UpdateUserPhoneNoOpts struct {
	ID      primitive.ObjectID `json:"id" validate:"required"`
	PhoneNo *PhoneNoOpts       `json:"phone_no" validate:"required"`
}

type SetUserGroupsOpts struct {
	UserID     primitive.ObjectID `json:"user_id" validate:"required"`
	UserGroups []model.UserGroup  `json:"user_groups" validate:"required"`
}

type GetKeeperUsersOpts struct {
	Query string `json:"query"`
	Page  uint   `json:"page"`
}

type GetKeeperUsersResp struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName   string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	UserGroups []model.UserGroup  `json:"user_groups,omitempty" bson:"user_groups,omitempty"`
}
