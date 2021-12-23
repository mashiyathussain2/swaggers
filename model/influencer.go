package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	InfluencerColl               string = "influencer"
	InfluencerAccountRequestColl string = "influencer_request"
	CommissionLedgerColl         string = "commission_ledger"
	DebitRequestColl             string = "debit_request"
)

const (
	AcceptedStatus string = "accepted"
	InReviewStatus string = "in_review"
	RejectedStatus string = "rejected"
)

const (
	CreditTransaction string = "credit"
	DebitTransaction  string = "debit"
)

type Influencer struct {
	ID                primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string               `json:"name,omitempty" bson:"name,omitempty"`
	Username          string               `json:"username,omitempty" bson:"username,omitempty"`
	CoverImg          *IMG                 `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage      *IMG                 `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	SocialAccount     *SocialAccount       `json:"social_account,omitempty" bson:"social_account,omitempty"`
	ExternalLinks     []string             `json:"external_links,omitempty" bson:"external_links,omitempty"`
	Bio               string               `json:"bio,omitempty" bson:"bio,omitempty"`
	FollowersID       []primitive.ObjectID `json:"followers_id,omitempty" bson:"followers_id,omitempty"`
	FollowingID       []primitive.ObjectID `json:"following_id,omitempty" bson:"following_id,omitempty"`
	FollowersCount    uint                 `json:"followers_count" bson:"followers_count"`
	FollowingCount    uint                 `json:"following_count" bson:"following_count"`
	PayoutInformation *PayoutInformation   `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	Balance           uint                 `json:"balance,omitempty" bson:"balance,omitempty"`
	TotalCommission   uint                 `json:"total_commission,omitempty" bson:"total_commission,omitempty"`
	CreatedAt         time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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

type PayoutInformation struct {
	UPIID           string           `json:"upi_id,omitempty" bson:"upi_id,omitempty"`
	PanCard         string           `json:"pan_card,omitempty" bson:"pan_card,omitempty"`
	BankInformation *BankInformation `json:"bank_information,omitempty" bson:"bank_information,omitempty"`
}

type BankInformation struct {
	AccountHolderName string `json:"account_holder_name,omitempty" bson:"account_holder_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty" bson:"account_number,omitempty"`
	IFSCCode          string `json:"ifsc_code,omitempty" bson:"ifsc_code,omitempty"`
}

type Transaction struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID      primitive.ObjectID `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	Type              string             `json:"type,omitempty" bson:"type,omitempty"`
	ItemID            primitive.ObjectID `json:"item_id,omitempty" bson:"item_id,omitempty"`
	OrderID           primitive.ObjectID `json:"order_id,omitempty" bson:"order_id,omitempty"`
	OrderNo           string             `json:"order_no,omitempty" bson:"order_no,omitempty"`
	CatalogInfo       *CatalogInfo       `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	OrderValue        *Price             `json:"order_value,omitempty" bson:"order_value,omitempty"`
	CommissionValue   float64            `json:"commission_value,omitempty" bson:"commission_value,omitempty"`
	CreatedAt         time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Balance           float64            `json:"balance" bson:"balance"`
	DebitAmount       float64            `json:"debit_amount,omitempty" bson:"debit_amount,omitempty"`
	PayoutInformation *PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
}

type DebitRequest struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID      primitive.ObjectID `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	Amount            float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	Status            string             `json:"status,omitempty" bson:"status,omitempty"`
	GranteeID         primitive.ObjectID `json:"grantee_id,omitempty" bson:"grantee_id,omitempty"`
	PayoutInformation *PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	CreatedAt         time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
