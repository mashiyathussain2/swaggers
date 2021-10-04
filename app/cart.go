package app

import (
	"bytes"
	"context"
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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Cart contains methods for Cart service functionality
type Cart interface {
	CreateCart(primitive.ObjectID) (primitive.ObjectID, error)
	AddToCart(*schema.AddToCartOpts) (*model.Cart, error)
	UpdateItemQty(*schema.UpdateItemQtyOpts) (*model.Cart, error)
	GetCartInfo(primitive.ObjectID) (*schema.GetCartInfoResp, error)
	SetCartAddress(*schema.AddressOpts) error
	CheckoutCart(primitive.ObjectID, string, string, string) (*schema.OrderInfo, error)
	ClearCart(primitive.ObjectID) error

	AddDiscountInCartItems(*schema.DiscountInCartItemsOpts)
	RemoveDiscountInCartItems(*schema.DiscountInCartItemsOpts)
	UpdateInventoryStatus(*schema.InventoryUpdateOpts)
	UpdateCatalogInfo(id primitive.ObjectID)
	ApplyCoupon(primitive.ObjectID, *schema.ApplyCouponOpts) error
	RemoveCoupon(primitive.ObjectID) error
	UpdateInventoryStatusInsideCatalogInfo(opts *schema.InventoryUpdateOpts)
}

// CartImpl implements Cart interface methods
type CartImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// CartImplOpts contains args required to create
type CartImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitCart returns new instance of Cart implementation
func InitCart(opts *CartImplOpts) Cart {
	ci := CartImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ci
}

func (ci *CartImpl) CreateCart(id primitive.ObjectID) (primitive.ObjectID, error) {
	ctx := context.TODO()
	cart := model.Cart{
		UserID:    id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	cartID, err := ci.DB.Collection(model.CartColl).InsertOne(ctx, cart)
	if err != nil {
		return primitive.NilObjectID, errors.Wrapf(err, "unable to create cart for user with id: %s", id)
	}

	filter := bson.M{
		"user_id": id,
	}
	update := bson.M{
		"$set": bson.M{
			"cart_id": cartID.InsertedID.(primitive.ObjectID),
		},
	}
	if _, err := ci.DB.Collection(model.CustomerColl).UpdateOne(ctx, filter, update); err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "failed to link cart and customer")
	}
	return cartID.InsertedID.(primitive.ObjectID), nil
}

