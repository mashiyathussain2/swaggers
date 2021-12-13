package app

import (
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"go-app/server/kafka"
	"sync"

	"github.com/pkg/errors"
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
	csp.Logger.Log().Msg("Initiating processing of content message")
	fmt.Println("Initiating processing of content message")
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

	// Only pushing pebble type content in elasticsearch
	if contentSchema.Type != model.PebbleType {
		return
	}

	// Removing content from index if is active set to false
	if !contentSchema.IsProcessed {
		// m := segKafka.Message{
		// 	Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
		// 	Value: nil,
		// }
		// csp.App.ContentFullProducer.Publish(m)
		csp.Logger.Log().Msg("pebble not processed yet")
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
				csp.Logger.Err(err).Interface("data", contentSchema).Msg("failed to get content catalog info")
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
		LikedBy:        contentSchema.LikedBy,
		CommentCount:   contentSchema.CommentCount,
		ViewCount:      contentSchema.ViewCount,
		Label:          contentSchema.Label,
		IsProcessed:    contentSchema.IsProcessed,
		IsActive:       contentSchema.IsActive,
		Paths:          contentSchema.Paths,
		Caption:        contentSchema.Caption,
		Hashtags:       contentSchema.Hashtags,
		CreatedAt:      contentSchema.CreatedAt,
		SeriesIDs:      contentSchema.SeriesIDs,
		CreatorID:      contentSchema.CreatorID,
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

	// If a new like is registered inside `like` collection
	if s.Meta.Operation == "i" || s.Meta.Operation == "" {
		if s.Data == nil {
			return
		}
		var likeSchema schema.ProcessLikeOpts
		likeByteData, err := json.Marshal(s.Data)
		if err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode like update data fields into bytes")
			return
		}
		if err := json.Unmarshal(likeByteData, &likeSchema); err != nil {
			csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
			return
		}

		val, err := json.Marshal(schema.ProcessLikeESResp{ID: likeSchema.ID, ResourceType: likeSchema.ResourceType, ResourceID: likeSchema.ResourceID, UserID: likeSchema.UserID, CreatedAt: likeSchema.CreatedAt})
		if err != nil {
			csp.Logger.Err(err).Interface("likes", likeSchema).Msg("failed to convert struct to json")
			return
		}
		csp.App.LikeProducer.Publish(segKafka.Message{Key: []byte(likeSchema.ID.Hex()), Value: val})
		return
	}

	// If a new like doc was deleted inside `like` collection
	if s.Meta.Operation == "d" {
		csp.App.LikeProducer.Publish(segKafka.Message{Key: []byte(s.Meta.ID.(primitive.ObjectID).Hex()), Value: nil})
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
	if s.Meta.Operation == "i" || s.Meta.Operation == "" {
		if s.Data == nil {
			return
		}
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
		val, err := json.Marshal(schema.ProcessViewESResp{ID: viewSchema.ID, ResourceType: viewSchema.ResourceType, ResourceID: viewSchema.ResourceID, UserID: viewSchema.UserID, CreatedAt: viewSchema.CreatedAt, Duration: viewSchema.Duration})
		if err != nil {
			csp.Logger.Err(err).Interface("views", viewSchema).Msg("failed to convert struct to json")
			return
		}
		csp.App.ViewProducer.Publish(segKafka.Message{Key: []byte(viewSchema.ID.Hex()), Value: val})
		return
	}

	// If a new like doc was deleted inside `like` collection
	if s.Meta.Operation == "d" {
		csp.App.ViewProducer.Publish(segKafka.Message{Key: []byte(s.Meta.ID.(primitive.ObjectID).Hex()), Value: nil})
		return
	}

	if s.Meta.Operation == "u" {
		if updates, ok := s.Meta.Updates.(bson.D).Map()["changed"]; ok {
			if sync, ok := updates.(primitive.D).Map()["sync"]; ok {
				if sync.(bool) == true {
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
					val, err := json.Marshal(schema.ProcessViewESResp{ID: viewSchema.ID, ResourceType: viewSchema.ResourceType, ResourceID: viewSchema.ResourceID, UserID: viewSchema.UserID, CreatedAt: viewSchema.CreatedAt, Duration: viewSchema.Duration})
					if err != nil {
						csp.Logger.Err(err).Interface("views", viewSchema).Msg("failed to convert struct to json")
						return
					}
					csp.App.ViewProducer.Publish(segKafka.Message{Key: []byte(viewSchema.ID.Hex()), Value: val})
				}
			}
		}

	}

}

func (csp *ContentUpdateProcessor) ProcessLiveOrder(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)

	var liveOrder schema.LiveOrderKafkaMessage
	if err := json.Unmarshal(message.Value, &liveOrder); err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode live order update data fields into bytes")
		return
	}
	if liveOrder.ID.IsZero() {
		csp.Logger.Err(errors.New("invalid live")).Interface("liveOrder", liveOrder).Msg("invalid live id")
		return
	}

	live, err := csp.App.Live.GetLiveStreamByID(liveOrder.ID)
	if err != nil {
		csp.Logger.Err(err).Interface("liveOrder", liveOrder).Msg("failed to query live order data")
		return
	}
	if live == nil {
		csp.Logger.Err(errors.New("invalid live")).Interface("liveOrder", liveOrder).Msg("failed to find live order data")
		return
	}

	opts := schema.PushNewOrderOpts{
		ARN:          live.IVS.Channel.ARN,
		Name:         liveOrder.Name,
		ProfileImage: liveOrder.ProfileImage,
	}
	csp.App.Live.PushOrder(&opts)
}

