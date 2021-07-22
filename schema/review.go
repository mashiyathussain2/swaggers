package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateReviewStoryOpts struct {
	FileName  string             `json:"file_name" validate:"required"`
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	BrandID   primitive.ObjectID `json:"brand_id"`
	CatalogID primitive.ObjectID `json:"catalog_id" validate:"required"`
	Rating    *uint              `json:"rating" validate:"required,gte=0,lte=5"`
}

type GetReviewStoryUploadURLResp struct {
	UploadURL string             `json:"token"`
	MediaID   primitive.ObjectID `json:"id"`
}

type GetReviewStoryUploadURLBodyResp struct {
	Success bool                         `json:"success"`
	Payload *GetReviewStoryUploadURLResp `json:"payload"`
	Error   []ErrorCMS                   `json:"error"`
}

type CreateReviewStoryResp struct {
	ID        primitive.ObjectID `json:"id"`
	UploadURL string             `json:"upload_url"`
}

type CreateVideoReviewContentOpts struct {
	FileName  string             `json:"file_name"`
	UserID    primitive.ObjectID `json:"user_id"`
	BrandID   primitive.ObjectID `json:"brand_id"`
	CatalogID primitive.ObjectID `json:"catalog_id"`
}

type ReviewUserInfoResp struct {
	Success bool            `json:"success"`
	Payload *ReviewUserInfo `json:"payload"`
	Error   []ErrorCMS      `json:"error"`
}

type ReviewUserInfo struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	FullName     string             `json:"full_name,omitempty"`
	ProfileImage *Img               `json:"profile_image,omitempty"`
}

type ReviewStoryFullMessage struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CatalogID primitive.ObjectID `json:"catalog_id,omitempty" bson:"catalog_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`

	BrandID     primitive.ObjectID      `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	Rating      *uint                   `json:"rating,omitempty" bson:"rating,omitempty"`
	StoryID     primitive.ObjectID      `json:"story_id,omitempty" bson:"story_id,omitempty"`
	IsProcessed bool                    `json:"is_processed,omitempty" bson:"is_processed,omitempty"`
	CreatedAt   time.Time               `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time               `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	StoryInfo   *CatalogContentInfoResp `json:"story_info,omitempty"`
	UserInfo    *ReviewUserInfo         `json:"user_info,omitempty"`
}
