package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"go-app/model"
	"mime/multipart"
	"net/textproto"
	"strings"
	"time"

	"hash/fnv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommissionInvoice interface {
	CreateCommissionInvoice(debit_request_id primitive.ObjectID) error
	GetInvoicePDF(orderNo string) (*bytes.Buffer, string, error)
	SendCommissionInvoice(sc mongo.SessionContext, invoiceNo string) error
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
func (ci *CommissionInvoiceImpl) generateInvoiceNo(sc mongo.SessionContext, influencer_id primitive.ObjectID) (string, error) {

	filter := bson.M{"influencer_id": influencer_id}
	cnt, err := ci.DB.Collection(model.CommissionInvoiceColl).CountDocuments(sc, filter)
	if err != nil {
		return "", errors.Wrapf(err, "error counting commission invoices")
	}
	h := fnv.New32a()
	h.Write([]byte(influencer_id.Hex()))
	base := h.Sum32()
	s := fmt.Sprintf("%d-%d", base, cnt+1)
	return s, nil
}

func (ci *CommissionInvoiceImpl) validateGenerateInvoice(sc mongo.SessionContext, debitRequestID primitive.ObjectID) error {
	// Checking if invoice exists with provided order_no
	filter := bson.M{
		"debit_request_id": debitRequestID,
	}
	count, err := ci.DB.Collection(model.CommissionInvoiceColl).CountDocuments(sc, filter)
	if err != nil {
		return errors.Wrapf(err, "failed to check for invoice with debit_request_id:%s", debitRequestID)
	}
	if count != 0 {
		return errors.Errorf("invoice already generated for debit_request_id: %s", debitRequestID)
	}
	return nil
}

//CreateCommissionInvoice creates invoice based on debit_request collection id
func (ci *CommissionInvoiceImpl) CreateCommissionInvoice(debitRequestID primitive.ObjectID) error {

	ctx := context.TODO()
	var invoice model.CommissionInvoice
	var debitReqInfo []model.DebitRequestAllInfo
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": debitRequestID,
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
	addFieldsStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"user_id": bson.M{
				"$first": "$user_info._id",
			},
		},
	}}
	lookupStage3 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "customer",
			"localField":   "user_id",
			"foreignField": "user_id",
			"as":           "customer_info",
		},
	}}
	// unwindStage := bson.D{{
	// 	Key: "$unwind", Value: bson.M{
	// 		"path": "$influencer_info",
	// 	},
	// }}
	cur, err := ci.DB.Collection(model.DebitRequestColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, lookupStage2, addFieldsStage, lookupStage3})
	if err != nil {
		return errors.Wrapf(err, "error getting debit request")
	}
	if err := cur.All(ctx, &debitReqInfo); err != nil {
		return errors.Wrap(err, "error decoding debit request")
	}
	if len(debitReqInfo) != 1 {
		return errors.New("error debit request incorrect")
	}

	// creating session for atomic updates
	session, err := ci.DB.Client().StartSession()
	if err != nil {
		return errors.Wrap(err, "failed to start session")
	}
	// Closing session at the end for function execution
	defer session.EndSession(ctx)

	// staring a new transaction
	if err := session.StartTransaction(); err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	if err = mongo.WithSession(context.TODO(), session, func(sc mongo.SessionContext) error {

		//validating request:
		err := ci.validateGenerateInvoice(sc, debitRequestID)
		if err != nil {
			return errors.Wrapf(err, "error validating invoice generation")
		}
		// get unique invoice no based on influencerid
		invoiceNo, err := ci.generateInvoiceNo(sc, debitReqInfo[0].InfluencerID)
		if err != nil {
			return errors.Wrapf(err, "error generating invoice no")
		}
		invoice = model.CommissionInvoice{
			InvoiceNo:         invoiceNo,
			DebitRequestID:    debitRequestID,
			InfluencerID:      debitReqInfo[0].InfluencerID,
			InfluencerInfo:    debitReqInfo[0].InfluencerInfo[0],
			CustomerInfo:      debitReqInfo[0].CustomerInfo[0],
			UserInfo:          debitReqInfo[0].UserInfo[0],
			Amount:            uint(debitReqInfo[0].Amount),
			PayoutInformation: debitReqInfo[0].PayoutInformation,
			RequestDate:       debitReqInfo[0].CreatedAt,
			CreatedAt:         time.Now(),
		}
		_, err = ci.DB.Collection(model.CommissionInvoiceColl).InsertOne(sc, invoice)
		if err != nil {
			return errors.Wrapf(err, "error generating invoice")
		}

		if err := session.CommitTransaction(sc); err != nil {
			return errors.Wrapf(err, "failed to commit transaction")
		}

		return nil
	}); err != nil {
		ci.Logger.Err(err).Msgf("failed to generate invoice for debit_request_id: %s", debitRequestID.Hex())
		return err
	}
	err = ci.SendCommissionInvoice(nil, invoice.InvoiceNo)
	if err != nil {
		return errors.Wrapf(err, "error sending invoice")
	}

	return nil
}

func (ci *CommissionInvoiceImpl) GetCIbyNo(sc mongo.SessionContext, invoiceNo string) (*model.CommissionInvoice, error) {
	var invoice model.CommissionInvoice
	if sc != nil {
		err := ci.DB.Collection(model.CommissionInvoiceColl).FindOne(sc, bson.M{"invoice_no": invoiceNo}).Decode(&invoice)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting invoice")
		}
	}
	ctx := context.TODO()
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

func (ci *CommissionInvoiceImpl) GetInvoicePDF(invoiceNo string) (*bytes.Buffer, string, error) {
	invoice, err := ci.GetCIbyNo(nil, invoiceNo)
	if err != nil {
		return nil, "", err
	}
	resp, fileName, err := ci.generateCommissionInvoicePDF(invoice)
	if err != nil {
		return nil, "", err
	}
	// ci.SendCommissionInvoice(invoiceNo)
	return resp, fileName, nil
}