func (ci *CartImpl) AddToCart(opts *schema.AddToCartOpts) (*model.Cart, error) {

	ctx := context.TODO()
	var s model.GetAllCatalogInfoResp

	url := ci.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + opts.CatalogID.Hex()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to request to get catalog info")
	}
	req.Header.Add("Authorization", ci.App.Config.HypdApiConfig.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch catlog data")
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrap(err, "failed to get catalog info")
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		ci.Logger.Err(errors.New("success false from catalog")).Str("body", string(body)).Msg("got success false response from catalog")
		return nil, errors.New("got success false response from catalog")
	}

	//checking if variant exist or not
	found := false
	for _, v := range s.Payload.Variants {
		if v.ID == opts.VariantID {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.Errorf("variant with id %s not found", opts.VariantID.Hex())
	}

	//checking if item already in cart
	findFilter := bson.M{
		"user_id":          opts.ID,
		"items.catalog_id": opts.CatalogID,
		"items.variant_id": opts.VariantID,
	}
	var cartMongo model.Cart
	mongoErr := ci.DB.Collection(model.CartColl).FindOne(ctx, findFilter).Decode(&cartMongo)
	if mongoErr != nil {
		if mongoErr != mongo.ErrNilDocument && mongoErr != mongo.ErrNoDocuments {
			return nil, errors.Wrapf(mongoErr, "unable to check cart for catalog")
		}
	}
	if cartMongo.UserID == opts.ID {
		return nil, errors.Errorf("item already in cart")
	}

	//calculate discount if available
	discount := uint(0)
	discountInfo := model.DiscountInfo{}

	if s.Payload.DiscountInfo != nil {
		for _, d := range s.Payload.DiscountInfo.VariantsID {
			if d == opts.VariantID {
				switch s.Payload.DiscountInfo.Type {
				case model.FlatOffType:
					discount = s.Payload.DiscountInfo.Value
				case model.PercentOffType:
					discount = uint(float64((s.Payload.DiscountInfo.Value * uint(s.Payload.RetailPrice.Value)) / 100.0))
					if discount > s.Payload.DiscountInfo.MaxValue && s.Payload.DiscountInfo.MaxValue > 0 {
						discount = s.Payload.DiscountInfo.MaxValue
					}
					discountInfo.MaxValue = s.Payload.DiscountInfo.MaxValue
				}
				discountInfo.Value = discount
				discountInfo.ID = s.Payload.DiscountInfo.ID
				discountInfo.Type = s.Payload.DiscountInfo.Type
			}
		}
	}
	item := model.Item{
		ID:        primitive.NewObjectID(),
		CatalogID: opts.CatalogID,
		BrandID:   s.Payload.BrandID,
		VariantID: opts.VariantID,
		BrandInfo: s.Payload.BrandInfo,
		CatalogInfo: &model.CatalogInfo{
			ID:            s.Payload.ID,
			BrandID:       s.Payload.BrandID,
			Name:          s.Payload.Name,
			FeaturedImage: s.Payload.FeaturedImage,

			VariantType:   s.Payload.VariantType,
			Variants:      s.Payload.Variants,
			HSNCode:       s.Payload.HSNCode,
			TransferPrice: s.Payload.TransferPrice,

			ETA:          s.Payload.ETA,
			Status:       s.Payload.Status,
			Tax:          s.Payload.Tax,
			DiscountInfo: s.Payload.DiscountInfo,
		},
		BasePrice:     &s.Payload.BasePrice,
		RetailPrice:   &s.Payload.RetailPrice,
		TransferPrice: &s.Payload.TransferPrice,
		Quantity:      opts.Quantity,
	}
	if s.Payload.DiscountInfo != nil {
		item.DiscountID = s.Payload.DiscountInfo.ID
		item.DiscountInfo = &discountInfo
	}

	updateQuery := bson.M{
		"$push": bson.M{
			"items": item,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}

	filter := bson.M{
		"user_id": opts.ID,
	}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var cart model.Cart
	err = ci.DB.Collection(model.CartColl).FindOneAndUpdate(ctx, filter, updateQuery, qOpts).Decode(&cart)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to update cart")
	}
	return &cart, nil
}

//UpdateItemQty function increases, decreases or removes the item from cart based on input qty
func (ci *CartImpl) UpdateItemQty(opts *schema.UpdateItemQtyOpts) (*model.Cart, error) {

	ctx := context.TODO()

	var s model.GetCatalogVariant

	url := ci.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + opts.CatalogID.Hex() + "/variant/" + opts.VariantID.Hex()
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request to get catalog and variant")
	}
	req.Header.Add("Authorization", ci.App.Config.HypdApiConfig.Token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to fetch catlog data")
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to read response from api %s", url)
		return nil, errors.Wrap(err, "failed to get brandinfo")
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !s.Success {
		ci.Logger.Err(errors.New("success false from entity")).Str("body", string(body)).Msg("got success false response from entity")
		return nil, errors.New("got success false response from catalog")
	}

	filterQuery := bson.M{
		"user_id":          opts.ID,
		"items.catalog_id": opts.CatalogID,
		"items.variant_id": opts.VariantID,
	}
	var updateQuery bson.M
	if opts.Quantity == 0 {
		updateQuery = bson.M{
			"$pull": bson.M{
				"items": bson.M{
					"catalog_id": opts.CatalogID,
					"variant_id": opts.VariantID,
				},
			},
			"$set": bson.M{
				"updated_at": time.Now().UTC(),
			},
		}
	} else {

		discount := uint(0)
		if s.Payload.DiscountInfo != nil {
			if s.Payload.DiscountInfo.Type == model.FlatOffType {
				discount = uint(s.Payload.DiscountInfo.Value)
			} else {
				discount = (s.Payload.DiscountInfo.Value * uint(s.Payload.RetailPrice.Value)) / 100
				if discount > s.Payload.DiscountInfo.MaxValue {
					discount = s.Payload.DiscountInfo.MaxValue
				}
			}
		}
		incPrice := s.Payload.RetailPrice
		incGrandTotal := uint(s.Payload.RetailPrice.Value) - discount
		updateQuery = bson.M{
			"$inc": bson.M{
				"items.$.quantity":     opts.Quantity,
				"total_price.value":    opts.Quantity * int(incPrice.Value),
				"grand_total.value":    opts.Quantity * int(incGrandTotal),
				"total_discount.value": opts.Quantity * int(discount),
			},
		}
	}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var cart model.Cart
	err = ci.DB.Collection(model.CartColl).FindOneAndUpdate(ctx, filterQuery, updateQuery, qOpts).Decode(&cart)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to update cart")
	}
	return &cart, nil
}

