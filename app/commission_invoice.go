package app

import (
	"bytes"
	"context"
	"fmt"
	"go-app/model"
	"time"

	"hash/fnv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommissionInvoice interface {
	CreateCommissionInvoice(debit_request_id primitive.ObjectID) error
	GetInvoicePDF(userID primitive.ObjectID, orderNo string) (*bytes.Buffer, string, error)
}

type CommissionInvoiceImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

type CommissionInvoiceOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

func InitCommissionInvoice(opts *CommissionInvoiceOpts) CommissionInvoice {
	ci := CommissionInvoiceImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ci
}

// generateInvoiceNo generates invoice no for influencer
func (ci *CommissionInvoiceImpl) generateInvoiceNo(influencer_id primitive.ObjectID) (string, error) {

	filter := bson.M{"influencer_id": influencer_id}
	cnt, err := ci.DB.Collection(model.CommissionInvoiceColl).CountDocuments(context.TODO(), filter)
	if err != nil {
		return "", errors.Wrapf(err, "error counting commission invoices")
	}
	h := fnv.New32a()
	h.Write([]byte(influencer_id.Hex()))
	base := h.Sum32()
	s := fmt.Sprintf("%d-%d", base, cnt+1)
	return s, nil
}

func (ci *CommissionInvoiceImpl) validateGenerateInvoice(sc mongo.SessionContext, debit_request_id primitive.ObjectID) error {
	// Checking if invoice exists with provided order_no
	filter := bson.M{
		"debit_request_id": debit_request_id,
	}
	count, err := ci.DB.Collection(model.CommissionInvoiceColl).CountDocuments(sc, filter)
	if err != nil {
		return errors.Wrapf(err, "failed to check for invoice with debit_request_id:%s", debit_request_id)
	}
	if count != 0 {
		return errors.Errorf("invoice already generated for debit_request_id: %s", debit_request_id)
	}
	return nil
}

//CreateCommissionInvoice creates invoice based on debit_request collection id
func (ci *CommissionInvoiceImpl) CreateCommissionInvoice(debit_request_id primitive.ObjectID) error {

	ctx := context.TODO()
	var debitReqInfo []model.DebitRequestAllInfo
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": debit_request_id,
		},
	}}
	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "influencer",
			"localField":   "influencer_id",
			"foreignField": "_id",
			"as":           "influencer_info",
		},
	}}
	lookupStage2 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "user",
			"localField":   "influencer_id",
			"foreignField": "influencer_id",
			"as":           "user_info",
		},
	}}
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path": "$influencer_info",
		},
	}}
	cur, err := ci.DB.Collection(model.DebitRequestColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, lookupStage2, unwindStage})
	if err != nil {
		return errors.Wrapf(err, "error getting debit request")
	}
	if err := cur.All(ctx, &debitReqInfo); err != nil {
		return errors.Wrap(err, "error decoding debit request")
	}
	if len(debitReqInfo) != 1 {
		return errors.New("error debit request incorrect")
	}

	// get unique invoice no based on influencerid
	invoiceNo, err := ci.generateInvoiceNo(debitReqInfo[0].InfluencerID)
	if err != nil {
		return errors.Wrapf(err, "error generating invoice no")
	}
	invoice := model.CommissionInvoice{
		InvoiceNo:         invoiceNo,
		DebitRequestID:    debit_request_id,
		InfluencerID:      debitReqInfo[0].InfluencerID,
		InfluencerInfo:    debitReqInfo[0].InfluencerInfo,
		UserInfo:          debitReqInfo[0].UserInfo,
		Amount:            uint(debitReqInfo[0].Amount),
		PayoutInformation: debitReqInfo[0].PayoutInformation,
		RequestDate:       debitReqInfo[0].CreatedAt,
		CreatedAt:         time.Now(),
	}
	_, err = ci.DB.Collection(model.CommissionInvoiceColl).InsertOne(ctx, invoice)
	if err != nil {
		return errors.Wrapf(err, "error generating invoice")
	}
	return nil
}

func (ci *CommissionInvoiceImpl) GetCIbyNo(invoiceNo string) (*model.CommissionInvoice, error) {
	ctx := context.TODO()
	var invoice model.CommissionInvoice
	err := ci.DB.Collection(model.CommissionInvoiceColl).FindOne(ctx, bson.M{"invoice_no": invoiceNo}).Decode(&invoice)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting invoice")
	}
	return &invoice, nil
}
func (ci *CommissionInvoiceImpl) generateCommissionInvoicePDF(invoice *model.CommissionInvoice) (*bytes.Buffer, string, error) {
	body, err := ParseTemplate(ci.App.Config.PDFConfig.CommissionInvoiceTemplatePath, invoice)
	if err != nil {
		ci.Logger.Err(err).Msg("failed to prepare pdf")
		return nil, "", err
	}
	buff, err := GeneratePDF(body)
	if err != nil {
		ci.Logger.Err(err).Msg("failed to generate pdf")
		return nil, "", err
	}
	return buff, fmt.Sprintf("%s.pdf", invoice.InvoiceNo), nil
}

func (ci *CommissionInvoiceImpl) GetInvoicePDF(userID primitive.ObjectID, orderNo string) (*bytes.Buffer, string, error) {
	invoice, err := ci.GetCIbyNo(orderNo)
	if err != nil {
		return nil, "", err
	}
	resp, fileName, err := ci.generateCommissionInvoicePDF(invoice)
	if err != nil {
		return nil, "", err
	}
	return resp, fileName, nil
}
