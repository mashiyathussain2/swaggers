package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetBrandSchema struct {
	ID   string   `json:"_id"`
	Slug string   `json:"slug"`
	Name string   `json:"name"`
	Logo *ImgResp `json:"logo"`
}

type IsLikedByUserBucketSchema struct {
	DocCount uint `json:"doc_count"`
}

type LikeCountBucketSchema struct {
	Key           string                      `json:"key"`
	DocCount      uint                        `json:"doc_count"`
	IsLikedByUser []IsLikedByUserBucketSchema `json:"is_liked_by_user"`
}

type LikeCountSchema struct {
	Buckets []LikeCountBucketSchema `json:"buckets"`
}

type LikeCountAggResp struct {
	LikeCount []LikeCountSchema `json:"like_count"`
}

type GetInfluencerProductESResp struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	InfluencerID primitive.ObjectID   `json:"influencer_id" bson:"influencer_id,omitempty"`
	CatalogIDs   []primitive.ObjectID `json:"catalog_ids" bson:"catalog_ids"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