//GetCartInfo function returns cart info
func (ci *CartImpl) GetCartInfo(id primitive.ObjectID) (*schema.GetCartInfoResp, error) {
	ctx := context.TODO()
	var cart schema.GetCartInfoResp
	err := ci.DB.Collection(model.CartColl).FindOne(ctx, bson.M{"user_id": id}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Errorf("cart with id :%s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "unable to query for cart")
	}
	tp := uint(0)
	rp := uint(0)
	td := uint(0)
	gt := uint(0)
	for i, cartItem := range cart.Items {
		tp = tp + uint(cartItem.BasePrice.Value)*cartItem.Quantity
		rp = rp + uint(cartItem.RetailPrice.Value)*cartItem.Quantity
		// gt = gt + uint(cartItem.RetailPrice.Value)*cartItem.Quantity
		td = td + uint(cartItem.BasePrice.Value-cartItem.RetailPrice.Value)*cartItem.Quantity
		if cartItem.CatalogInfo.DiscountInfo != nil {
			applied := false
			for _, v := range cartItem.CatalogInfo.DiscountInfo.VariantsID {
				if v == cartItem.VariantID {
					var dp *model.Price
					cart.Items[i].DiscountInfo = &model.DiscountInfo{
						ID:    cartItem.CatalogInfo.DiscountInfo.ID,
						Type:  cartItem.CatalogInfo.DiscountInfo.Type,
						Value: cartItem.CatalogInfo.DiscountInfo.Value,
					}
					switch cartItem.CatalogInfo.DiscountInfo.Type {
					case model.FlatOffType:
						dp = model.SetINRPrice(cartItem.RetailPrice.Value - float32(cartItem.CatalogInfo.DiscountInfo.Value))
						td = td + cartItem.CatalogInfo.DiscountInfo.Value*cartItem.Quantity
					case model.PercentOffType:
						cart.Items[i].DiscountInfo.MaxValue = cartItem.CatalogInfo.DiscountInfo.MaxValue
						d := uint(float64((cartItem.CatalogInfo.DiscountInfo.Value * uint(cartItem.RetailPrice.Value)) / 100.0))
						if d > cartItem.CatalogInfo.DiscountInfo.MaxValue && cartItem.CatalogInfo.DiscountInfo.MaxValue > 0 {
							d = cartItem.CatalogInfo.DiscountInfo.MaxValue
						}
						dp = model.SetINRPrice(cartItem.RetailPrice.Value - float32(d))
						td = td + (d * cartItem.Quantity)
					default:
					}
					cart.Items[i].DiscountedPrice = dp
					cart.Items[i].DiscountID = cartItem.CatalogInfo.DiscountInfo.ID
					cart.Items[i].TransferPrice = model.SetINRPrice(0)
					gt = gt + uint(dp.Value)*cartItem.Quantity
					applied = true
				}
			}
			if applied == false {
				gt += uint(cartItem.RetailPrice.Value) * cartItem.Quantity
			}
		} else {
			gt += uint(cartItem.RetailPrice.Value) * cartItem.Quantity
		}

	}

	cart.TotalPrice = model.SetINRPrice(float32(tp))
	cart.TotalDiscount = model.SetINRPrice(float32(td))
	cart.GrandTotal = model.SetINRPrice(float32(gt))
	return &cart, nil
}

//GetCartItems function sets the shipping address for cart
func (ci *CartImpl) SetCartAddress(opts *schema.AddressOpts) error {
	ctx := context.TODO()
	address := model.Address{
		ID:                opts.AddressID,
		DisplayName:       opts.DisplayName,
		Line1:             opts.Line1,
		Line2:             opts.Line2,
		District:          opts.District,
		City:              opts.City,
		State:             opts.State,
		PostalCode:        opts.PostalCode,
		Country:           opts.Country,
		PlainAddress:      opts.PlainAddress,
		IsBillingAddress:  opts.IsBillingAddress,
		IsShippingAddress: opts.IsShippingAddress,
		IsDefaultAddress:  opts.IsDefaultAddress,
		ContactNumber:     opts.ContactNumber,
	}
	findQuery := bson.M{
		"user_id": opts.ID,
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"shipping_address": address,
			"billing_address":  address,
			"updated_at":       time.Now().UTC(),
		},
	}

	res, err := ci.DB.Collection(model.CartColl).UpdateOne(ctx, findQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to set the address")
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find cart with id: %s", opts.ID.Hex())
	}

	return nil
}

