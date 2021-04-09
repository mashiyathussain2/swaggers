package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name in DB
const (
	UserColl       string = "user"
	KeeperUserColl string = "keeper_user"
)

// list of supported user types
const (
	CustomerType   string = "customer"
	InfluencerType string = "influencer"
	BrandType      string = "brand"
	KeeperType     string = "keeper"
)

// list of supported roles
const (
	UserRole       string = "user"
	AdminRole      string = "admin"
	SuperadminRole string = "superadmin"
)

// LoginOTP stores contains otp for login functionality
type LoginOTP struct {
	Type      string    `json:"type,omitempty" bson:"type,omitempty"`
	OTP       string    `json:"otp,omitempty" bson:"otp,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// PhoneLoginOTPType contains list of supported otp login type
const (
	PhoneLoginOTPType string = "phone_no"
)

// list of CreatedViaOptions
const (
	CreateViaMobile    string = "mobile"
	CreatedViaGoogle   string = "google"
	CreatedViaFacebook string = "facebook"
)

// User represents an entity-user
/*
	Entity user can be of three types
		-> Customer (app user)
		-> Influencer
		-> Brand
*/
type User struct {
	ID                    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type                  string             `json:"type,omitempty" bson:"type,omitempty"`
	Role                  string             `json:"role,omitempty" bson:"role,omitempty"`
	Email                 string             `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNo               *PhoneNumber       `json:"phone_no,omitempty" bson:"phone_no,omitempty"`
	LoginOTP              *LoginOTP          `json:"login_otp,omitempty" bson:"login_otp,omitempty"`
	Username              string             `json:"username,omitempty" bson:"username,omitempty"`
	Password              string             `json:"password,omitempty" bson:"password,omitempty"`
	PasswordResetCode     string             `json:"password_reset_code,omitempty" bson:"password_reset_code,omitempty"`
	EmailVerificationCode string             `json:"email_verification_code,omitempty" bson:"email_verification_code,omitempty"`
	PhoneVerificationCode string             `json:"phone_verification_code,omitempty" bson:"phone_verification_code,omitempty"`
	EmailVerifiedAt       time.Time          `json:"email_verified_at,omitempty" bson:"email_verified_at,omitempty"`
	PhoneVerifiedAt       time.Time          `json:"phone_verified_at,omitempty" bson:"phone_verified_at,omitempty"`
	CreatedVia            string             `json:"created_via,omitempty" bson:"created_via,omitempty"`
	CreatedAt             time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt             time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type KeeperUser struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	FullName     string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	ProfileImage *IMG               `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
