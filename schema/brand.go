package schema

import (
	"go-app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckBrandIDExistsResp struct {
	Success bool `json:"success"`
	Payload bool `json:"payload"`
}

// CreateBrandOpts serialize the create brand args
type CreateBrandOpts struct {
	Name             string `json:"name" validate:"required"`
	RegisteredName   string `json:"registered_name"`
	Description      string `json:"description" validate:"required"`
	WebsiteLink      string `json:"website_link" validate:"url"`
	FulfillmentEmail string `json:"fulfillment_email" validate:"required,email"`
}

// CreateBrandResp contains fields that is to be returned for create brand operation
type CreateBrandResp struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug        string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	WebsiteLink string             `json:"website_link,omitempty" bson:"website_link,omitempty"`
	Fulfillment *model.Fulfillment `json:"fulfillment,omitempty" bson:"fulfillment,omitempty"`
}

type BrandInfoResp struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Slug        string             `json:"slug,omitempty"`
	Description string             `json:"description,omitempty"`
	Logo        *model.IMG         `json:"logo"`
}

type GetBrandInfoResp struct {
	Success bool            `json:"success"`
	Payload []BrandInfoResp `json:"payload"`
}