func (ci *CartImpl) CheckoutCart(id primitive.ObjectID, source, platform, userName string) (*schema.OrderInfo, error) {

	ctx := context.TODO()

	isWeb := false
	if platform == "web" {
		isWeb = true
	}
	grandTotal := 0
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"user_id": id,
		},
	}}

	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path": "$items",
		},
	}}

	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$items.brand_id",
			"items": bson.M{
				"$push": "$items",
			},
			"cartInfo": bson.M{
				"$first": "$$ROOT",
			},
		},
	}}
	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"cartInfo.items":    "$items",
			"cartInfo.brand_id": "$_id",
		},
	}}

	replaceRootStage := bson.D{{
		Key: "$replaceRoot", Value: bson.M{
			"newRoot": "$cartInfo",
		},
	}}

	projectStage := bson.D{{
		Key: "$project", Value: bson.M{
			"items._id": 0,
		},
	}}
	cartCursor, err := ci.DB.Collection(model.CartColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		unwindStage,
		groupStage,
		setStage,
		replaceRootStage,
		projectStage,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "unable to get cart data")
	}

	var cartUnwindBrands []schema.CartUnwindBrand

	if err := cartCursor.All(ctx, &cartUnwindBrands); err != nil {
		return nil, errors.Wrap(err, "error decoding cart")
	}
	var orderItemsOpts []schema.OrderItemOpts

	outOfStockString := ""

	for _, c := range cartUnwindBrands {
		order := schema.OrderItemOpts{
			UserID:          c.UserID,
			BrandID:         c.BrandID,
			ShippingAddress: c.ShippingAddress,
			BillingAddress:  c.BillingAddress,
			OrderItems:      []schema.OrderItem{},
			Source:          source,
			IsWeb:           isWeb,
			Platform:        platform,
			CartType:        model.CartCheckout,
		}

		displayName := strings.ToLower(c.ShippingAddress.DisplayName)
		if displayName == "home" || displayName == "other" || displayName == "work" || displayName == "" {
			order.ShippingAddress.DisplayName = userName
			order.BillingAddress.DisplayName = userName
		}

		for _, item := range c.Items {

			var cv model.GetCatalogVariant
			// url := "http://localhost:8000" + "/api/keeper/catalog/" + item.CatalogID.Hex() + "/variant/" + item.VariantID.Hex()

			url := ci.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + item.CatalogID.Hex() + "/variant/" + item.VariantID.Hex()
			client := http.Client{}
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate request to get catalog & variant")
			}
			req.Header.Add("Authorization", ci.App.Config.HypdApiConfig.Token)
			resp, err := client.Do(req)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to fetch catlog data")
			}
			defer resp.Body.Close()

			//Read the response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ci.Logger.Err(err).Msgf("failed to read response from api %s", url)
				return nil, errors.Wrap(err, "failed to get brandinfo")
			}
			if err := json.Unmarshal(body, &cv); err != nil {
				ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
				return nil, errors.Wrap(err, "failed to decode body into struct")
			}
			if !cv.Success {
				ci.Logger.Err(errors.New("success false from inventory")).Str("body", string(body)).Msg("got success false response from inventory")
				return nil, errors.New("got success false response from catalog")
			}
			if cv.Payload.InventoryInfo.UnitInStock == 0 || cv.Payload.InventoryInfo.Status.Value == model.OutOfStockStatus {
				outOfStockString = outOfStockString + fmt.Sprintf("item %s is out of stock", item.CatalogInfo.Name)
				continue
			}
			if cv.Payload.InventoryInfo.UnitInStock < int(item.Quantity) {
				outOfStockString = outOfStockString + fmt.Sprintf("only %d unit available for item %s in stock", cv.Payload.InventoryInfo.UnitInStock, item.CatalogInfo.Name) + "\n"
				continue
			}
			item.CatalogInfo.TransferPrice = cv.Payload.TransferPrice
			// item.CatalogInfo.BasePrice = cv.Payload.BasePrice
			// item.CatalogInfo.RetailPrice = cv.Payload.RetailPrice

			var dp *model.Price

			if cv.Payload.DiscountInfo != nil {

				switch cv.Payload.DiscountInfo.Type {
				case model.FlatOffType:
					dp = model.SetINRPrice(cv.Payload.RetailPrice.Value - float32(cv.Payload.DiscountInfo.Value))
				case model.PercentOffType:
					d := uint(float64((cv.Payload.DiscountInfo.Value * uint(cv.Payload.RetailPrice.Value)) / 100.0))
					if d > cv.Payload.DiscountInfo.MaxValue && cv.Payload.DiscountInfo.MaxValue > 0 {
						d = cv.Payload.DiscountInfo.MaxValue
					}
					dp = model.SetINRPrice(cv.Payload.RetailPrice.Value - float32(d))
				default:
				}
			}

			it := schema.OrderItem{
				CatalogID: item.CatalogID,
				VariantID: item.VariantID,
				CatalogInfo: schema.OrderCatalogInfo{
					ID:      item.CatalogID,
					BrandID: item.BrandID,
					Name:    item.CatalogInfo.Name,
					FeaturedImage: schema.Img{
						SRC: cv.Payload.FeaturedImage.SRC,
					},
					VariantType: item.CatalogInfo.VariantType,
					Variant: schema.OrderVariant{
						ID:        item.VariantID,
						Attribute: cv.Payload.Variant.Attribute,
						SKU:       cv.Payload.Variant.SKU,
					},
					ETA:           item.CatalogInfo.ETA,
					HSNCode:       item.CatalogInfo.HSNCode,
					TransferPrice: cv.Payload.TransferPrice,
				},
				Tax:         item.CatalogInfo.Tax,
				BasePrice:   &cv.Payload.BasePrice,
				RetailPrice: &cv.Payload.RetailPrice,
				Quantity:    item.Quantity,
			}
			if !cv.Payload.DiscountInfo.ID.IsZero() {
				it.DiscountID = cv.Payload.DiscountInfo.ID
				it.DiscountInfo = cv.Payload.DiscountInfo
				it.DiscountedPrice = dp
				grandTotal -= int(dp.Value)
			}
			grandTotal += int(cv.Payload.RetailPrice.Value)
			order.OrderItems = append(order.OrderItems, it)
		}
		orderItemsOpts = append(orderItemsOpts, order)
	}

	if len(outOfStockString) > 0 {
		return nil, errors.Errorf(outOfStockString)
	}

	var coupon schema.CouponOrderOpts

	if cartUnwindBrands[0].Coupon != nil {
		coupon.ID = cartUnwindBrands[0].Coupon.ID
		coupon.Code = cartUnwindBrands[0].Coupon.Code
		if cartUnwindBrands[0].Coupon.Type == model.FlatOffType {
			coupon.AppliedValue = model.SetINRPrice(float32(cartUnwindBrands[0].Coupon.Value))
		} else if cartUnwindBrands[0].Coupon.Type == model.PercentOffType {
			av := grandTotal * cartUnwindBrands[0].Coupon.Value
			if av > int(cartUnwindBrands[0].Coupon.MaxDiscount.Value) {
				av = int(cartUnwindBrands[0].Coupon.MaxDiscount.Value)
			}
			coupon.AppliedValue = model.SetINRPrice(float32(av))
		}

		if cartUnwindBrands[0].Coupon.Type != model.FreeDelivery {
			for i := range orderItemsOpts {
				orderItemsOpts[i].Coupon = &coupon
			}
		}

	}

	// b, err := json.MarshalIndent(orderItemsOpts, "", "  ")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Print(string(b))

	//Create Order
	coURL := ci.App.Config.HypdApiConfig.OrderApi + "/api/order"

	var orderResp schema.OrderResp
	reqBody, err := json.Marshal(orderItemsOpts)
	if err != nil {
		ci.Logger.Err(err).Interface("orderItemsOpts", orderItemsOpts).Msgf("failed to prepare request json to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, coURL, bytes.NewBuffer(reqBody))
	if err != nil {
		ci.Logger.Err(err).Interface("orderItemsOpts", orderItemsOpts).Msgf("failed to create request to create order %s", coURL)
		return nil, errors.Wrap(err, "failed to create request to generete order")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ci.App.Config.HypdApiConfig.Token)
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		ci.Logger.Err(err).RawJSON("responseBody", reqBody).Msgf("failed to send request to api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).RawJSON("reqBody", reqBody).Msgf("failed to read response from api %s", coURL)
		return nil, errors.Wrap(err, "failed to get order info")
	}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return nil, errors.Wrap(err, "failed to decode body into struct")
	}
	if !orderResp.Success {
		ci.Logger.Err(errors.New("success false from order")).Str("body", string(body)).Msg("got success false response from order")
		return nil, errors.Errorf("%s - got success false response from order", string(body))
	}

	return &orderResp.Payload, nil
}