func (csp *ContentUpdateProcessor) ProcessSeriesMessage(msg kafka.Message) {
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
		csp.App.PebbleSeriesProducer.Publish(m)
		return
	}
	// When content is added/updated it will be sync with %pebble_series_full topic.
	var seriesSchema schema.PebbleSeriesKafkaUpdateOpts
	contentByteData, err := json.Marshal(s.Data)
	if err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode catalog update data fields into bytes")
		return
	}
	if err := json.Unmarshal(contentByteData, &seriesSchema); err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}
	// Removing content from index if is active set to false
	if !seriesSchema.IsActive {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		csp.App.PebbleSeriesProducer.Publish(m)
		return
	}
	if len(seriesSchema.PebbleIds) > 0 {
		isActive := true
		mi, err := csp.App.Series.GetContentForPebbleSeries(&schema.GetContentFilter{IDs: seriesSchema.PebbleIds, IsActive: &isActive})
		if err != nil {
			csp.Logger.Err(err).Interface("id", seriesSchema).Msg("failed to get pebble media info")
			return
		}
		seriesSchema.PebbleInfo = mi

	}

	val, err := json.Marshal(model.PebbleSeries{
		ID:         seriesSchema.ID,
		Name:       seriesSchema.Name,
		Thumbnail:  seriesSchema.Thumbnail,
		PebbleIds:  seriesSchema.PebbleIds,
		PebbleInfo: seriesSchema.PebbleInfo,
		Label:      seriesSchema.Label,
		IsActive:   seriesSchema.IsActive,
		CreatedAt:  seriesSchema.CreatedAt,
		UpdatedAt:  seriesSchema.UpdatedAt,
	})
	if err != nil {
		csp.Logger.Err(err).Interface("contentSchema", seriesSchema).Msg("failed to convert contentSchema to json")
		return
	}
	m := segKafka.Message{
		Key:   []byte(seriesSchema.ID.Hex()),
		Value: val,
	}
	csp.App.PebbleSeriesProducer.Publish(m)
	return
}

