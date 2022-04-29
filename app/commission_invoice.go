package app

import (
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
	GenerateCommissionInvoice(debit_request_id primitive.ObjectID) error
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

//GenerateCommissionInvoice generates invoice based on debit_request collection id
func (ci *CommissionInvoiceImpl) GenerateCommissionInvoice(debit_request_id primitive.ObjectID) error {

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
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path": "$influencer_info",
		},
	}}
	cur, err := ci.DB.Collection(model.DebitRequestColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, unwindStage})
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
