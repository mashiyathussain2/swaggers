//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_brand.go -package=mock go-app/app Brand

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go-app/schema"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Brand service contains all the CRUD operations related to brands
type Brand interface {
	CheckBrandIDExists(ctx context.Context, id primitive.ObjectID) (bool, error)
}

// BrandImpl implements brand service methods
type BrandImpl struct {
	App    *App
	Logger *zerolog.Logger
}

// BrandOpts contains args required to create a new instance of brand service
type BrandOpts struct {
	App    *App
	Logger *zerolog.Logger
}

// InitBrand returns brand service implementation instance
func InitBrand(opts *BrandOpts) Brand {
	return &BrandImpl{
		App:    opts.App,
		Logger: opts.Logger,
	}
}

// // CheckBrandIDExists return true/false based on if passed id exists in brand collection
// func (b *BrandImpl) CheckBrandIDExists(ctx context.Context, id primitive.ObjectID) (bool, error) {
// 	filter := bson.M{
// 		"_id": id,
// 	}
// 	count, err := b.DB.Collection(model.BrandColl).CountDocuments(ctx, filter)
// 	if err != nil {
// 		return false, err
// 	}
// 	if count != 0 {
// 		return true, nil
// 	}
// 	return false, nil
// }

// CheckBrandIDExists return true/false based on if passed id exists in brand collection
func (b *BrandImpl) CheckBrandIDExists(ctx context.Context, id primitive.ObjectID) (bool, error) {

	url := fmt.Sprintf("%s/api/keeper/brand/%s/check", b.App.Config.HypdApiConfig.EntityApi, id.Hex())
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	var res schema.CheckBrandIDExistsResp
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return false, errors.Wrap(err, "unable to decode response")
	}
	if res.Status != false {
		return false, errors.Errorf("unable to check if brand id exists")
	}

	return res.Found, nil
}
