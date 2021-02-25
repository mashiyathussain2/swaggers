package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

//LabelOpts will hold the keywords related to pebbles.
type LabelOpts struct {
	Interests []string `json:"interests" validate:"required,min=1"`
	AgeGroup  []string `json:"age_group"`
	Gender    []string `json:"gender" validate:"required,min=1,dive,oneof=M F O"`
}

// CreatePebbleOpts contains and validates args required to create a pebble
type CreatePebbleOpts struct {
	Caption       string               `json:"caption" validate:"required"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" validate:"required,min=1"`
	BrandIDs      []primitive.ObjectID `json:"brand_ids" validate:"required,min=1"`
	CatalogIDs    []primitive.ObjectID `json:"catalog_ids"`
	Label         *LabelOpts           `json:"label" validate:"required"`
}

//CreatePebbleResp returns token required for uploading the content to S3 in the background
type CreatePebbleResp struct {
	ID    primitive.ObjectID `json:"id"`
	Token string             `json:"token"`
}
