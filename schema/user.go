package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PhoneNoOpts contains fields and validations for mobile no
type PhoneNoOpts struct {
	Prefix string `json:"prefix" validate:"required"`
	Number string `json:"number" validate:"required"`
}

// CreateUserOpts contains fields and validations required to create a new user.
type CreateUserOpts struct {
	Type            string       `json:"type" validate:"required,oneof=customer influencer brand"`
	MobileNo        *PhoneNoOpts `json:"phone_no"`
	Email           string       `json:"email" validate:"required_without=MobileNo|email"`
	Password        string       `json:"password" validate:"required,min=6"`
	ConfirmPassword string       `json:"confirm_password" validate:"required,eqfield=Password"`
}

// CreateUserResp contains fields to be required in response to create user
type CreateUserResp struct {
	ID      primitive.ObjectID `json:"id"`
	Type    string             `json:"type"`
	Email   string             `json:"email,omitempty"`
	PhoneNo *model.PhoneNumber `json:"phone_no,omitempty"`
}

// VerifyEmailOpts contains fields and validations required to verify an email
type VerifyEmailOpts struct {
	Email            string `json:"email" validate:"required,email"`
	VerificationCode string `json:"verification_code" validate:"required"`
}

type CheckEmailOpts struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckPhoneNoOpts struct {
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
}

// VerifyEmailOpts contains fields and validations required to verify an email
type VerifyPhoneNoOpts struct {
	PhoneNo          *PhoneNoOpts `json:"phone_no" validate:"required"`
	VerificationCode string       `json:"verification_code" validate:"required"`
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
type ForgotPasswordOpts struct {
	Email string `json:"email" validate:"required,email"`
}

//ResendVerificationEmailOpts contains fields and validations required to send otp to email to reset password
type ResendVerificationEmailOpts struct {
	Email string `json:"email" validate:"required,email"`
}

//ResetPasswordOpts contains fields and validations required to change existing user password
type ResetPasswordOpts struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	OTP             string `json:"otp" validate:"required"`
}

// MobileLoginCustomerUserOpts contains field and validations required to allow mobile login for customer
type MobileLoginCustomerUserOpts struct {
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
	OTP     string       `json:"otp" validate:"required"`
}

// GenerateMobileLoginOTPOpts contains fields and validations to generate mobile login otp
type GenerateMobileLoginOTPOpts struct {
	PhoneNo *PhoneNoOpts `json:"phone_no" validate:"required"`
}

// LoginWithSocial contains fields and validations required to allow customer login for social apps
type LoginWithSocial struct {
	Type         string `json:"type" validate:"required,oneof=google facebook"`
	Email        string `json:"email" validate:"required,email"`
	FullName     string `json:"full_name" validate:"required"`
	ProfileImage *Img   `json:"profile_image" validate:"required"`
}

type LoginWithApple struct {
	Type     string `json:"type" validate:"required,oneof=apple"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	AppleID  string `json:"apple_id" validate:"required"`
}

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
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UpdateUserEmailOpts struct {
	ID    primitive.ObjectID `json:"id" validate:"required"`
	Email string             `json:"email" validate:"email"`
}

type UpdateUserPhoneNoOpts struct {
	ID      primitive.ObjectID `json:"id" validate:"required"`
	PhoneNo *model.PhoneNumber `json:"phone_no" validate:"required"`
}