func (csp *ContentUpdateProcessor) ProcessCollectionMessage(msg kafka.Message) {
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
		csp.App.PebbleSeriesProducer.Publish(m)
		return
	}
	// When collection is added/updated it will be sync with %pebble_collection_full topic.
	var collectionSchema schema.PebbleCollectionKafkaUpdateOpts
	contentByteData, err := json.Marshal(s.Data)
	if err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode catalog update data fields into bytes")
		return
	}
	if err := json.Unmarshal(contentByteData, &collectionSchema); err != nil {
		csp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}
	// // Removing content from index if is active set to false
	// if !seriesSchema.IsActive {
	// 	m := segKafka.Message{
	// 		Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
	// 		Value: nil,
	// 	}
	// 	csp.App.PebbleSeriesProducer.Publish(m)
	// 	return
	// }
	switch collectionSchema.Type {

	case model.BrandCollection:
		res, err := csp.App.Content.GetBrandInfo(collectionSchema.BrandIDs)
		if err != nil {
			csp.Logger.Err(err).Interface("collectionSchema", collectionSchema).Msg("failed to get BrandInfo")
			return
		}
		collectionSchema.BrandInfo = res

	case model.InfluencerCollection:
		res, err := csp.App.Content.GetInfluencerInfo(collectionSchema.InfluencerIDs)
		if err != nil {
			csp.Logger.Err(err).Interface("collectionSchema", collectionSchema).Msg("failed to get BrandInfo")
			return
		}
		collectionSchema.InfluencerInfo = res

	}

	val, err := json.Marshal(model.Collection{
		ID:                  collectionSchema.ID,
		Name:                collectionSchema.Name,
		Type:                collectionSchema.Type,
		Genders:             collectionSchema.Genders,
		Hashtags:            collectionSchema.Hashtags,
		BrandIDs:            collectionSchema.BrandIDs,
		BrandInfo:           collectionSchema.BrandInfo,
		InfluencerIDs:       collectionSchema.InfluencerIDs,
		InfluencerInfo:      collectionSchema.InfluencerInfo,
		SeriesSubCollection: collectionSchema.SeriesSubCollection,
		Status:              collectionSchema.Status,
		CreatedAt:           collectionSchema.CreatedAt,
		UpdatedAt:           collectionSchema.UpdatedAt,
	})
	if err != nil {
		csp.Logger.Err(err).Interface("collectionSchema", collectionSchema).Msg("failed to convert collectionSchema to json")
		return
	}
	m := segKafka.Message{
		Key:   []byte(collectionSchema.ID.Hex()),
		Value: val,
	}
	csp.App.PebbleCollectionProducer.Publish(m)
	return
}

// func (csp *ContentUpdateProcessor) ProcessPebbleSeries(msg kafka.Message) {
// 	var s *schema.KafkaMessage
// 	message := msg.(segKafka.Message)
// 	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
// 		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode Pebble Series message")
// 		return
// 	}

// 	csp.ProcessSeriesMessage(message)
// }

func (csp *ContentUpdateProcessor) ProcessContentMessageForSeries(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		csp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}

	// When content is added/updated it will be sync with %content_full topic.
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

	// Only pushing pebble type content in elasticsearch
	if contentSchema.Type != model.PebbleType {
		return
	}

	// Removing content from index if is active set to false
	if !contentSchema.IsActive {
		csp.App.Series.UpdateSeriesLastSync(contentSchema.ID)
		return
	}

	if s.Meta.Operation == "u" {
		if updates, ok := s.Meta.Updates.(bson.D).Map()["changed"]; ok {
			if _, ok := updates.(primitive.D).Map()["is_active"]; ok {

				csp.Logger.Info().Msg("syncing content in series")
				csp.App.Series.UpdateSeriesLastSync(contentSchema.ID)

			} else if _, ok := updates.(primitive.D).Map()["like_count"]; ok {

				csp.Logger.Info().Msg("syncing content in series")
				csp.App.Series.UpdateSeriesLastSync(contentSchema.ID)

			}
		}
	}

	return
}

func (csp *ContentUpdateProcessor) ProcessLikeForSeries(msg kafka.Message) {
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
		csp.App.Series.UpdateSeriesLastSync(likeSchema.ResourceID)
		return
	}

	// unliking
	if s.Meta.Operation == "d" {
		likeSchema := schema.ProcessLikeOpts{
			ID: s.Meta.ID.(primitive.ObjectID),
		}
		csp.App.Series.UpdateSeriesLastSync(likeSchema.ResourceID)
		return
	}
}
