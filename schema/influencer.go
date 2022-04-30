package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EditInfluencerOptsOpts create fields and validations required to create an new instance of influencer
type CreateInfluencerOpts struct {
	Name          string             `json:"name" validate:"required"`
	Username      string             `json:"username" validate:"required"`
	Bio           string             `json:"bio"`
	CoverImg      *Img               `json:"cover_img" validate:"required"`
	ProfileImage  *Img               `json:"profile_image" validate:"required"`
	ExternalLinks []string           `json:"external_links" validate:"required,min=1,dive,min=6"`
	SocialAccount *SocialAccountOpts `json:"social_account"`
}

// CreateInfluencerResp contains fields to be returned in response to create influencer
type CreateInfluencerResp struct {
	ID            primitive.ObjectID   `json:"id"`
	Name          string               `json:"name"`
	Username      string               `json:"username"`
	Bio           string               `json:"bio"`
	CoverImg      *model.IMG           `json:"cover_img"`
	ProfileImage  *model.IMG           `json:"profile_image"`
	ExternalLinks []string             `json:"external_links"`
	SocialAccount *model.SocialAccount `json:"social_account"`
	CreatedAt     time.Time            `json:"created_at"`
}

// EditInfluencerOpts contains fields and validations required to edit existing influencer
type EditInfluencerOpts struct {
	ID            primitive.ObjectID `json:"id" validate:"required"`
	Name          string             `json:"name"`
	Username      string             `json:"username"`
	Bio           string             `json:"bio"`
	CoverImg      *Img               `json:"cover_img"`
	ProfileImage  *Img               `json:"profile_image"`
	ExternalLinks []string           `json:"external_links"`
	SocialAccount *SocialAccountOpts `json:"social_account"`
}

// EditInfluencerResp contains fields to be returned in response to edit influencer

// swagger:model EditInfluencerResp
type EditInfluencerResp struct {
	ID            primitive.ObjectID   `json:"id"`
	Name          string               `json:"name"`
	Username      string               `json:"username"`
	Bio           string               `json:"bio"`
	CoverImg      *model.IMG           `json:"cover_img"`
	ProfileImage  *model.IMG           `json:"profile_image"`
	ExternalLinks []string             `json:"external_links"`
	SocialAccount *model.SocialAccount `json:"social_account"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// GetInfluencersByIDOpts contains fields and validations required to get multiple influencer by matching id
type GetInfluencersByIDOpts struct {
	IDs []primitive.ObjectID `json:"id" validate:"required,min=1"`
}

type GetInfluencersByNameOpts struct {
	Name string `json:"name" validate:"required,min=3"`
}

// GetInfluencerResp contains fields to be returned for get influencer function
type GetInfluencerResp struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string               `json:"name,omitempty" bson:"name,omitempty"`
	Username       string               `json:"username,omitempty" bson:"username,omitempty"`
	CoverImg       *model.IMG           `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage   *model.IMG           `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	SocialAccount  *model.SocialAccount `json:"social_account,omitempty" bson:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty" bson:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty" bson:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id,omitempty" bson:"followers_id,omitempty"`
	FollowingID    []primitive.ObjectID `json:"following_id,omitempty" bson:"following_id,omitempty"`
	FollowersCount uint                 `json:"followers_count" bson:"followers_count"`
	FollowingCount uint                 `json:"following_count" bson:"following_count"`
}

// swagger:model AddInfluencerFollowerOpts
type AddInfluencerFollowerOpts struct {
	// swagger:strfmt ObjectID
	// required:true
	InfluencerID primitive.ObjectID `json:"id" validate:"required"`
	// swagger:strfmt ObjectID
	// required:true
	CustomerID primitive.ObjectID `json:"customer_id" validate:"required"`
}

