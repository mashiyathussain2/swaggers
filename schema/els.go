package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model GetBrandBasicESEesp
type GetBrandBasicESEesp struct {
	ID               primitive.ObjectID `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	Username         string             `json:"username,omitempty"`
	Logo             *model.IMG         `json:"logo,omitempty"`
	IsFollowedByUser bool               `json:"is_followed_by_user,omitempty"`
}

// swagger:model GetBrandsListOpts
type GetBrandsListOpts struct {
	Page int `json:"page,omitempty"`
	Size int `json:"size,omitempty"`
}

// swagger:model GetActiveBrandsListESEesp
type GetActiveBrandsListESEesp struct {
	ID   primitive.ObjectID `json:"id,omitempty"`
	Name string             `json:"name,omitempty"`
	Logo *model.IMG         `json:"logo,omitempty"`
}

// swagger:model GetBrandInfoEsResp
type GetBrandInfoEsResp struct {
	ID               primitive.ObjectID   `json:"id,omitempty"`
	Name             string               `json:"name,omitempty"`
	LName            string               `json:"lname,omitempty"`
	Username         string               `json:"username,omitempty"`
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

// swagger:model GetBrandsByIDBasicOpts
type GetBrandsByIDBasicOpts struct {
	IDs        []string `json:"ids"`
	CustomerID primitive.ObjectID
}

// swagger:model GetBrandsInfoByIDOpts
type GetBrandsInfoByIDOpts struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID `json:"user_id"`
}

//swagger:model GetBrandsByUsernameBasicOpts
type GetBrandsByUsernameBasicOpts struct {
	Usernames  []string `json:"usernames"`
	CustomerID primitive.ObjectID
}

// swagger:model GetBrandsInfoByUsernameOpts
type GetBrandsInfoByUsernameOpts struct {
	Username   string             `json:"username"`
	CustomerID primitive.ObjectID `json:"user_id"`
}

// swagger:model GetInfluencerBasicESEesp
type GetInfluencerBasicESEesp struct {
	ID               primitive.ObjectID `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	Username         string             `json:"username,omitempty"`
	ProfileImage     *model.IMG         `json:"profile_image,omitempty"`
	IsFollowedByUser bool               `json:"is_followed_by_user,omitempty"`
}

// swagger:model GetInfluencerInfoEsResp
type GetInfluencerInfoEsResp struct {
	ID               primitive.ObjectID             `json:"id,omitempty"`
	Name             string                         `json:"name,omitempty"`
	Username         string                         `json:"username,omitempty"`
	CoverImg         *model.IMG                     `json:"cover_img,omitempty"`
	ProfileImage     *model.IMG                     `json:"profile_image,omitempty"`
	SocialAccount    *model.SocialAccount           `json:"social_account,omitempty"`
	ExternalLinks    []string                       `json:"external_links"`
	Bio              string                         `json:"bio"`
	FollowersCount   uint                           `json:"followers_count"`
	CreatedAt        time.Time                      `json:"created_at,omitempty"`
	UpdatedAt        time.Time                      `json:"updated_at,omitempty"`
	IsFollowedByUser bool                           `json:"is_followed_by_user,omitempty"`
	ContentCount     *GetInfluencerContentCountResp `json:"content_count,omitempty"`
}

// swagger:model GetInfluencersByIDBasicOpts
type GetInfluencersByIDBasicOpts struct {
	IDs        []string `json:"ids"`
	CustomerID primitive.ObjectID
}

// swagger:model GetInfluencerInfoByIDOpts
type GetInfluencerInfoByIDOpts struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID
}

// swagger:model GetInfluencersByUsernameBasicOpts
type GetInfluencersByUsernameBasicOpts struct {
	Usernames  []string `json:"usernames"`
	CustomerID primitive.ObjectID
}

// swagger:model GetInfluencerInfoByUsernameOpts
type GetInfluencerInfoByUsernameOpts struct {
	Username   string `json:"username"`
	CustomerID primitive.ObjectID
}

type GetInfluencerContentCount struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type GetInfluencerContentCountResp struct {
	Pebbles     int64   `json:"pebbles"`
	Products    float64 `json:"products"`
	Collections int64   `json:"collections"`
}