func (ci *CartImpl) ClearCart(id primitive.ObjectID) error {

	findQuery := bson.M{
		"user_id": id,
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"total_price.value":    0,
			"total_discount.value": 0,
			"grand_total.value":    0,
		},
		"$unset": bson.M{
			"items":  "",
			"coupon": "",
		},
	}

	res, err := ci.DB.Collection(model.CartColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to query for document")
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find cart for user with id: %s", id)
	}

	return nil
}

func (ci *CartImpl) AddDiscountInCartItems(opts *schema.DiscountInCartItemsOpts) {
	if opts.CatalogID.IsZero() || len(opts.VariantsID) == 0 {
		return
	}
	filter := bson.M{
		"items.catalog_id": opts.CatalogID,
		"items.variant_id": bson.M{
			"$in": opts.VariantsID,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"items.$.discount_id": opts.ID,
			"items.$.discount_info": model.DiscountInfo{
				ID:       opts.ID,
				Type:     opts.Type,
				Value:    opts.Value,
				MaxValue: opts.MaxValue,
			},
		},
	}

	if _, err := ci.DB.Collection(model.CartColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to add discount in cart items")
	}
}

func (ci *CartImpl) RemoveDiscountInCartItems(opts *schema.DiscountInCartItemsOpts) {
	if opts.CatalogID.IsZero() || len(opts.VariantsID) == 0 {
		return
	}
	filter := bson.M{
		"catalog_id": opts.CatalogID,
		"variant_id": bson.M{
			"$in": opts.VariantsID,
		},
	}

	update := bson.M{
		"$unset": bson.M{
			"items.$.discount_id":   1,
			"items.$.discount_info": 1,
		},
	}

	if _, err := ci.DB.Collection(model.CartColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to add discount in cart items")
	}
}

