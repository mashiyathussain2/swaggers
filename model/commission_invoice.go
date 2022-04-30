package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Define name of the commission invoice collection
const (
	CommissionInvoiceColl string = "commission_invoice"
)

type CommissionInvoice struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DebitRequestID    primitive.ObjectID `json:"debit_request_id,omitempty" bson:"debit_request_id,omitempty"`
	InvoiceNo         string             `json:"invoice_no,omitempty" bson:"invoice_no,omitempty"`
	InfluencerID      primitive.ObjectID `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	InfluencerInfo    Influencer         `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	UserInfo          User               `json:"user_info,omitempty" bson:"user_info,omitempty"`
	Amount            uint               `json:"amount,omitempty" bson:"amount,omitempty"`
	PayoutInformation *PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	RequestDate       time.Time          `json:"request_date,omitempty" bson:"request_date,omitempty"`
	CreatedAt         time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type DebitRequestAllInfo struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID      primitive.ObjectID `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	Amount            float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	Status            string             `json:"status,omitempty" bson:"status,omitempty"`
	PayoutInformation *PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	InfluencerInfo    Influencer         `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	UserInfo          User               `json:"user_info,omitempty" bson:"user_info,omitempty"`
	GranteeID         primitive.ObjectID `json:"grantee_id,omitempty" bson:"grantee_id,omitempty"`
	CreatedAt         time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// InvoiceDate returns invoice creation date in IST timezone
func (i *CommissionInvoice) InvoiceDate() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return i.RequestDate.In(loc).Local().Format("02-Jan-2006")
}
