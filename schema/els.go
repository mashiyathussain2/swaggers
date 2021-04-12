package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetBrandBasicESEesp struct {
	ID               primitive.ObjectID `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	Logo             *model.IMG         `json:"logo,omitempty"`
	IsFollowedByUser bool               `json:"is_followed_by_user,omitempty"`
}
type GetBrandInfoEsResp struct {
	ID               primitive.ObjectID   `json:"id,omitempty"`
	Name             string               `json:"name,omitempty"`
	LName            string               `json:"lname,omitempty"`
	Website          string               `json:"website,omitempty"`
	Logo             *model.IMG           `json:"logo,omitempty"`
	FollowersCount   uint                 `json:"followers_count,omitempty"`
	Bio              string               `json:"bio,omitempty"`
	CoverImg         *model.IMG           `json:"cover_img,omitempty"`
	SocialAccount    *model.SocialAccount `json:"social_account,omitempty"`
	CreatedAt        time.Time            `json:"created_at,omitempty"`
	UpdatedAt        time.Time            `json:"updated_at,omitempty"`
	IsFollowedByUser bool                 `json:"is_followed_by_user,omitempty"`
}

type GetBrandsByIDBasicOpts struct {
	IDs        []string `json:"ids"`
	CustomerID primitive.ObjectID
}

type GetBrandsInfoByIDOpts struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID `json:"user_id"`
}

type GetInfluencerBasicESEesp struct {
	ID               primitive.ObjectID `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	ProfileImage     *model.IMG         `json:"profile_image,omitempty"`
	IsFollowedByUser bool               `json:"is_followed_by_user,omitempty"`
}
type GetInfluencerInfoEsResp struct {
	ID               primitive.ObjectID   `json:"id,omitempty"`
	Name             string               `json:"name,omitempty"`
	CoverImg         *model.IMG           `json:"cover_img,omitempty"`
	ProfileImage     *model.IMG           `json:"profile_image,omitempty"`
	SocialAccount    *model.SocialAccount `json:"social_account,omitempty"`
	ExternalLinks    []string             `json:"external_links,omitempty"`
	Bio              string               `json:"bio,omitempty"`
	FollowersCount   uint                 `json:"followers_count"`
	CreatedAt        time.Time            `json:"created_at,omitempty"`
	UpdatedAt        time.Time            `json:"updated_at,omitempty"`
	IsFollowedByUser bool                 `json:"is_followed_by_user,omitempty"`
}

type GetInfluencersByIDBasicOpts struct {
	IDs        []string `json:"ids"`
	CustomerID primitive.ObjectID
}

type GetInfluencerInfoByIDOpts struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID
}