type InfluencerKafkaMessage struct {
	ID             primitive.ObjectID   `json:"_id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Username       string               `json:"username,omitempty"`
	CoverImg       *model.IMG           `json:"cover_img,omitempty"`
	ProfileImage   *model.IMG           `json:"profile_image,omitempty"`
	SocialAccount  *model.SocialAccount `json:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id,omitempty"`
	FollowingID    []primitive.ObjectID `json:"following_id,omitempty"`
	FollowersCount uint                 `json:"followers_count"`
	FollowingCount uint                 `json:"following_count"`
	CreatedAt      time.Time            `json:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty"`
}
type InfluencerFullKafkaMessageOpts struct {
	ID             primitive.ObjectID   `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Username       string               `json:"username,omitempty"`
	CoverImg       *model.IMG           `json:"cover_img,omitempty"`
	ProfileImage   *model.IMG           `json:"profile_image,omitempty"`
	SocialAccount  *model.SocialAccount `json:"social_account,omitempty"`
	ExternalLinks  []string             `json:"external_links,omitempty"`
	Bio            string               `json:"bio,omitempty"`
	FollowersID    []primitive.ObjectID `json:"followers_id"`
	FollowingID    []primitive.ObjectID `json:"following_id"`
	FollowersCount uint                 `json:"followers_count"`
	FollowingCount uint                 `json:"following_count"`
	CreatedAt      time.Time            `json:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty"`
}

type LinkUserAccountOpts struct {
	RequestID    primitive.ObjectID `json:"request_id" validate:"required"`
	InfluencerID primitive.ObjectID `json:"influencer_id" validate:"required"`
	UserID       primitive.ObjectID `json:"user_id" validate:"required"`
}

// swagger:model InfluencerAccountRequestOpts
type InfluencerAccountRequestOpts struct {
	// swagger:strfmt ObjectID
	// required:true
	UserID primitive.ObjectID `json:"user_id" validate:"required"`
	// swagger:strfmt ObjectID
	// required:true
	CustomerID primitive.ObjectID `json:"customer_id" validate:"required"`
	// InfluencerID  primitive.ObjectID `json:"influencer_id" validate:"required"`
	FullName      string             `json:"full_name" validate:"required"`
	Username      string             `json:"username,omitempty"`
	ProfileImage  Img                `json:"profile_image" validate:"required"`
	CoverImage    Img                `json:"cover_image" validate:"required"`
	Bio           string             `json:"bio" validate:"required"`
	Website       string             `json:"website"`
	SocialAccount *SocialAccountOpts `json:"social_account"`
	Source        interface{}
}

type UpdateInfluencerAccountRequestStatusOpts struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	Grant     *bool              `json:"grant" validate:"required"`
	GranteeID primitive.ObjectID
}

type InfluencerAccountRequestInfluencerInfo struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	CoverImg     *model.IMG         `json:"cover_img,omitempty" bson:"cover_img,omitempty"`
	ProfileImage *model.IMG         `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
}

type InfluencerAccountRequestCustomerInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"full_name" bson:"full_name"`
	Gender   string             `json:"gender" bson:"gender"`
	DOB      time.Time          `json:"dob" bson:"dob"`
}

type InfluencerAccountRequestUserInfo struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNo *model.PhoneNumber `json:"phone_no" bson:"phone_no"`
	Email   string             `json:"email" bson:"email"`
}

type InfluencerAccountRequestResp struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// InfluencerID   primitive.ObjectID                      `json:"influencer_id" bson:"influencer_id"`
	Name          string                                `json:"name" bson:"name"`
	ProfileImage  *model.IMG                            `json:"profile_image" bson:"profile_image"`
	CoverImage    *model.IMG                            `json:"cover_image" bson:"cover_image"`
	Bio           string                                `json:"bio" bson:"bio"`
	Website       string                                `json:"website" bson:"website"`
	SocialAccount *model.SocialAccount                  `json:"social_account" bson:"social_account"`
	CustomerInfo  *InfluencerAccountRequestCustomerInfo `json:"customer_info" bson:"customer_info"`
	UserInfo      *InfluencerAccountRequestUserInfo     `json:"user_info" bson:"user_info"`
	IsActive      bool                                  `json:"is_active" bson:"is_active"`
	GranteeID     primitive.ObjectID                    `json:"grantee_id" bson:"grantee_id"`
	CreatedAt     time.Time                             `json:"created_at" bson:"created_at"`
	GrantedAt     time.Time                             `json:"granted_at" bson:"granted_at"`
	Status        string                                `json:"status,omitempty" bson:"status,omitempty"`
	Source        *map[string]string                    `json:"source,omitempty" bson:"source,omitempty"`
}

