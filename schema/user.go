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

// GetUserResp returns fields in response to get user
type GetUserResp struct {
	ID            primitive.ObjectID `json:"id,omitempty"`
	Type          string             `json:"type,omitempty"`
	Role          string             `json:"role,omitempty"`
	Email         string             `json:"email,omitempty"`
	PhoneNo       *model.PhoneNumber `json:"phone_no,omitempty"`
	Username      string             `json:"username,omitempty"`
	Password      string             `json:"password,omitempty"`
	EmailVerified bool               `json:"email_verified,omitempty"`
	PhoneVerified bool               `json:"phone_verified,omitempty"`
	CreatedVia    string             `json:"created_via,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty"`
}

type GetUserInfoResp struct {
	ID            primitive.ObjectID `json:"id"`
	CustomerID    primitive.ObjectID `json:"customer_id"`
	Type          string             `json:"type"`
	FullName      string             `json:"full_name"`
	ProfileImage  *model.IMG         `json:"profile_image"`
	Email         string             `json:"email"`
	PhoneNo       *model.PhoneNumber `json:"phone_no"`
	Username      string             `json:"username"`
	EmailVerified bool               `json:"email_verified"`
	PhoneVerified bool               `json:"phone_verified"`
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

type GetUserInfoByIDOpts struct {
	ID primitive.ObjectID `json:"id" validate:"required"`
}
