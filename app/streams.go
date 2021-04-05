package app

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("got brand update")
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		bp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode brand update message")
		return
	}
	fmt.Println(string(message.Value))
	if s.Meta.Operation == "d" {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		bp.App.BrandFullProducer.Publish(m)
		return
	}

	fmt.Println(s.Data)
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
