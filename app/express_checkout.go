package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExpressCheckout contains methods for ExpressCheckout service functionality
type ExpressCheckout interface {
	ExpressCheckout(*schema.ExpressCheckoutOpts) (*schema.OrderInfo, error)
	ExpressCheckoutComplete(*schema.ExpressCheckoutOpts, string, string) (*schema.OrderInfo, error)
	ExpressCheckoutWeb(*schema.ExpressCheckoutWebOpts, string) (*schema.OrderInfo, error)
	ExpressCheckoutRTO(opts *schema.ExpressCheckoutWebOpts, userName string, userAgent, ipAddress, email string) (interface{}, error)
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
		Source:    opts.Items[0].Source,
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
	orderOpts := []schema.OrderItemOpts{
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

func (ec *ExpressCheckoutImpl) ExpressCheckoutComplete(opts *schema.ExpressCheckoutOpts, userName, platform string) (*schema.OrderInfo, error) {

	var orderOpts []schema.OrderItemOpts
	// var orderItems []schema.OrderItem
	grandTotal := 0
	displayName := strings.ToLower(opts.Address.DisplayName)
	if displayName == "home" || displayName == "other" || displayName == "work" || displayName == "" {
		opts.Address.DisplayName = userName
	}

	oiBrandMap := make(map[primitive.ObjectID][]schema.OrderItem)

	for _, item := range opts.Items {
		orderItem := schema.OrderItem{
			CatalogID: item.CatalogID,
			VariantID: item.VariantID,
			Quantity:  uint(item.Quantity),
			Source:    opts.Items[0].Source,
		}
		var variant schema.OrderVariant

		var s model.GetAllCatalogInfoResp

		url := ec.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + item.CatalogID.Hex()
		// resp, err := http.Get(url)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "unable to fetch catlog data")
		// }
		// defer resp.Body.Close()
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate request to get catalog & variant")
		}
		req.Header.Add("Authorization", ec.App.Config.HypdApiConfig.Token)
		resp, err := client.Do(req)
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

			TransferPrice:  s.Payload.TransferPrice,
			ETA:            s.Payload.ETA,
			CommissionRate: s.Payload.CommissionRate,
		}
		if !orderItem.DiscountID.IsZero() {
			grandTotal += int(orderItem.DiscountedPrice.Value)
		} else {
			grandTotal += (int(orderItem.RetailPrice.Value))

		}
		orderItem.Tax = s.Payload.Tax
		// orderItems = append(orderItems, orderItem)
		oiBrandMap[orderItem.CatalogInfo.BrandID] = append(oiBrandMap[orderItem.CatalogInfo.BrandID], orderItem)
	}
	var couponOrderOpts schema.CouponOrderOpts
	if opts.Coupon != "" {
		appliedValue := model.SetINRPrice(0)
		coupon, err := ec.App.Cart.GetCoupon(opts.Coupon)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting coupon")
		}
		if coupon.Status != "active" {
			return nil, errors.Errorf("coupon is not active")
		}
		if coupon.Type == model.FlatOffType {
			appliedValue = model.SetINRPrice(float32(coupon.Value))
		} else if coupon.Type == model.PercentOffType {
			av := (grandTotal * coupon.Value) / 100
			if coupon.MaxDiscount != nil {
				if av > int(coupon.MaxDiscount.Value) {
					av = int(coupon.MaxDiscount.Value)
				}
			}
			appliedValue = model.SetINRPrice(float32(av))
		}
		couponOrderOpts = schema.CouponOrderOpts{
			ID:           coupon.ID,
			Code:         coupon.Code,
			AppliedValue: appliedValue,
		}
	}

	for brand, oi := range oiBrandMap {
		orderItem := schema.OrderItemOpts{
			UserID:          opts.UserID,
			BrandID:         brand,
			ShippingAddress: opts.Address,
			BillingAddress:  opts.Address,
			Source:          opts.Source,
			OrderItems:      oi,
			SourceID:        &opts.SourceID,
			Platform:        platform,
			CartType:        model.ExpressCheckout,
		}
		if opts.Coupon != "" {
			orderItem.Coupon = &couponOrderOpts
		}
		orderOpts = append(orderOpts, orderItem)
	}

	//Create Order
	coURL := ec.App.Config.HypdApiConfig.OrderApi + "/api/order"

	var orderResp schema.OrderResp
	reqBody, err := json.Marshal(orderOpts)
	if err != nil {
		ec.Logger.Err(err).Interface("orderOpts", orderOpts).Msgf("failed to prepare request json to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	// resp, err := http.Post(coURL, "application/json", bytes.NewBuffer(reqBody))
	// //Handle Error
	// if err != nil {
	// 	ec.Logger.Err(err).RawJSON("responseBody", reqBody).Msgf("failed to send request to api %s", coURL)
	// 	return nil, errors.Wrap(err, "failed to get order info")
	// }
	// defer resp.Body.Close()

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, coURL, bytes.NewBuffer(reqBody))
	if err != nil {
		ec.Logger.Err(err).Interface("reqBody", reqBody).Msgf("failed to generate order")
		return nil, errors.Wrap(err, "failed to generate request to generate order")
	}
	req.Header.Add("Authorization", ec.App.Config.HypdApiConfig.Token)
	resp, err := client.Do(req)
	if err != nil {
		ec.Logger.Err(err).Interface("reqBody", reqBody).Msgf("unable to fetch order info")
		return nil, errors.Wrapf(err, "unable to fetch order info")
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

func (ec *ExpressCheckoutImpl) ExpressCheckoutWeb(opts *schema.ExpressCheckoutWebOpts, userName string) (*schema.OrderInfo, error) {

	var orderOpts []schema.OrderItemOpts
	// var orderItems []schema.OrderItem
	grandTotal := 0
	isCouponApplied := false
	var coupon *model.Coupon
	var err error
	if opts.Coupon != "" {
		isCouponApplied = true
		coupon, err = ec.App.Cart.GetCoupon(opts.Coupon)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting coupon")
		}
	}
	displayName := strings.ToLower(opts.Address.DisplayName)
	if displayName == "home" || displayName == "other" || displayName == "work" || displayName == "" {
		opts.Address.DisplayName = userName
	}

	oiBrandMap := make(map[primitive.ObjectID][]schema.OrderItem)

	for _, item := range opts.Items {
		orderItem := schema.OrderItem{
			CatalogID: item.CatalogID,
			VariantID: item.VariantID,
			Quantity:  uint(item.Quantity),
			Source:    opts.Items[0].Source,
		}
		var variant schema.OrderVariant

		var s model.GetAllCatalogInfoResp

		url := ec.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + item.CatalogID.Hex()
		// resp, err := http.Get(url)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "unable to fetch catlog data")
		// }
		// defer resp.Body.Close()
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate request to get catalog & variant")
		}
		req.Header.Add("Authorization", ec.App.Config.HypdApiConfig.Token)
		resp, err := client.Do(req)
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
			if isCouponApplied {
				toAdd := false
				fmt.Println("apply coupon on ", coupon.ApplicableON.Name)
				switch coupon.ApplicableON.Name {
				case "brand":
					if orderItem.CatalogInfo.BrandID == coupon.ApplicableON.IDs[0] {
						toAdd = true
					}
				case "influencer":
					if orderItem.Source.ID == coupon.ApplicableON.IDs[0].Hex() {
						toAdd = true
					}
				case "cart":
					toAdd = true
				}
				if toAdd {
					grandTotal += int(dp.Value) * int(orderItem.Quantity)
				}
			} else {
				grandTotal += int(dp.Value) * int(orderItem.Quantity)
			}
		} else {
			if isCouponApplied {
				toAdd := false
				switch coupon.ApplicableON.Name {
				case "brand":
					if orderItem.CatalogInfo.BrandID == coupon.ApplicableON.IDs[0] {
						toAdd = true
					}
				case "influencer":
					if orderItem.Source.ID == coupon.ApplicableON.IDs[0].Hex() {
						toAdd = true
					}
				case "cart":
					toAdd = true
				}
				if toAdd {
					grandTotal += int(s.Payload.RetailPrice.Value) * int(orderItem.Quantity)
				}
			} else {
				grandTotal += int(s.Payload.RetailPrice.Value) * int(orderItem.Quantity)
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

			TransferPrice:  s.Payload.TransferPrice,
			ETA:            s.Payload.ETA,
			CommissionRate: s.Payload.CommissionRate,
		}
		orderItem.Tax = s.Payload.Tax
		// orderItems = append(orderItems, orderItem)
		oiBrandMap[orderItem.CatalogInfo.BrandID] = append(oiBrandMap[orderItem.CatalogInfo.BrandID], orderItem)
	}

	var couponOrderOpts schema.CouponOrderOpts
	if isCouponApplied {
		appliedValue := model.SetINRPrice(0)
		if coupon.Status != "active" {
			return nil, errors.Errorf("coupon is not active")
		}
		if coupon.Type == model.FlatOffType {
			appliedValue = model.SetINRPrice(float32(coupon.Value))
		} else if coupon.Type == model.PercentOffType {
			av := (grandTotal * coupon.Value) / 100
			if coupon.MaxDiscount != nil {
				if av > int(coupon.MaxDiscount.Value) {
					av = int(coupon.MaxDiscount.Value)
				}
			}
			appliedValue = model.SetINRPrice(float32(av))
		}
		couponOrderOpts = schema.CouponOrderOpts{
			ID:           coupon.ID,
			Code:         coupon.Code,
			AppliedValue: appliedValue,
			ApplicableON: coupon.ApplicableON,
		}
		fmt.Println("coupon order opts ", couponOrderOpts)
	}

	for brand, oi := range oiBrandMap {
		orderItem := schema.OrderItemOpts{
			UserID:          opts.UserID,
			BrandID:         brand,
			ShippingAddress: opts.Address,
			BillingAddress:  opts.Address,
			OrderItems:      oi,
			Platform:        "web",
			CartType:        model.ExpressCheckout,
			IsWeb:           true,
			Source:          opts.Source,
			SourceID:        &opts.SourceID,
			IsCOD:           opts.IsCOD,
			RequestID:       opts.RequestID,
		}
		if opts.Coupon != "" {
			orderItem.Coupon = &couponOrderOpts
		}
		orderOpts = append(orderOpts, orderItem)
	}
	//Create Order
	coURL := ec.App.Config.HypdApiConfig.OrderApi + "/api/order"

	var orderResp schema.OrderResp
	reqBody, err := json.Marshal(orderOpts)
	if err != nil {
		ec.Logger.Err(err).Interface("orderOpts", orderOpts).Msgf("failed to prepare request json to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	// resp, err := http.Post(coURL, "application/json", bytes.NewBuffer(reqBody))
	// //Handle Error
	// if err != nil {
	// 	ec.Logger.Err(err).RawJSON("responseBody", reqBody).Msgf("failed to send request to api %s", coURL)
	// 	return nil, errors.Wrap(err, "failed to get order info")
	// }
	// defer resp.Body.Close()

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, coURL, bytes.NewBuffer(reqBody))
	if err != nil {
		ec.Logger.Err(err).Interface("reqBody", reqBody).Msgf("failed to generate order")
		return nil, errors.Wrap(err, "failed to generate request to generate order")
	}
	req.Header.Add("Authorization", ec.App.Config.HypdApiConfig.Token)
	resp, err := client.Do(req)
	if err != nil {
		ec.Logger.Err(err).Interface("reqBody", reqBody).Msgf("unable to fetch order info")
		return nil, errors.Wrapf(err, "unable to fetch order info")
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

func (ec *ExpressCheckoutImpl) ExpressCheckoutRTO(opts *schema.ExpressCheckoutWebOpts, userName string, userAgent, ipAddress, email string) (interface{}, error) {

	// totalPrice := 0
	totalItems := 0
	grandTotal := 0
	displayName := strings.ToLower(opts.Address.DisplayName)
	if displayName == "home" || displayName == "other" || displayName == "work" || displayName == "" {
		opts.Address.DisplayName = userName
	}
	var lineItems []schema.GoKwikLineItems

	// var orderItems []schema.OrderItem
	for _, item := range opts.Items {
		orderItem := schema.OrderItem{
			CatalogID: item.CatalogID,
			VariantID: item.VariantID,
			Quantity:  uint(item.Quantity),
			Source:    opts.Items[0].Source,
		}
		totalItems += int(item.Quantity)

		var variant schema.OrderVariant
		var s model.GetAllCatalogInfoResp

		url := ec.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + item.CatalogID.Hex()
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate request to get catalog & variant")
		}
		req.Header.Add("Authorization", ec.App.Config.HypdApiConfig.Token)
		resp, err := client.Do(req)
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
			grandTotal += int(dp.Value) * int(orderItem.Quantity)
		} else {
			grandTotal += int(s.Payload.RetailPrice.Value) * int(orderItem.Quantity)
		}

		// orderItem.BasePrice = &s.Payload.BasePrice
		// orderItem.RetailPrice = &s.Payload.RetailPrice
		// orderItem.CatalogInfo = schema.OrderCatalogInfo{
		// 	ID:      item.CatalogID,
		// 	BrandID: s.Payload.BrandID,
		// 	Name:    s.Payload.Name,
		// 	Variant: variant,
		// 	FeaturedImage: schema.Img{
		// 		SRC: s.Payload.FeaturedImage.SRC,
		// 	},

		// 	VariantType: s.Payload.VariantType,
		// 	HSNCode:     s.Payload.HSNCode,

		// 	TransferPrice:  s.Payload.TransferPrice,
		// 	ETA:            s.Payload.ETA,
		// 	CommissionRate: s.Payload.CommissionRate,
		// }
		// orderItem.Tax = s.Payload.Tax
		// orderItems = append(orderItems, orderItem)
		lineItems = append(lineItems, schema.GoKwikLineItems{
			Sku:                 variant.SKU,
			Price:               float64(s.Payload.RetailPrice.Value),
			Quantity:            int(item.Quantity),
			Total:               int(item.Quantity) * int(s.Payload.RetailPrice.Value),
			ProductThumbnailURL: s.Payload.FeaturedImage.SRC,
			ProductURL:          ec.App.Config.HypdApiConfig.CatalogURL + "/product?id=" + item.CatalogID.Hex(),
		})
		// orderItems = append(orderItems, orderItem)
	}
	nameParts := strings.Split(opts.Address.DisplayName, " ")
	address := schema.GoKwikShippingAddress{
		FirstName: nameParts[0],
		LastName:  strings.Join(nameParts[1:], " "),
		Address1:  opts.Address.Line1,
		Address2:  opts.Address.Line2,
		City:      opts.Address.City,
		State:     opts.Address.State.Name,
		Postcode:  opts.Address.PostalCode,
		Phone:     opts.Address.ContactNumber.Number,
		Email:     email,
	}
	baddress := schema.GoKwikBillingAddress{
		Address1: opts.Address.Line1,
		Address2: opts.Address.Line2,
		City:     opts.Address.City,
		State:    opts.Address.State.Name,
		Postcode: opts.Address.PostalCode,
	}
	eOpts := schema.CheckCODEligiblityOpts{
		Customer: schema.GoKwikCustomer{},
		Order: schema.GoKwikOrder{
			OrderDate:              time.Now(),
			Subtotal:               grandTotal,
			TotalLineItems:         len(opts.Items),
			TotalLineItemsQuantity: totalItems,
			TotalDiscount:          0,
			Total:                  grandTotal,
			PromoCode:              "",
			LineItems:              lineItems,
			ShippingAddress:        address,
			BillingAddress:         baddress,
			Session: schema.GoKwikSession{
				Source:            "organic",
				CustomerUserAgent: userAgent,
				CustomerIP:        ipAddress,
			},
		},
	}

	coURL := ec.App.Config.GoKwikConfig.RTOApi

	var rtoResp interface{}
	reqBody, err := json.Marshal(eOpts)

	fmt.Println(4)
	fmt.Println(string(reqBody))
	if err != nil {
		ec.Logger.Err(err).Interface("eOpts", eOpts).Msgf("failed to prepare request json to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get cart info")
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, coURL, bytes.NewBuffer(reqBody))
	if err != nil {
		ec.Logger.Err(err).Interface("eOpts", eOpts).Msgf("failed to create request to check eligiblity %s", coURL)
		return nil, errors.Wrap(err, "failed to create request to check eligiblity")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appid", ec.App.Config.GoKwikConfig.AppID)
	req.Header.Add("appsecret", ec.App.Config.GoKwikConfig.AppSecret)

	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		ec.Logger.Err(err).RawJSON("responseBody", reqBody).Msgf("failed to send request to to check eligiblity %s", coURL)
		return nil, errors.Wrap(err, "failed to send request to check eligiblity")
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(5)

	if err != nil {
		ec.Logger.Err(err).RawJSON("reqBody", reqBody).Msgf("failed to read response from gokwik api %s", coURL)
		return nil, errors.Wrap(err, "failed to read gokwik info")
	}
	if err := json.Unmarshal(body, &rtoResp); err != nil {
		ec.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}

	return &rtoResp, nil
}
