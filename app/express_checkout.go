package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExpressCheckout contains methods for ExpressCheckout service functionality
type ExpressCheckout interface {
	ExpressCheckout(*schema.ExpressCheckoutOpts) (*schema.OrderInfo, error)
	ExpressCheckoutComplete(*schema.ExpressCheckoutOpts) (*schema.OrderInfo, error)
}

// ExpressCheckoutImpl implements ExpressCheckout interface methods
type ExpressCheckoutImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// ExpressCheckoutImplOpts contains args required to create
type ExpressCheckoutImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitExpressCheckout returns new instance of ExpressCheckout implementation
func InitExpressCheckout(opts *ExpressCheckoutImplOpts) ExpressCheckout {
	ec := ExpressCheckoutImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ec
}

func (ec *ExpressCheckoutImpl) ExpressCheckout(opts *schema.ExpressCheckoutOpts) (*schema.OrderInfo, error) {

	orderItem := schema.OrderItem{
		CatalogID: opts.Items[0].CatalogID,
		VariantID: opts.Items[0].VariantID,
		Quantity:  uint(opts.Items[0].Quantity),
	}
	var variant schema.OrderVariant

	var s model.GetAllCatalogInfoResp

	url := ec.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + opts.Items[0].CatalogID.Hex()
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch catlog data")
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ec.Logger.Err(err).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrap(err, "failed to get catalog info")
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ec.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		ec.Logger.Err(errors.New("success false from catalog")).Str("body", string(body)).Msg("got success false response from catalog")
		return nil, errors.New("got success false response from catalog")
	}

	//checking if variant exist or not
	found := false
	for _, v := range s.Payload.Variants {
		if v.ID == opts.Items[0].VariantID {
			found = true
			variant = schema.OrderVariant{
				ID:        v.ID,
				SKU:       v.SKU,
				Attribute: v.Attribute,
			}
			break
		}
	}
	if !found {
		return nil, errors.Errorf("variant with id %s not found", opts.Items[0].VariantID.Hex())
	}

	//calculate discount if available
	discount := uint(0)
	discountInfo := model.DiscountInfo{}
	var dp *model.Price
	if s.Payload.DiscountInfo != nil {
		fmt.Println(s.Payload.DiscountInfo)
		for _, d := range s.Payload.DiscountInfo.VariantsID {
			if d == opts.Items[0].VariantID {
				switch s.Payload.DiscountInfo.Type {
				case model.FlatOffType:
					discount = s.Payload.DiscountInfo.Value
					dp = model.SetINRPrice(s.Payload.RetailPrice.Value - float32(discount))
				case model.PercentOffType:
					discount = uint(float64((s.Payload.DiscountInfo.Value * uint(s.Payload.RetailPrice.Value)) / 100.0))
					if discount > s.Payload.DiscountInfo.MaxValue && s.Payload.DiscountInfo.MaxValue > 0 {
						discount = s.Payload.DiscountInfo.MaxValue
					}
					discountInfo.MaxValue = s.Payload.DiscountInfo.MaxValue
					dp = model.SetINRPrice(s.Payload.RetailPrice.Value - float32(discount))
				}
				discountInfo.Value = discount
				discountInfo.ID = s.Payload.DiscountInfo.ID
				discountInfo.Type = s.Payload.DiscountInfo.Type
				discountInfo.Value = discount
				orderItem.DiscountID = s.Payload.DiscountInfo.ID
				orderItem.DiscountInfo = &discountInfo
				orderItem.DiscountedPrice = dp

			}
		}
	}
	orderItem.BasePrice = &s.Payload.BasePrice
	orderItem.RetailPrice = &s.Payload.RetailPrice
	orderItem.CatalogInfo = schema.OrderCatalogInfo{
		ID:      opts.Items[0].CatalogID,
		BrandID: s.Payload.BrandID,

		Name:    s.Payload.Name,
		Variant: variant,
		FeaturedImage: schema.Img{
			SRC: s.Payload.FeaturedImage.SRC,
		},

		VariantType: s.Payload.VariantType,
		HSNCode:     s.Payload.HSNCode,

		TransferPrice: s.Payload.TransferPrice,
		ETA:           s.Payload.ETA,
	}
	orderOpts := []schema.OrderOpts{
		{
			UserID:          opts.UserID,
			BrandID:         s.Payload.BrandID,
			ShippingAddress: opts.Address,
			BillingAddress:  opts.Address,
			Source:          opts.Source,
			OrderItems:      []schema.OrderItem{orderItem},
		},
	}

	//Create Order
	coURL := ec.App.Config.HypdApiConfig.OrderApi + "/api/order"

	var orderResp schema.OrderResp
	reqBody, err := json.Marshal(orderOpts)
	fmt.Println(string(reqBody))
	if err != nil {
		ec.Logger.Err(err).Interface("orderOpts", orderOpts).Msgf("failed to prepare request json to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	resp, err = http.Post(coURL, "application/json", bytes.NewBuffer(reqBody))
	//Handle Error
	if err != nil {
		ec.Logger.Err(err).RawJSON("responseBody", reqBody).Msgf("failed to send request to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	defer resp.Body.Close()
	//Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		ec.Logger.Err(err).RawJSON("reqBody", reqBody).Msgf("failed to read response from api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		ec.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !orderResp.Success {
		ec.Logger.Err(errors.New("success false from order")).Str("body", string(body)).Msg("got success false response from order")
		return nil, errors.New("got success false response from order")
	}

	return &orderResp.Payload, nil
}

func (ec *ExpressCheckoutImpl) ExpressCheckoutComplete(opts *schema.ExpressCheckoutOpts) (*schema.OrderInfo, error) {

	var orderOpts []schema.OrderOpts
	// var orderItems []schema.OrderItem

	oiBrandMap := make(map[primitive.ObjectID][]schema.OrderItem)

	for _, item := range opts.Items {
		orderItem := schema.OrderItem{
			CatalogID: item.CatalogID,
			VariantID: item.VariantID,
			Quantity:  uint(item.Quantity),
		}
		var variant schema.OrderVariant

		var s model.GetAllCatalogInfoResp

		url := ec.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + item.CatalogID.Hex()
		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to fetch catlog data")
		}
		defer resp.Body.Close()

		//Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ec.Logger.Err(err).Msgf("failed to read response from api %s", url)
			return nil, errors.Wrap(err, "failed to get catalog info")
		}
		if err := json.Unmarshal(body, &s); err != nil {
			ec.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
			return nil, errors.Wrap(err, "failed to decode body into struct")
		}
		if !s.Success {
			ec.Logger.Err(errors.New("success false from catalog")).Str("body", string(body)).Msg("got success false response from catalog")
			return nil, errors.New("got success false response from catalog")
		}
		//checking if variant exist or not
		found := false
		for _, v := range s.Payload.Variants {
			if v.ID == item.VariantID {
				found = true
				variant = schema.OrderVariant{
					ID:        v.ID,
					SKU:       v.SKU,
					Attribute: v.Attribute,
				}
				break
			}
		}
		if !found {
			return nil, errors.Errorf("variant with id %s not found", item.VariantID.Hex())
		}
		//calculate discount if available
		discount := uint(0)
		discountInfo := model.DiscountInfo{}
		var dp *model.Price
		if s.Payload.DiscountInfo != nil {
			fmt.Println(s.Payload.DiscountInfo)
			for _, d := range s.Payload.DiscountInfo.VariantsID {
				if d == item.VariantID {
					switch s.Payload.DiscountInfo.Type {
					case model.FlatOffType:
						discount = s.Payload.DiscountInfo.Value
						dp = model.SetINRPrice(s.Payload.RetailPrice.Value - float32(discount))
					case model.PercentOffType:
						discount = uint(float64((s.Payload.DiscountInfo.Value * uint(s.Payload.RetailPrice.Value)) / 100.0))
						if discount > s.Payload.DiscountInfo.MaxValue && s.Payload.DiscountInfo.MaxValue > 0 {
							discount = s.Payload.DiscountInfo.MaxValue
						}
						discountInfo.MaxValue = s.Payload.DiscountInfo.MaxValue
						dp = model.SetINRPrice(s.Payload.RetailPrice.Value - float32(discount))
					}
					discountInfo.Value = discount
					discountInfo.ID = s.Payload.DiscountInfo.ID
					discountInfo.Type = s.Payload.DiscountInfo.Type
					discountInfo.Value = discount
					orderItem.DiscountID = s.Payload.DiscountInfo.ID
					orderItem.DiscountInfo = &discountInfo
					orderItem.DiscountedPrice = dp

				}
			}
		}
		orderItem.BasePrice = &s.Payload.BasePrice
		orderItem.RetailPrice = &s.Payload.RetailPrice
		orderItem.CatalogInfo = schema.OrderCatalogInfo{
			ID:      item.CatalogID,
			BrandID: s.Payload.BrandID,
			Name:    s.Payload.Name,
			Variant: variant,
			FeaturedImage: schema.Img{
				SRC: s.Payload.FeaturedImage.SRC,
			},

			VariantType: s.Payload.VariantType,
			HSNCode:     s.Payload.HSNCode,

			TransferPrice: s.Payload.TransferPrice,
			ETA:           s.Payload.ETA,
		}
		// orderItems = append(orderItems, orderItem)
		oiBrandMap[orderItem.CatalogInfo.BrandID] = append(oiBrandMap[orderItem.CatalogInfo.BrandID], orderItem)
	}

	for brand, oi := range oiBrandMap {
		orderOpts = append(orderOpts, schema.OrderOpts{
			UserID:          opts.UserID,
			BrandID:         brand,
			ShippingAddress: opts.Address,
			BillingAddress:  opts.Address,
			Source:          opts.Source,
			OrderItems:      oi,
		})
	}
	//Create Order
	coURL := ec.App.Config.HypdApiConfig.OrderApi + "/api/order"

	var orderResp schema.OrderResp
	reqBody, err := json.Marshal(orderOpts)
	fmt.Println(string(reqBody))
	if err != nil {
		ec.Logger.Err(err).Interface("orderOpts", orderOpts).Msgf("failed to prepare request json to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	resp, err := http.Post(coURL, "application/json", bytes.NewBuffer(reqBody))
	//Handle Error
	if err != nil {
		ec.Logger.Err(err).RawJSON("responseBody", reqBody).Msgf("failed to send request to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ec.Logger.Err(err).RawJSON("reqBody", reqBody).Msgf("failed to read response from api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		ec.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !orderResp.Success {
		ec.Logger.Err(errors.New("success false from order")).Str("body", string(body)).Msg("got success false response from order")
		return nil, errors.New("got success false response from order")
	}

	return &orderResp.Payload, nil

}
