package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	InfluencerColl               string = "influencer"
	InfluencerAccountRequestColl string = "influencer_request"
)

const (
	AcceptedStatus string = "accepted"
	InReviewStatus string = "in_review"
	RejectedStatus string = "rejected"
)

type Influencer struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string               `json:"name,omitempty" bson:"name,omitempty"`
	Username       string               `json:"username,omitempty" bson:"username,omitempty"`
	CoverImg       *IMG                 `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage   *IMG                 `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	SocialAccount  *SocialAccount       `json:"social_account,omitempty" bson:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty" bson:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty" bson:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id,omitempty" bson:"followers_id,omitempty"`
	FollowingID    []primitive.ObjectID `json:"following_id,omitempty" bson:"following_id,omitempty"`
	FollowersCount uint                 `json:"followers_count" bson:"followers_count"`
	FollowingCount uint                 `json:"following_count" bson:"following_count"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type InfluencerAccountRequest struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	// InfluencerID  primitive.ObjectID `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	Name          string         `json:"name,omitempty" bson:"name,omitempty"`
	Username      string         `json:"username,omitempty" bson:"username,omitempty"`
	ProfileImage  *IMG           `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	CoverImage    *IMG           `json:"cover_image,omitempty" bson:"cover_image,omitempty"`
	Bio           string         `json:"bio,omitempty" bson:"bio,omitempty"`
	Website       string         `json:"website,omitempty" bson:"website,omitempty"`
	SocialAccount *SocialAccount `json:"social_account,omitempty" bson:"social_account,omitempty"`
	IsActive      bool           `json:"is_active,omitempty" bson:"is_active,omitempty"`
	// IsGranted     *bool              `json:"is_granted,omitempty" bson:"is_granted,omitempty"`
	GranteeID primitive.ObjectID `json:"grantee_id,omitempty" bson:"grantee_id,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	GrantedAt time.Time          `json:"granted_at,omitempty" bson:"granted_at,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
}
