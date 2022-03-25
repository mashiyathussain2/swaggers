package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/pasztorpisti/qs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) addInfluencerProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddInfluencerProductsOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var err error
	s.InfluencerID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err = a.App.InfluencerProducts.AddInfluencerProductsOpts(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}

func (a *API) removeInfluencerProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.RemoveInfluencerProductsOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var err error
	s.InfluencerID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err = a.App.InfluencerProducts.RemoveInfluencerProductsOpts(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}

func (a *API) getInfluencerProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerProducts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var resp *schema.GetInfluencerProductESResp
	var err error
	if s.Type == "self" {
		s.InfluencerID = requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID
		resp, err = a.App.InfluencerProducts.GetInfluencerProductsOpts(s.InfluencerID)
		fmt.Println("resp from db", resp)
	} else {
		if errs := a.Validator.Validate(&s); errs != nil {
			requestCTX.SetErrs(errs, http.StatusBadRequest)
			return
		}
		resp, err = a.App.Elasticsearch.GetInfluencerProducts(s.InfluencerID, s.Page)
	}
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusCreated)
}
