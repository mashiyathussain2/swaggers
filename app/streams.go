package app

import (
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"go-app/server/kafka"

	"github.com/rs/zerolog"
	segKafka "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type CatalogProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

// InitCatalogProcessor returns a new instance of CatalogProcessor
func InitCatalogProcessor(opts *CatalogProcessorOpts) *CatalogProcessor {
	cp := CatalogProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}
	return &cp
}

func (cp *CatalogProcessor) ProcessCatalogUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}
	cp.Logger.Info().RawJSON("json", message.Value).Msg("got the catalog")
	// If delete operation is performed then removing the document from index as well
	if s.Meta.Operation == "d" {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		cp.App.CatalogFullProducer.Publish(m)
		return
	}

	var catalog schema.CatalogKafkaMessage
	catalogByteData, err := json.Marshal(s.Data)
	if err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode catalog update data fields into bytes")
		return
	}
	if err := json.Unmarshal(catalogByteData, &catalog); err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	// not doing anything for unpublished catalog
	if catalog.Status.Value != model.Publish {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		cp.App.CatalogFullProducer.Publish(m)
		return
	}

	updatedCatalogSchema, err := cp.App.KeeperCatalog.GetAllCatalogInfo(catalog.ID)
	if err != nil {
		cp.Logger.Err(err).Interface("catalog", catalog).Msgf("failed to sync catalog with id:%s", catalog.ID)
		return
	}

	val, err := json.Marshal(updatedCatalogSchema)
	if err != nil {
		cp.Logger.Err(err).Interface("catalog", catalog).Msgf("failed to convert catalog with id:%s into json", catalog.ID)
		return
	}

	m := segKafka.Message{
		Key:   []byte(updatedCatalogSchema.ID.Hex()),
		Value: val,
	}
	cp.App.CatalogFullProducer.Publish(m)
}

func (cp *CatalogProcessor) ProcessDiscountUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode discount update message")
		return
	}

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
	cp.App.KeeperCatalog.SyncCatalog(discount.CatalogID)
}

func (cp *CatalogProcessor) ProcessInventoryUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode inventory update message")
		return
	}

	var inventory model.Inventory
	inventoryBytes, err := json.Marshal(s.Data)
	if err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode inventory update data fields into bytes")
		return
	}
	if err := json.Unmarshal(inventoryBytes, &inventory); err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}
	cp.App.KeeperCatalog.SyncCatalog(inventory.CatalogID)
}

func (cp *CatalogProcessor) ProcessCatalogContentUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode inventory update message")
		return
	}

	var content schema.CatalogContentKafkaMessage
	contentBytes, err := json.Marshal(s.Data)
	if err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode content update data fields into bytes")
		return
	}
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	if content.Type == "catalog_content" {
		cp.Logger.Info().Msg("syncing catalog content")
		cp.App.KeeperCatalog.SyncCatalogContent(content.ID)
		return
	}

	if content.Type == "review_story" {
		if s.Meta.Operation == "u" {
			if updates, ok := s.Meta.Updates.(bson.D).Map()["changed"]; ok {
				if val, ok := updates.(primitive.D).Map()["is_active"]; ok {
					if val.(bool) {
						cp.Logger.Info().Msg("syncing catalog review")
						cp.App.Review.ProcessReviewStory(s.Meta.ID.(primitive.ObjectID))
					}
				}
			}
		}
	}
}

func (cp *CatalogProcessor) ProcessGroupUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode inventory update message")
		return
	}

	var group schema.GroupChangeKafkaMessage
	groupBytes, err := json.Marshal(s.Data)
	if err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode group update data fields into bytes")
		return
	}
	if err := json.Unmarshal(groupBytes, &group); err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	if group.Status.Value == model.Publish {
		cp.App.KeeperCatalog.SyncCatalogs(group.CatalogIDs)
	}
}

func (cp *CatalogProcessor) ProcessBrandUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode brand update message")
		return
	}
	if s.Meta.Operation == "u" {
		if updates, ok := s.Meta.Updates.(bson.D).Map()["changed"]; ok {
			var updateCatalog bool

			if _, ok := updates.(primitive.D).Map()["name"]; ok {
				updateCatalog = true
			}
			if _, ok := updates.(primitive.D).Map()["logo"]; ok {
				updateCatalog = true
			}

			if updateCatalog {
				cp.App.KeeperCatalog.SyncCatalogByBrandID(s.Meta.ID.(primitive.ObjectID))
			}
		}
	}
}

type CollectionProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type CollectionProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

// InitCollectionProcessor returns a new instance of CollectionProcessor
func InitCollectionProcessor(opts *CollectionProcessorOpts) *CollectionProcessor {
	cp := CollectionProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}
	return &cp
}

