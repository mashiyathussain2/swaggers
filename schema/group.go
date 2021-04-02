package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateCatalogGroupOpts contains fields which are passed on to create catalog group
type CreateCatalogGroupOpts struct {
	Basis string               `json:"basis" validate:"required"`
	IDs   []primitive.ObjectID `json:"ids" validate:"required,min=1"`
}

//CreateGroupResp contain group data to be returned
type CreateGroupResp struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Basis      string               `json:"basis,omitempty" bson:"basis,omitempty"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	Status     string               `json:"status,omitempty" bson:"status,omitempty"`
}

//GroupResp contains the response for GetGroups
type GroupResp struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Status      model.GroupStatus  `json:"status" bson:"status,omitempty" `
	Basis       string             `json:"basis" bson:"basis,omitempty" `
	Minimum     float32            `json:"minimum" bson:"minimum,omitempty" `
	Maximum     float32            `json:"maximum" bson:"maximum,omitempty" `
	CatalogInfo []model.Catalog    `json:"catalog_info" bson:"catalog_info,omitempty" `
}

//GetCatalogByGroupIDResp contains fields for Catalogs to be returned in the GetGroupByIDResp
type GetCatalogByGroupIDResp struct {
	ID            primitive.ObjectID          `json:"id" bson:"_id,omitempty"`
	Name          string                      `json:"name,omitempty" bson:"name,omitempty"`
	BrandID       primitive.ObjectID          `json:"brand_id,omitempty" bson:"brand_id,omitempty"`
	Paths         []model.Path                `json:"category_path,omitempty" bson:"category_path,omitempty"`
	BasePrice     *model.Price                `json:"base_price,omitempty" bson:"base_price,omitempty"`
	RetailPrice   *model.Price                `json:"retail_price,omitempty" bson:"retail_price,omitempty"`
	FeaturedImage *model.CatalogFeaturedImage `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
}

//GetGroupsOpts contains fields to retrieve Groups
type GetGroupsOpts struct {
	Page   int    `json:"page" validate:"gte=0"`
	Status string `json:"status" validate:"oneof=unlist archive publish all"`
}

//GetGroupsByCatalogIDOpts contains fields to retrieve Groups
type GetGroupsByCatalogIDOpts struct {
	ID   primitive.ObjectID `json:"id,omitempty" validate:"required"`
	Page int                `json:"page" validate:"gte=0"`
}

//KeeperGetGroupsByCatalogIDOpts contains fields to retrieve Groups
type KeeperGetGroupsByCatalogIDOpts struct {
	ID     primitive.ObjectID `json:"id,omitempty" validate:"required"`
	Page   int                `json:"page" validate:"gte=0"`
	Status string             `json:"status"`
}

//AddCatalogsInTheGroupOpts contains fields to add Catalogs to the Group with given ID.
type AddCatalogsInTheGroupOpts struct {
	ID         primitive.ObjectID   `json:"id" validate:"required"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids" validate:"gt=0"`
}

//UpdateGroupStatusOpts contains fields to update group status.
type UpdateGroupStatusOpts struct {
	ID     primitive.ObjectID `json:"id" validate:"required"`
	Status string             `json:"status" validate:"required,oneof=publish unlist archive"`
}

//GetGroupsByCatalogIDResp contains the response for GetGroups
type GetGroupsByCatalogIDResp struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Basis       string             `json:"basis" bson:"basis,omitempty" `
	CatalogInfo []model.Catalog    `json:"catalog_info" bson:"catalog_info,omitempty" `
}

type GetGroupsByCatalogNameResp struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Minimum     int                `json:"minimum,omitempty" bson:"minimum,omitempty"`
	Maximum     int                `json:"maximum,omitempty" bson:"maximum,omitempty"`
	GroupStatus model.GroupStatus  `json:"group_status,omitempty" bson:"group_status,omitempty"`
	CatalogInfo []model.Catalog    `json:"catalog_info,omitempty" bson:"catalog_info,omitempty"`
	Count       int                `json:"count,omitempty" bson:"count,omitempty"`
}

//EditGroupOpts contains fields to add Catalogs to the Group with given ID.
type EditGroupOpts struct {
	ID         primitive.ObjectID   `json:"id" validate:"required"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids"`
	Basis      string               `json:"basis"`
}

//EditGroupResp contains fields to add Catalogs to the Group with given ID.
type EditGroupResp struct {
	ID         primitive.ObjectID   `json:"id" validate:"required"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids"`
	Basis      string               `json:"basis"`
	Status     model.GroupStatus    `json:"status"`
}

type GroupChangeKafkaMessage struct {
	ID         primitive.ObjectID   `json:"_id,omitempty"`
	Basis      string               `json:"basis"`
	CatalogIDs []primitive.ObjectID `json:"catalog_ids,omitempty"`
	Status     model.GroupStatus    `json:"status"`
	CreatedAt  time.Time            `json:"created_at,omitempty"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty"`
}
