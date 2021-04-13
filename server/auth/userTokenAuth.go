package auth

import (
	"encoding/base64"
	"encoding/json"

	"go-app/model"
	"go-app/server/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenAuthentication contains authentication related attributes and methods
type TokenAuthentication struct {
	Config *config.TokenAuthConfig
}

// NewTokenAuthentication returns new instance of TokenAuthentication
func NewTokenAuthentication(c *config.TokenAuthConfig) *TokenAuthentication {
	return &TokenAuthentication{Config: c}
}

// UserAuth contains encoded token info and user info
type UserAuth struct {
	UserClaim *UserClaim
	JWTToken  JWTToken
}

// UserClaim contains customer related info for jwt token
type UserClaim struct {
	ID            string             `json:"id"`
	KeeperUserID  string             `json:"keeper_user_id,omitempty"`
	CustomerID    string             `json:"customer_id,omitempty"`
	CartID        string             `json:"cart_id,omitempty"`
	Type          string             `json:"type"`
	Role          string             `json:"role"`
	FullName      string             `json:"full_name"`
	DOB           time.Time          `json:"dob,omitempty"`
	Email         string             `json:"email"`
	PhoneNo       *model.PhoneNumber `json:"phone_no,omitempty"`
	ProfileImage  *model.IMG         `json:"profile_image"`
	Gender        string             `json:"gender,omitempty"`
	EmailVerified bool               `json:"email_verified,omitempty"`
	PhoneVerified bool               `json:"phone_verified,omitempty"`
	CreatedVia    string             `json:"created_via,omitempty"`
	jwt.StandardClaims
}

// GetJWTToken return jwt.Token with claimInfo from user claim fields
func (uc *UserClaim) GetJWTToken() *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	return token
}

// ToJSON := converting struct to json
func (uc *UserClaim) ToJSON() string {
	json, _ := json.Marshal(uc)
	return string(json)
}

// IsAdmin if user is an admin user
func (uc *UserClaim) IsAdmin() bool {
	if uc.Role == model.AdminRole {
		return true
	}
	return false
}

// IsSudo if user is a keeper user
func (uc *UserClaim) IsSudo() bool {
	if uc.Type == model.KeeperType {
		return true
	}
	return false
}

// IsInternal if user is a keeper user
func (uc *UserClaim) IsInternal() bool {
	if uc.Type == model.InternalType {
		return true
	}
	return false
}

// SignToken sign and encodes jwt.Token as a string
func (t *TokenAuthentication) SignToken(claim Claim) (string, error) {
	userClaim := claim.(*UserClaim)
	if t.Config.JWTExpiresAt != 0 {
		expirationTime := time.Now().Add(time.Duration(t.Config.JWTExpiresAt) * time.Minute)
		userClaim.StandardClaims.ExpiresAt = expirationTime.Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, _ := token.SignedString([]byte(t.Config.JWTSignKey))
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), nil
}

// SignToken sign and encodes jwt.Token as a string
func (t *TokenAuthentication) SignKeeperToken(claim Claim) (string, error) {
	userClaim := claim.(*UserClaim)
	if t.Config.JWTExpiresAt != 0 {
		expirationTime := time.Now().Add(time.Duration(t.Config.JWTExpiresAt) * time.Minute)
		userClaim.StandardClaims.ExpiresAt = expirationTime.Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, _ := token.SignedString([]byte(t.Config.JWTSignKey))
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), nil
}

// VerifyToken first verifies the authenticity of the jwt token string and then parse the token string into struct
func (t *TokenAuthentication) VerifyToken(tokenString string) (Claim, error) {
	uc := UserClaim{}
	data, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(string(data), &uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Config.JWTSignKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return &uc, nil
}
