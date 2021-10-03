package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SeriesSubCollectionOpts struct {
	Thumbnail Img      `json:"thumbnail" validate:"required"`
	SeriesIDs []string `json:"series_ids" validate:"required"`
}

type CreateCollectionOpts struct {
	Name                string                    `json:"name" validate:"required" `
	Type                string                    `json:"type" validate:"required,oneof=hashtag influencer brand series" `
	Genders             []string                  `json:"genders" validate:"required,dive,oneof=M F O"`
	Hashtags            []string                  `json:"hashtags"`
	BrandIDs            []string                  `json:"brand_ids"`
	InfluencerIDs       []string                  `json:"influencer_ids"`
	SeriesSubCollection []SeriesSubCollectionOpts `json:"series_subcollection"`
}

type PebbleCollectionKafkaUpdateOpts struct {
	ID                  primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                string                      `json:"name,omitempty" bson:"name,omitempty"`
	Type                string                      `json:"type,omitempty" bson:"type,omitempty"`
	Genders             []string                    `json:"genders,omitempty" bson:"genders,omitempty"`
	Hashtags            []string                    `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	BrandIDs            []string                    `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo           []model.BrandInfo           `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	InfluencerIDs       []string                    `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo      []model.InfluencerInfo      `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	SeriesSubCollection []model.SeriesSubCollection `json:"series_subcollection,omitempty" bson:"series_subcollection,omitempty"`
	Status              string                      `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt           time.Time                   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time                   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type GetCollectionFilter struct {
	Genders []string `json:"genders,omitempty" queryparam:"genders"`
	Page    int      `json:"page,omitempty" queryparam:"page"`
}

type GetPebbleCollectionESResp struct {
	ID                  primitive.ObjectID          `json:"id,omitempty" bson:"_id,omitempty"`
	Name                string                      `json:"name,omitempty" bson:"name,omitempty"`
	Type                string                      `json:"type,omitempty" bson:"type,omitempty"`
	Genders             []string                    `json:"genders,omitempty" bson:"genders,omitempty"`
	Hashtags            []string                    `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	BrandIDs            []string                    `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	BrandInfo           []model.BrandInfo           `json:"brand_info,omitempty" bson:"brand_info,omitempty"`
	InfluencerIDs       []string                    `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluencerInfo      []model.InfluencerInfo      `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	SeriesSubCollection []model.SeriesSubCollection `json:"series_subcollection,omitempty" bson:"series_subcollection,omitempty"`
	Status              string                      `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt           time.Time                   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           time.Time                   `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateCollectionOpts struct {
	ID                  primitive.ObjectID        `json:"id" validate:"required" `
	Name                string                    `json:"name"`
	Genders             []string                  `json:"genders"`
	Hashtags            []string                  `json:"hashtags"`
	BrandIDs            []string                  `json:"brand_ids"`
	InfluencerIDs       []string                  `json:"influencer_ids"`
	SeriesSubCollection []SeriesSubCollectionOpts `json:"series_subcollection"`
	Status              string                    `json:"status"`
}

type GetCollectionsKeeperFilter struct {
	Page   int      `qs:"page"`
	Status []string `qs:"status"`
}

// CollectionResp serialize the get collections api response
type CollectionResp struct {
	ID                  primitive.ObjectID          `json:"id,omitempty" bson:"_id,omitempty"`
	Name                string                      `json:"name" bson:"name,omitempty"`
	Type                string                      `json:"type,omitempty" bson:"type,omitempty"`
	Genders             []string                    `json:"genders,omitempty" bson:"genders,omitempty"`
	Hashtags            []string                    `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	BrandIDs            []string                    `json:"brand_ids,omitempty" bson:"brand_ids,omitempty"`
	InfluencerIDs       []string                    `json:"influencer_ids,omitempty" bson:"influencer_ids,omitempty"`
	InfluenncerInfo     []model.InfluencerInfo      `json:"influencer_info,omitempty" bson:"influencer_info,omitempty"`
	SeriesSubCollection []model.SeriesSubCollection `json:"series_subcollection,omitempty" bson:"series_subcollection,omitempty"`
	Status              string                      `json:"status,omitempty" bson:"status,omitempty"`
}