func (ci *CartImpl) UpdateInventoryStatus(opts *schema.InventoryUpdateOpts) {
	filter := bson.M{
		"items.catalog_id": opts.CatalogID,
		"items.variant_id": opts.VariantID,
	}

	var update bson.M
	if opts.UnitInStock > 0 {
		update = bson.M{
			"$set": bson.M{
				"items.$.in_stock": true,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"items.$.in_stock": false,
			},
		}
	}
	if _, err := ci.DB.Collection(model.CartColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to update stock in cart items")
	}
}

func (ci *CartImpl) UpdateCatalogInfo(id primitive.ObjectID) {
	var s model.GetAllCatalogInfoResp

	url := ci.App.Config.HypdApiConfig.CatalogApi + "/api/keeper/catalog/" + id.Hex()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ci.Logger.Err(errors.Wrapf(err, "failed to request to get catalog info"))
	}
	req.Header.Add("Authorization", ci.App.Config.HypdApiConfig.Token)
	resp, err := client.Do(req)
	if err != nil {
		ci.Logger.Err(errors.Wrapf(err, "unable to fetch catlog data"))
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ci.Logger.Err(err).Msgf("failed to read response from api %s", url)
		ci.Logger.Err(errors.Wrap(err, "failed to get catalog info"))
	}
	if err := json.Unmarshal(body, &s); err != nil {
		ci.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		ci.Logger.Err(errors.Wrap(err, "failed to decode body into struct"))
	}
	if !s.Success {
		ci.Logger.Err(errors.New("success false from catalog")).Str("body", string(body)).Msg("got success false response from catalog")
		ci.Logger.Err(errors.New("got success false response from catalog"))
	}

	filter := bson.M{
		"catalog_id": id,
	}

	catalogInfo := model.CatalogInfo{
		ID:            s.Payload.ID,
		Name:          s.Payload.Name,
		BrandID:       s.Payload.BrandID,
		FeaturedImage: s.Payload.FeaturedImage,
		VariantType:   s.Payload.VariantType,
		Variants:      s.Payload.Variants,
		HSNCode:       s.Payload.HSNCode,
		ETA:           s.Payload.ETA,
		Status:        s.Payload.Status,
		DiscountInfo:  s.Payload.DiscountInfo,
		TransferPrice: s.Payload.TransferPrice,
	}

	update := bson.M{
		"$set": bson.M{
			"items.$.catalog_info": catalogInfo,
			"items.$.brand_info":   s.Payload.BrandInfo,
			"items.$.base_price":   s.Payload.BasePrice,
			"items.$.retail_price": s.Payload.RetailPrice,
		},
	}

	if _, err := ci.DB.Collection(model.CartColl).UpdateMany(context.TODO(), filter, update); err != nil {
		ci.Logger.Err(err).Interface("id", id).Msg("failed to update catalog info in cart items")
	}
}

