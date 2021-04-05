package api

import (
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.CreateInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
	return
}

func (a *API) editInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencersByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByIDOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencersByID(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerByName(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByNameOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerByName(s.Name)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencersBasic(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByIDBasicOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetInfluencersByIDBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["influencerID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid influencer id:%s in url", mux.Vars(r)["influencerID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetInfluencerInfoByID(&schema.GetInfluencerInfoByIDOpts{ID: id})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
