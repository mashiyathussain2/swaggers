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
	ID                primitive.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	DebitRequestID    primitive.ObjectID       `json:"debit_request_id,omitempty" bson:"debit_request_id,omitempty"`
	InvoiceNo         string                   `json:"invoice_no,omitempty" bson:"invoice_no,omitempty"`
	InfluencerID      primitive.ObjectID       `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	InfluencerInfo    InfluencerInfoForInvoice `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	CustomerInfo      CustomerInfoForInvoice   `json:"customer_info,omitempty" bson:"customer_info,omitempty"`
	UserInfo          UserInfoForInvoice       `json:"user_info,omitempty" bson:"user_info,omitempty"`
	Amount            uint                     `json:"amount,omitempty" bson:"amount,omitempty"`
	PayoutInformation *PayoutInformation       `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	RequestDate       time.Time                `json:"request_date,omitempty" bson:"request_date,omitempty"`
	CreatedAt         time.Time                `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type DebitRequestAllInfo struct {
	ID                primitive.ObjectID         `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID      primitive.ObjectID         `json:"influencer_id,omitempty" bson:"influencer_id,omitempty"`
	Amount            float64                    `json:"amount,omitempty" bson:"amount,omitempty"`
	Status            string                     `json:"status,omitempty" bson:"status,omitempty"`
	PayoutInformation *PayoutInformation         `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	InfluencerInfo    []InfluencerInfoForInvoice `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	UserInfo          []UserInfoForInvoice       `json:"user_info,omitempty" bson:"user_info,omitempty"`
	CustomerInfo      []CustomerInfoForInvoice   `json:"customer_info,omitempty" bson:"customer_info,omitempty"`
	GranteeID         primitive.ObjectID         `json:"grantee_id,omitempty" bson:"grantee_id,omitempty"`
	CreatedAt         time.Time                  `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         time.Time                  `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UserInfoForInvoice struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email   string             `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNo *PhoneNumber       `json:"phone_no,omitempty" bson:"phone_no,omitempty"`
}

type InfluencerInfoForInvoice struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name,omitempty" bson:"name,omitempty"`
	PayoutInformation *PayoutInformation `json:"payout_information,omitempty" bson:"payout_information,omitempty"`
	Balance           uint               `json:"balance,omitempty" bson:"balance,omitempty"`
	TotalCommission   uint               `json:"total_commission,omitempty" bson:"total_commission,omitempty"`
}

type CustomerInfoForInvoice struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Address  []Address          `json:"address,omitempty" bson:"address,omitempty"`
}

// InvoiceDate returns invoice creation date in IST timezone
func (i *CommissionInvoice) InvoiceDate() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return i.RequestDate.In(loc).Local().Format("02-Jan-2006")
}

func (i *CommissionInvoice) IsBankTransfer() bool {
	return i.PayoutInformation.BankInformation != nil
}
