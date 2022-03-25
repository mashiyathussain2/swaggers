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

func (a *API) createInfluencerCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	// s.InfluencerID,err:=
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.InfluencerCollection.CreateInfluencerCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) keeperGetInfluencerCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	fmt.Println(s)
	res, err := a.App.InfluencerCollection.KeeperGetInfluencerCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) editInfluencerCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.InfluencerCollection.EditInfluencerCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) getActiveInfluencerCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetActiveInfluencerCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetActiveInfluencerCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) appGetInfluencerCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetActiveInfluencerCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetInfluencerCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) createInfluencerCollectionApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerCollectionOpts
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
	res, err := a.App.InfluencerCollection.CreateInfluencerCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) editInfluencerCollectionApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerCollectionAppOpts
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
	res, err := a.App.InfluencerCollection.EditInfluencerCollectionApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) getActiveInfluencerCollectionByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resp, err := a.App.Elasticsearch.GetActiveInfluencerCollectionByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
