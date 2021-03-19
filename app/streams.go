package app

import (
	"encoding/json"
	"go-app/model"
	"go-app/schema"
	"go-app/server/kafka"
	"sync"

	"github.com/rs/zerolog"
	segKafka "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContentUpdateProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type ContentUpdateProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

// InitContentUpdateProcessor returns a new instance of InitContentUpdateProcessor
func InitContentUpdateProcessor(opts *ContentUpdateProcessorOpts) *ContentUpdateProcessor {
	csp := ContentUpdateProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}
	return &csp
}

// ProcessBrandMessage processes the brand change in content
func (csp *ContentUpdateProcessor) ProcessBrandMessage(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode brand update message")
		return
	}
	// Ignoring brand insert operation
	if s.Meta.Operation == "i" {
		return
	}
	// Ignoring brand delete operation
	if s.Meta.Operation == "d" {

	}
	// Handling brand update operation
	if s.Meta.Operation == "u" {
		var brandSchema schema.UpdateContentBrandInfoOpts
		brandByteData, err := json.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode brand update data fields into bytes")
			return
		}
		if err := json.Unmarshal(brandByteData, &brandSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		csp.App.Content.UpdateContentBrandInfo(&brandSchema)
	}
}

// ProcessInfluencerMessage processes the influencer changes in content
func (csp *ContentUpdateProcessor) ProcessInfluencerMessage(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode influencer update message")
		return
	}
	// Ignoring brand insert operation
	if s.Meta.Operation == "i" {
		return
	}
	// Ignoring brand delete operation
	if s.Meta.Operation == "d" {

	}
	// Handling influencer update operation
	if s.Meta.Operation == "u" {
		var influencerSchema schema.UpdateContentInfluencerInfoOpts
		influencerByteData, err := bson.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode influencer update data fields into bytes")
			return
		}
		if err := json.Unmarshal(influencerByteData, &influencerSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		csp.App.Content.UpdateContentInfluencerInfo(&influencerSchema)
	}
}

// ProcessCatalogMessage processes the catalog changes in content
func (csp *ContentUpdateProcessor) ProcessCatalogMessage(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}
	// Ignoring brand insert operation
	if s.Meta.Operation == "i" {
		return
	}
	// Ignoring brand delete operation
	if s.Meta.Operation == "d" {

	}
	// Handling influencer update operation
	if s.Meta.Operation == "u" {
		var catalogSchema schema.UpdateContentCatalogInfoOpts
		catalogByteData, err := bson.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode catalog update data fields into bytes")
			return
		}
		if err := json.Unmarshal(catalogByteData, &catalogSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		csp.App.Content.UpdateContentCatalogInfo(&catalogSchema)
	}
}

func (csp *ContentUpdateProcessor) ProcessContentMessage(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}
	if s.Meta.Operation == "i" {
		var wg sync.WaitGroup
		var contentSchema schema.ContentAddOpts
		contentByteData, err := json.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode catalog update data fields into bytes")
			return
		}
		if err := json.Unmarshal(contentByteData, &contentSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}

		if len(contentSchema.BrandIDs) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Fetching full brand info object
				brandInfo, err := csp.App.Content.GetBrandInfo(contentSchema.BrandIDs)
				if err != nil {
					csp.Logger.Err(err).Interface("data", contentSchema).Msg("failed to get content brand info")
					return
				}
				contentSchema.BrandInfo = brandInfo
			}()
		}

		if len(contentSchema.InfluencerIDs) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Fetching full influencer info object
				influencerInfo, err := csp.App.Content.GetInfluencerInfo(contentSchema.InfluencerIDs)
				if err != nil {
					csp.Logger.Err(err).Interface("data", contentSchema).Msg("failed to get content brand info")
					return
				}
				contentSchema.InfluencerInfo = influencerInfo
			}()
		}

		if len(contentSchema.CatalogIDs) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Fetching full catalog info object
				catalogInfo, err := csp.App.Content.GetCatalogInfo(contentSchema.CatalogIDs)
				if err != nil {
					csp.Logger.Err(err).Interface("data", contentSchema).Msg("failed to get content brand info")
				}
				contentSchema.CatalogInfo = catalogInfo
			}()
		}

		if contentSchema.MediaID != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				switch contentSchema.MediaType {
				case model.ImageType:
					Id, err := primitive.ObjectIDFromHex(contentSchema.MediaID)
					if err != nil {
						csp.Logger.Err(err).Interface("id", contentSchema.MediaID).Msg("invalid media id")
						return
					}
					mi, err := csp.App.Media.GetImageMediaByID(Id)
					if err != nil {
						csp.Logger.Err(err).Interface("id", contentSchema.MediaID).Msg("failed to get image media info")
						return
					}
					contentSchema.MediaInfo = mi
				case model.VideoType:
					Id, err := primitive.ObjectIDFromHex(contentSchema.MediaID)
					if err != nil {
						csp.Logger.Err(err).Interface("id", contentSchema.MediaID).Msg("invalid media id")
						return
					}
					mi, err := csp.App.Media.GetVideoMediaByID(Id)
					if err != nil {
						csp.Logger.Err(err).Interface("id", contentSchema.MediaID).Msg("failed to get video media info")
						return
					}
					contentSchema.MediaInfo = mi
				}
			}()
		}

		wg.Wait()
		val, err := json.Marshal(contentSchema)
		if err != nil {
			csp.Logger.Err(err).Interface("contentSchema", contentSchema).Msg("failed to convert contentSchema to json")
			return
		}
		m := segKafka.Message{
			Key:   []byte(contentSchema.ID),
			Value: val,
		}
		csp.App.ContentFullProducer.Publish(m)
		return
	}
}
