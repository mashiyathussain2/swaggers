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
		influencerByteData, err := json.Marshal(s.Data)
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
		catalogByteData, err := json.Marshal(s.Data)
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

	// If delete operation is performed then removing the document from index as well
	if s.Meta.Operation == "d" {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		csp.App.ContentFullProducer.Publish(m)
		return
	}

	// When content is added/updated it will be sync with %content_full topic.
	var wg sync.WaitGroup
	var contentSchema schema.ContentUpdateOpts
	contentByteData, err := json.Marshal(s.Data)
	if err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode catalog update data fields into bytes")
		return
	}
	if err := json.Unmarshal(contentByteData, &contentSchema); err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	// Removing content from index if is active set to false
	if !contentSchema.IsActive {
		return
	}

	if len(contentSchema.BrandIDs) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Fetching full brand info object
			var ids []string
			for _, id := range contentSchema.BrandIDs {
				ids = append(ids, id.Hex())
			}
			brandInfo, err := csp.App.Content.GetBrandInfo(ids)
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
			var ids []string
			for _, id := range contentSchema.InfluencerIDs {
				ids = append(ids, id.Hex())
			}
			influencerInfo, err := csp.App.Content.GetInfluencerInfo(ids)
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
			var ids []string
			for _, id := range contentSchema.CatalogIDs {
				ids = append(ids, id.Hex())
			}
			catalogInfo, err := csp.App.Content.GetCatalogInfo(ids)
			if err != nil {
				csp.Logger.Err(err).Interface("data", contentSchema).Msg("failed to get content brand info")
			}
			contentSchema.CatalogInfo = catalogInfo
		}()
	}

	if !contentSchema.MediaID.IsZero() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch contentSchema.MediaType {
			case model.ImageType:
				mi, err := csp.App.Media.GetImageMediaByID(contentSchema.MediaID)
				if err != nil {
					csp.Logger.Err(err).Interface("id", contentSchema.MediaID).Msg("failed to get image media info")
					return
				}
				contentSchema.MediaInfo = mi
			case model.VideoType:
				mi, err := csp.App.Media.GetVideoMediaByID(contentSchema.MediaID)
				if err != nil {
					csp.Logger.Err(err).Interface("id", contentSchema.MediaID).Msg("failed to get video media info")
					return
				}
				contentSchema.MediaInfo = mi
			}
		}()
	}

	wg.Wait()

	val, err := json.Marshal(model.Content{
		ID:             contentSchema.ID,
		Type:           contentSchema.Type,
		MediaType:      contentSchema.MediaType,
		MediaID:        contentSchema.MediaID,
		MediaInfo:      contentSchema.MediaInfo,
		BrandIDs:       contentSchema.BrandIDs,
		BrandInfo:      contentSchema.BrandInfo,
		InfluencerIDs:  contentSchema.InfluencerIDs,
		InfluencerInfo: contentSchema.InfluencerInfo,
		CatalogIDs:     contentSchema.CatalogIDs,
		CatalogInfo:    contentSchema.CatalogInfo,
		LikeCount:      contentSchema.LikeCount,
		LikeIDs:        contentSchema.LikeIDs,
		CommentCount:   contentSchema.CommentCount,
		ViewCount:      contentSchema.ViewCount,
		Label:          contentSchema.Label,
		IsProcessed:    contentSchema.IsProcessed,
		IsActive:       contentSchema.IsActive,
		Caption:        contentSchema.Caption,
		Hashtags:       contentSchema.Hashtags,
		CreatedAt:      contentSchema.CreatedAt,
	})
	if err != nil {
		csp.Logger.Err(err).Interface("contentSchema", contentSchema).Msg("failed to convert contentSchema to json")
		return
	}
	m := segKafka.Message{
		Key:   []byte(contentSchema.ID.Hex()),
		Value: val,
	}
	csp.App.ContentFullProducer.Publish(m)
	return
}

func (csp *ContentUpdateProcessor) ProcessLike(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}

	// creating a like
	if s.Meta.Operation == "i" {
		var likeSchema schema.ProcessLikeOpts
		commentByteData, err := json.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode comment update data fields into bytes")
			return
		}
		if err := json.Unmarshal(commentByteData, &likeSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		csp.App.Content.AddContentLike(&likeSchema)
		return
	}
	// unliking
	if s.Meta.Operation == "d" {
		likeSchema := schema.ProcessLikeOpts{
			ID: s.Meta.ID.(primitive.ObjectID),
		}
		csp.App.Content.DeleteContentLike(&likeSchema)
		return
	}
}

func (csp *ContentUpdateProcessor) ProcessComment(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}

	// creating a comment
	if s.Meta.Operation == "i" {
		var commentSchema schema.ProcessCommentOpts
		commentByteData, err := json.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode comment update data fields into bytes")
			return
		}
		if err := json.Unmarshal(commentByteData, &commentSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		csp.App.Content.AddContentComment(&commentSchema)
	}
}

func (csp *ContentUpdateProcessor) ProcessView(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}

	// creating a view
	if s.Meta.Operation == "i" {
		var viewSchema schema.ProcessViewOpts
		viewByteData, err := json.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode view update data fields into bytes")
			return
		}
		if err := json.Unmarshal(viewByteData, &viewSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}
		csp.App.Content.AddContentView(&viewSchema)
	}
}
