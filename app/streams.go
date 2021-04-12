package app

import (
	"context"
	"encoding/json"
	"go-app/model"
	"go-app/schema"
	"go-app/server/kafka"

	"github.com/rs/zerolog"
	segKafka "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BrandProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type BrandProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

func InitBrandProcessor(opts *BrandProcessorOpts) *BrandProcessor {
	bp := BrandProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}
	return &bp
}

func (bp *BrandProcessor) ProcessBrandUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		bp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode brand update message")
		return
	}
	if s.Meta.Operation == "d" {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		bp.App.BrandFullProducer.Publish(m)
		return
	}

	var brand schema.BrandKafkaMessage
	brandByteData, err := json.Marshal(s.Data)
	if err != nil {
		bp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode brand update data fields into bytes")
		return
	}
	if err := json.Unmarshal(brandByteData, &brand); err != nil {
		bp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	brandFullOpts := schema.BrandFullKafkaMessageOpts{
		ID:                 brand.ID,
		Name:               brand.Name,
		LName:              brand.LName,
		Logo:               brand.Logo,
		FulfillmentEmail:   brand.FulfillmentEmail,
		FulfillmentCCEmail: brand.FulfillmentCCEmail,
		RegisteredName:     brand.RegisteredName,
		Domain:             brand.Domain,
		Website:            brand.Website,
		FollowersCount:     brand.FollowersCount,
		FollowingCount:     brand.FollowingCount,
		FollowersID:        brand.FollowersID,
		FollowingID:        brand.FollowingID,
		Bio:                brand.Bio,
		CoverImg:           brand.CoverImg,
		SocialAccount:      brand.SocialAccount,
		CreatedAt:          brand.CreatedAt,
		UpdatedAt:          brand.UpdatedAt,
	}

	val, err := json.Marshal(brandFullOpts)
	if err != nil {
		bp.Logger.Err(err).Interface("brand", brand).Msgf("failed to convert brand with id:%s into json", brand.ID)
		return
	}
	m := segKafka.Message{
		Key:   []byte(brand.ID.Hex()),
		Value: val,
	}
	bp.App.BrandFullProducer.Publish(m)
}

type InfluencerProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type InfluencerProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

func InitInfluencerProcessor(opts *InfluencerProcessorOpts) *InfluencerProcessor {
	bp := InfluencerProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}
	return &bp
}

func (ip *InfluencerProcessor) ProcessInfluencerUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		ip.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode influencer update message")
		return
	}

	if s.Meta.Operation == "d" {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		ip.App.InfluencerFullProducer.Publish(m)
		return
	}

	var influencer schema.InfluencerKafkaMessage
	influencerByteData, err := json.Marshal(s.Data)
	if err != nil {
		ip.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode influencer update data fields into bytes")
		return
	}
	if err := json.Unmarshal(influencerByteData, &influencer); err != nil {
		ip.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	influencerFullOpts := schema.InfluencerFullKafkaMessageOpts{
		ID:             influencer.ID,
		Name:           influencer.Name,
		CoverImg:       influencer.CoverImg,
		ProfileImage:   influencer.ProfileImage,
		SocialAccount:  influencer.SocialAccount,
		ExternalLinks:  influencer.ExternalLinks,
		Bio:            influencer.Bio,
		FollowersID:    influencer.FollowersID,
		FollowingID:    influencer.FollowingID,
		FollowersCount: influencer.FollowersCount,
		FollowingCount: influencer.FollowingCount,
		CreatedAt:      influencer.CreatedAt,
		UpdatedAt:      influencer.UpdatedAt,
	}
	val, err := json.Marshal(influencerFullOpts)
	if err != nil {
		ip.Logger.Err(err).Interface("influencer", influencer).Msgf("failed to convert influencer with id:%s into json", influencer.ID)
		return
	}
	m := segKafka.Message{
		Key:   []byte(influencer.ID.Hex()),
		Value: val,
	}
	ip.App.InfluencerFullProducer.Publish(m)
}

type UserProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type UserProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

func InitUserProcessorOpts(opts *UserProcessorOpts) *UserProcessor {
	up := UserProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}

	return &up
}

func (up *UserProcessor) ProcessCustomerUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		up.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode user update message")
		return
	}
	if s.Meta.Operation == "i" {
		var customer model.Customer
		customerBytes, err := json.Marshal(s.Data)
		if err != nil {
			up.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode customer update data fields into bytes")
			return
		}
		if err := json.Unmarshal(customerBytes, &customer); err != nil {
			up.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		user, err := up.App.User.GetUserByID(customer.UserID)
		if err != nil {
			up.Logger.Err(err).Interface("customer", customer).Msg("failed to get customer user")
			return
		}
		if user.Type == model.CustomerType {
			_, err = up.App.Cart.CreateCart(user.ID)
			if err != nil {
				up.Logger.Err(err).Msg("failed to create cart")
				return
			}
		}
	}

	up.App.CustomerChanges.Commit(context.TODO(), msg)
}

type CartProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type CartProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

func InitCartProcessorOpts(opts *CartProcessorOpts) *CartProcessor {
	cp := CartProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}

	return &cp
}

func (cp *CartProcessor) ProcessDiscountUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode discount update message")
		return
	}

	if s.Meta.Operation == "u" {
		var discount schema.DiscountKafkaMessage
		discountBytes, err := json.Marshal(s.Data)
		if err != nil {
			cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode discount update data fields into bytes")
			return
		}
		if err := json.Unmarshal(discountBytes, &discount); err != nil {
			cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}

		opts := schema.DiscountInCartItemsOpts{
			ID:         discount.ID,
			CatalogID:  discount.CatalogID,
			VariantsID: discount.VariantsID,
			Type:       discount.Type,
			Value:      discount.Value,
			IsActive:   discount.IsActive,
			IsDisabled: discount.IsDisabled,
			MaxValue:   discount.MaxValue,
		}
		if discount.IsActive {
			cp.App.Cart.AddDiscountInCartItems(&opts)
		} else {
			cp.App.Cart.RemoveDiscountInCartItems(&opts)
		}
	}

}

func (cp *CartProcessor) ProcessInventoryUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode discount update message")
		return
	}

	if s.Meta.Operation == "u" {
		var inventory schema.InventoryUpdateKafkaMessage
		inventoryBytes, err := json.Marshal(s.Data)
		if err != nil {
			cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode inventory update data fields into bytes")
			return
		}
		if err := json.Unmarshal(inventoryBytes, &inventory); err != nil {
			cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}

		opts := schema.InventoryUpdateOpts{
			ID:          inventory.ID,
			CatalogID:   inventory.CatalogID,
			VariantID:   inventory.VariantID,
			SKU:         inventory.SKU,
			UnitInStock: inventory.UnitInStock,
		}

		cp.App.Cart.UpdateInventoryStatus(&opts)
	}

}

func (cp *CartProcessor) ProcessCatalogUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode discount update message")
		return
	}
	if s.Meta.Operation == "u" {
		cp.App.Cart.UpdateCatalogInfo(s.Meta.ID.(primitive.ObjectID))
	}
}