func (ci *CartImpl) ApplyCoupon(user_id primitive.ObjectID, opts *schema.ApplyCouponOpts) error {

	coupon := model.Coupon{
		ID:               opts.CouponID,
		Code:             opts.Code,
		Description:      opts.Description,
		Type:             opts.Type,
		Value:            opts.Value,
		ApplicableON:     opts.ApplicableON,
		MaxDiscount:      opts.MaxDiscount,
		MinPurchaseValue: opts.MinPurchaseValue,
		ValidAfter:       opts.ValidAfter,
		ValidBefore:      opts.ValidBefore,
		Status:           opts.Status,
	}

	findQuery := bson.M{"user_id": user_id}
	updateQuery := bson.M{"$set": bson.M{
		"coupon": coupon,
	}}

	res, err := ci.DB.Collection(model.CartColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to add coupon to cart")
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find cart for user")
	}

	return nil
}

func (ci *CartImpl) RemoveCoupon(user_id primitive.ObjectID) error {
	findQuery := bson.M{"user_id": user_id}
	updateQuery := bson.M{"$unset": bson.M{
		"coupon": 0,
	}}

	res, err := ci.DB.Collection(model.CartColl).UpdateOne(context.TODO(), findQuery, updateQuery)
	if err != nil {
		return errors.Wrapf(err, "unable to add coupon to cart")
	}
	if res.MatchedCount == 0 {
		return errors.Errorf("unable to find cart for user")
	}
	return nil
}

func (ci *CartImpl) UpdateInventoryStatusInsideCatalogInfo(opts *schema.InventoryUpdateOpts) {
	filter := bson.M{
		"items.catalog_id": opts.CatalogID,
	}

	var update bson.M
	if opts.UnitInStock > 0 {
		update = bson.M{
			"$set": bson.M{
				"items.$.catalog_info.variants.$[elem].inventory_info.unit_in_stock": opts.UnitInStock,
				"items.$.catalog_info.variants.$[elem].inventory_info.status.value":  model.InStockStatus,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"items.$.catalog_info.variants.$[elem].inventory_info.unit_in_stock": 0,
				"items.$.catalog_info.variants.$[elem].inventory_info.status.value":  model.OutOfStockStatus,
			},
		}
	}
	updateOpts := options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: bson.A{
				bson.M{
					"elem._id": opts.VariantID,
				},
			},
		},
	}

	if _, err := ci.DB.Collection(model.CartColl).UpdateMany(context.TODO(), filter, update, &updateOpts); err != nil {
		ci.Logger.Err(err).Interface("opts", opts).Msg("failed to update stock in cart items")
	}
}