func (ci *CommissionInvoiceImpl) GetPreInvoicePDF(debitRequestID primitive.ObjectID) (*bytes.Buffer, string, error) {

	ctx := context.TODO()
	var debitReqInfo []model.DebitRequestAllInfo
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": debitRequestID,
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
	addFieldsStage := bson.D{{
		Key: "$addFields", Value: bson.M{
			"user_id": bson.M{
				"$first": "$user_info._id",
			},
		},
	}}
	lookupStage3 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "customer",
			"localField":   "user_id",
			"foreignField": "user_id",
			"as":           "customer_info",
		},
	}}
	cur, err := ci.DB.Collection(model.DebitRequestColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, lookupStage2, addFieldsStage, lookupStage3})
	if err != nil {
		return nil, "", errors.Wrapf(err, "error getting debit request")
	}
	if err := cur.All(ctx, &debitReqInfo); err != nil {
		return nil, "", errors.Wrap(err, "error decoding debit request")
	}
	if len(debitReqInfo) != 1 {
		return nil, "", errors.New("error debit request incorrect")
	}

	// get unique invoice no based on influencerid
	invoiceNo, err := ci.generateInvoiceNo(nil, debitReqInfo[0].InfluencerID)
	if err != nil {
		return nil, "", errors.Wrapf(err, "error generating invoice no")
	}
	invoice := model.CommissionInvoice{
		InvoiceNo:         invoiceNo,
		DebitRequestID:    debitRequestID,
		InfluencerID:      debitReqInfo[0].InfluencerID,
		InfluencerInfo:    debitReqInfo[0].InfluencerInfo[0],
		UserInfo:          debitReqInfo[0].UserInfo[0],
		Amount:            uint(debitReqInfo[0].Amount),
		PayoutInformation: debitReqInfo[0].PayoutInformation,
		RequestDate:       debitReqInfo[0].CreatedAt,
		CreatedAt:         time.Now(),
	}
	resp, fileName, err := ci.generateCommissionInvoicePDF(&invoice)
	if err != nil {
		return nil, "", err
	}
	return resp, fileName, nil
}

func (ci *CommissionInvoiceImpl) commissionInvoiceMailTemplate(invoice *model.CommissionInvoice) string {
	t := fmt.Sprintf(`
	Hey %s <br>
	-- <br>
	Congratulations! Your commission request for Amount ₹%d, is accepted and will be transferred within 2 business days ⚡⚡. <br>
	<br>
	PFA the invoice for the same. <br>
	Keep sharing your links , and keep earning more! 💰💰<br>
	<br>
	<br>
	Regards, <br>
	Team Hypd <br>
	`, invoice.CustomerInfo.FullName, invoice.Amount)
	return t
}

func (ci *CommissionInvoiceImpl) prepareCommissionInvoiceEmail(message, attachmentFilename string, destination, cc []string, file []byte) (*ses.SendRawEmailInput, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// preparing email main header
	h := make(textproto.MIMEHeader)
	h.Set("From", ci.App.Config.HypdFinanceEmail)
	for _, i := range destination {
		h.Add("To", i)
	}

	cc = append(cc, ci.App.Config.HypdFinanceEmail)

	for _, i := range cc {
		h.Add("Cc", i)
	}
	h.Set("Subject", "Commission Request Accepted")
	h.Set("Content-Language", "en-IN")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// body:
	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Type", "text/html; charset=us-ascii")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	// file attachment:
	fn := attachmentFilename
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment;filename="+fn)
	h.Set("Content-Type", "application/pdf; name=\""+fn+"\"")
	h.Set("Content-Transfer-Encoding", "base64")
	sEnc := base64.StdEncoding.EncodeToString([]byte(file))
	part, err = writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(sEnc))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Strip boundary line before header (doesn't work with it present)
	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		return nil, fmt.Errorf("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]

	raw := ses.RawMessage{
		Data: []byte(s),
	}

	var dest []*string
	for _, i := range destination {
		dest = append(dest, aws.String(i))
	}
	for _, c := range cc {
		dest = append(dest, aws.String(c))
	}
	input := &ses.SendRawEmailInput{
		Destinations: dest,
		Source:       aws.String(ci.App.Config.HypdFinanceEmail),
		RawMessage:   &raw,
	}

	return input, nil
}

func (ci *CommissionInvoiceImpl) SendCommissionInvoice(sc mongo.SessionContext, invoiceNo string) error {
	invoice, err := ci.GetCIbyNo(sc, invoiceNo)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to get invoice by invoice no: %s", invoiceNo)
		return err
	}
	if invoice == nil {
		ci.Logger.Err(err).Msgf("invoice not found by invoice no: %s", invoiceNo)
		return err
	}
	file, fn, err := ci.generateCommissionInvoicePDF(invoice)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to generate Commission Invoice PDF: %s", invoiceNo)
		return err
	}
	attachmentFilename := fn
	message := ci.commissionInvoiceMailTemplate(invoice)
	destination := invoice.UserInfo.Email
	cc := []string{}

	email, err := ci.prepareCommissionInvoiceEmail(message, attachmentFilename, []string{destination}, cc, file.Bytes())
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to prepare email to send to creator for invoice no: %s", invoiceNo)
		return err
	}
	resp, err := ci.App.SES.SendRawEmail(email)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to send email to creator for invoice no: %s", invoiceNo)
		return err
	}
	ci.Logger.Debug().Interface("resp", resp).Msgf("sent email to creator for invoice no: %s", invoiceNo)
	return nil
}