// EditInfluencerAppOpts contains fields and validations required to edit existing influencer

// swagger:model EditInfluencerAppOpts
type EditInfluencerAppOpts struct {
	// swagger:strfmt ObjectID
	ID primitive.ObjectID `json:"id" validate:"required"`
	// Name          string             `json:"name"`
	Username string `json:"username"`
	// Bio           string             `json:"bio"`
	// CoverImg      *Img               `json:"cover_img"`
	// ProfileImage  *Img               `json:"profile_image"`
	// ExternalLinks []string           `json:"external_links"`
	// SocialAccount *SocialAccountOpts `json:"social_account"`
	PayoutInformation *PayoutInformationOpts `json:"payout_information"`
}

type PayoutInformationOpts struct {
	UPIID           string                 `json:"upi_id"`
	PanCard         string                 `json:"pan_card"`
	BankInformation *model.BankInformation `json:"bank_information"`
}

type BankInformationOpts struct {
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	IFSCCode          string `json:"ifsc_code"`
}

type ProcessInsertOrderOpts struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderID      string             `json:"order_id,omitempty" bson:"order_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	BrandID      primitive.ObjectID `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	InfluencerID primitive.ObjectID `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
}

// swagger:model CommissionDebitRequest
type CommissionDebitRequest struct {
	// swagger:strfmt ObjectID
	ID                primitive.ObjectID    `json:"id" validate:"required"`
	Amount            uint                  `json:"amount" validate:"required,gte=1000"`
	PayoutInformation PayoutInformationOpts `json:"payout_information" validate:"required"`
}

type UpdateCommissionDebitRequest struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	Status    string             `json:"status" validate:"required"`
	GranteeID primitive.ObjectID `json:"grantee_id" validate:"required"`
	// PayoutInformation *PayoutInformationOpts `json:"payout_information" validate:"required"`
}

type GetDebitRequestResponse struct {
	ID                    primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Status                string                `json:"status,omitempty" bson:"status,omitempty"`
	Amount                uint                  `json:"amount,omitempty" bson:"amount,omitempty"`
	InfluencerInfo        GetInfluencerResp     `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	Email                 string                `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNo               string                `json:"phone_no,omitempty" bson:"phone_no,omitempty"`
	PayoutInformationOpts PayoutInformationResp `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
}

// swagger:model GetInfluencerDashboardResp
type GetInfluencerDashboardResp struct {
	MonthlyData []MonthlyData             `json:"monthly_data,omitempty" bson:"monthly_data,omitempty"`
	Ledger      []GetInfluencerLedgerResp `json:"ledger,omitempty" bson:"ledger,omitempty"`
	OverallData OverallData               `json:"overall_data,omitempty" bson:"overall_data,omitempty"`
}
type MonthlyData struct {
	Month uint `json:"month,omitempty" bson:"_id,omitempty"`
	Count uint `json:"count,omitempty" bson:"count,omitempty"`
}

type OverallData struct {
	Revenue         uint `json:"revenue,omitempty" bson:"revenue,omitempty"`
	TotalCommission uint `json:"total_commission,omitempty" bson:"total_commission,omitempty"`
}

// swagger:model GetInfluencerLedgerOpts
type GetInfluencerLedgerOpts struct {
	// swagger:strfmt ObjectID
	ID        primitive.ObjectID `json:"id" validate:"required"`
	Page      int                `json:"page"`
	Type      string             `json:"type" validate:"required"`
	StartDate *time.Time         `json:"start_date"`
	EndDate   *time.Time         `json:"end_date"`
}

// swagger:model GetInfluencerDashboardOpts
type GetInfluencerDashboardOpts struct {
	// swagger:strfmt ObjectID
	ID        primitive.ObjectID `json:"id" validate:"required"`
	StartDate *time.Time         `json:"start_date"`
	EndDate   *time.Time         `json:"end_date"`
}