func (cp *CollectionProcessor) ProcessCollectionUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}

	if s.Meta.Operation == "d" {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		cp.App.CollectionFullProducer.Publish(m)
		return
	}

	if s.Meta.Operation == "i" {
		cp.App.Collection.AddCatalogInfoToCollection(s.Meta.ID.(primitive.ObjectID))
		return
	}

	if s.Meta.Operation == "u" {
		if updates, ok := s.Meta.Updates.(bson.D).Map()["changed"]; ok {
			if _, ok := updates.(primitive.D).Map()["sub_collections.0.catalog_ids"]; ok {
				cp.App.Collection.AddCatalogInfoToCollection(s.Meta.ID.(primitive.ObjectID))
			}
		}
	}

	var collectionSchema schema.CollectionKafkaMessageResp
	collByteData, err := json.Marshal(s.Data)
	if err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to decode collection update data fields into bytes")
		return
	}
	if err := json.Unmarshal(collByteData, &collectionSchema); err != nil {
		cp.Logger.Err(err).Interface("data", s.Data).Msg("failed to convert bson to struct")
		return
	}

	if collectionSchema.Status != model.Publish {
		m := segKafka.Message{
			Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
			Value: nil,
		}
		cp.App.CollectionFullProducer.Publish(m)
		return
	}

	collData := schema.CollectionInfoResp{
		ID:        collectionSchema.ID,
		Name:      collectionSchema.Name,
		Type:      collectionSchema.Type,
		Genders:   collectionSchema.Genders,
		Title:     collectionSchema.Title,
		CreatedAt: collectionSchema.CreatedAt,
		UpdatedAt: collectionSchema.UpdatedAt,
		Status:    collectionSchema.Status,
		Order:     collectionSchema.Order,
	}

	for _, subColl := range collectionSchema.SubCollections {
		subCollData := schema.SubCollectionInfoResp{
			ID:                 subColl.ID,
			Name:               subColl.Name,
			CatalogIDs:         subColl.CatalogIDs,
			FeaturedCatalogIDs: subColl.FeaturedCatalogIDs,
			Image:              subColl.Image,
			CreatedAt:          subColl.CreatedAt,
			UpdatedAt:          subColl.UpdatedAt,
		}
		for _, catalogInfo := range subColl.CatalogInfo {
			subCollCatData := schema.SubCollectionCatalogInfoSchema{
				ID:            catalogInfo.ID,
				BrandID:       catalogInfo.BrandID,
				BrandInfo:     catalogInfo.BrandInfo,
				Name:          catalogInfo.Name,
				FeaturedImage: catalogInfo.FeaturedImage,
				Slug:          catalogInfo.Slug,
				VariantType:   catalogInfo.VariantType,
				Status:        catalogInfo.Status,
				BasePrice:     catalogInfo.BasePrice,
				RetailPrice:   catalogInfo.RetailPrice,
				DiscountID:    catalogInfo.DiscountID,
			}
			if catalogInfo.DiscountInfo != nil {
				subCollCatData.DiscountInfo = &schema.SubCollectionCatalogInfoDiscountInfoResp{
					ID:       catalogInfo.DiscountInfo.ID,
					Type:     catalogInfo.DiscountInfo.Type,
					MaxValue: catalogInfo.DiscountInfo.MaxValue,
					Value:    catalogInfo.DiscountInfo.Value,
				}
			}
			for _, variant := range catalogInfo.Variants {
				subCollCatVariantData := schema.SubCollectionCatalogInfoVariantsResp{
					ID:        variant.ID,
					Attribute: variant.Attribute,
					IsDeleted: variant.IsDeleted,
				}
				subCollCatData.Variants = append(subCollCatData.Variants, subCollCatVariantData)
			}
			subCollData.CatalogInfo = append(subCollData.CatalogInfo, subCollCatData)
		}
		collData.SubCollections = append(collData.SubCollections, subCollData)
	}

	val, err := json.Marshal(collData)
	if err != nil {
		cp.Logger.Err(err).Interface("collectionSchema", collData).Msg("failed to convert collectionschema to json")
		return
	}
	m := segKafka.Message{
		Key:   []byte(collData.ID.Hex()),
		Value: val,
	}
	cp.App.CollectionFullProducer.Publish(m)
}

func (cp *CollectionProcessor) ProcessCatalogUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		cp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}
	if s.Meta.Operation == "u" {
		cp.App.Collection.UpdateCollectionCatalogInfo(s.Meta.ID.(primitive.ObjectID))
		return
	}
}

type ReviewProcessor struct {
	App    *App
	Logger *zerolog.Logger
}

type ReviewProcessorOpts struct {
	App    *App
	Logger *zerolog.Logger
}

// InitCollectionProcessor returns a new instance of CollectionProcessor
func InitReviewProcessor(opts *ReviewProcessorOpts) *ReviewProcessor {
	rp := ReviewProcessor{
		App:    opts.App,
		Logger: opts.Logger,
	}
	return &rp
}

func (rp *ReviewProcessor) ProcessReviewUpdate(msg kafka.Message) {
	var s *schema.KafkaMessage
	message := msg.(segKafka.Message)
	if err := bson.UnmarshalExtJSON(message.Value, false, &s); err != nil {
		rp.Logger.Err(err).Interface("msg", message.Value).Msg("failed to decode catalog update message")
		return
	}

	if s.Meta.Operation == "u" {
		if updates, ok := s.Meta.Updates.(bson.D).Map()["changed"]; ok {
			if val, ok := updates.(primitive.D).Map()["is_processed"]; ok {
				if val.(bool) {
					rp.Logger.Info().Msg("syncing review stories")
					resp, err := rp.App.Review.GetReviewInfo(s.Meta.ID.(primitive.ObjectID))
					if err != nil {
						rp.Logger.Err(err).Msgf("failed to sync review story with id: %s", s.Meta.ID.(primitive.ObjectID).Hex())
						return
					}
					fmt.Println("got review info", resp)
					val, err := json.Marshal(resp)
					if err != nil {
						rp.Logger.Err(err).Interface("storySchema", resp).Msg("failed to convert storySchema to json")
						return
					}
					fmt.Println("got review info val", string(val))
					m := segKafka.Message{
						Key:   []byte(s.Meta.ID.(primitive.ObjectID).Hex()),
						Value: val,
					}
					rp.App.ReviewFullProducer.Publish(m)
				}
			}
		}
	}
}