// swagger:model GetInfluencerLedgerResp
type GetInfluencerLedgerResp struct {
	Date       string       `json:"date,omitempty" bson:"_id,omitempty"`
	Ledger     []LedgerResp `json:"ledger,omitempty" bson:"ledger,omitempty"`
	Commission uint         `json:"commission,omitempty" bson:"commission,omitempty"`
	Revenue    uint         `json:"revenue,omitempty" bson:"revenue,omitempty"`
}

type LedgerResp struct {
	ID                primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	Type              string                   `json:"type,omitempty" bson:"type,omitempty"`
	OrderNo           string                   `json:"order_no,omitempty" bson:"order_no,omitempty"`
	CatalogInfo       *model.CatalogInfo       `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	OrderValue        *model.Price             `json:"order_value,omitempty" bson:"order_value,omitempty"`
	CommissionValue   float64                  `json:"commission_value,omitempty" bson:"commission_value,omitempty"`
	CreatedAt         time.Time                `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Balance           float64                  `json:"balance,omitempty" bson:"balance,omitempty"`
	DebitAmount       float64                  `json:"debit_amount,omitempty" bson:"debit_amount,omitempty"`
	PayoutInformation *model.PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
}

// swagger:model GetPayoutInfoResp
type GetPayoutInfoResp struct {
	// swagger:strfmt ObjectID
	ID                primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	Balance           float64                  `json:"balance,omitempty" bson:"balance,omitempty"`
	PayoutInformation *model.PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
}

type PayoutInformationResp struct {
	UPIID           string                 `json:"upi_id,omitempty" bson:"upi_id,omitempty"`
	PanCard         string                 `json:"pan_card,omitempty" bson:"pan_card,omitempty"`
	BankInformation *model.BankInformation `json:"bank_information,omitempty" bson:"bank_information,omitempty"`
}

// swagger:model GetCommissionAndRevenueOpts
type GetCommissionAndRevenueOpts struct {
	// swagger:strfmt ObjectID
	ID        primitive.ObjectID `json:"id" validate:"required"`
	StartDate *time.Time         `json:"start_date" validate:"required" `
	EndDate   *time.Time         `json:"end_date" validate:"required"`
}

// swagger:model GetCommissionAndRevenueResp
type GetCommissionAndRevenueResp struct {
	Commission uint `json:"commission,omitempty" bson:"commission,omitempty"`
	Revenue    uint `json:"revenue,omitempty" bson:"revenue,omitempty"`
	Balance    uint `json:"balance,omitempty" bson:"balance,omitempty"`
}

// EditInfluencerAppV2Opts contains fields and validations required to edit existing influencer

// swagger:model EditInfluencerAppV2Opts
type EditInfluencerAppV2Opts struct {
	// swagger:strfmt ObjectID
	ID                primitive.ObjectID     `json:"id" validate:"required"`
	Name              string                 `json:"name,omitempty"`
	Username          string                 `json:"username,omitempty"`
	Bio               string                 `json:"bio,omitempty"`
	CoverImg          *Img                   `json:"cover_img,omitempty"`
	ProfileImage      *Img                   `json:"profile_image,omitempty"`
	ExternalLinks     []string               `json:"external_links,omitempty"`
	SocialAccount     *SocialAccountOpts     `json:"social_account,omitempty"`
	PayoutInformation *PayoutInformationOpts `json:"payout_information,omitempty"`
}

// swagger:model InfluencerAccountRequestV2Opts
type InfluencerAccountRequestV2Opts struct {
	// swagger:strfmt ObjectID
	UserID primitive.ObjectID `json:"user_id" validate:"required"`
	// swagger:strfmt ObjectID
	CustomerID      primitive.ObjectID `json:"customer_id" validate:"required"`
	FullName        string             `json:"full_name" validate:"required"`
	Username        string             `json:"username,omitempty"`
	Email           string             `json:"email,omitempty"`
	Phone           *PhoneNoOpts       `json:"phone,omitempty"`
	ProfileImage    Img                `json:"profile_image" validate:"required"`
	CoverImage      Img                `json:"cover_image" validate:"required"`
	SocialAccount   *SocialAccountOpts `json:"social_account"`
	AreaOfExpertise string             `json:"area_of_expertise"`
	Source          interface{}
}
